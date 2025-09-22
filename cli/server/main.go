package server

import (
	"context"
	"embed"
	"sync"

	"github.com/rs/zerolog"
	"github.com/urfave/cli/v3"

	"github.com/zobtube/zobtube/internal/config"
	"github.com/zobtube/zobtube/internal/controller"
	"github.com/zobtube/zobtube/internal/http"
	"github.com/zobtube/zobtube/internal/model"
	"github.com/zobtube/zobtube/internal/provider"
	"github.com/zobtube/zobtube/internal/runner"
	"github.com/zobtube/zobtube/internal/task/video"
)

type Parameters struct {
	Ctx     context.Context
	Cmd     *cli.Command
	Logger  *zerolog.Logger
	Version string
	Commit  string
	Date    string
	WebFS   *embed.FS
}

// channel for http server shutdown
var wg sync.WaitGroup
var shutdownChannel chan int

func Start(params *Parameters) error {
	// setup log level
	zerolog.SetGlobalLevel(zerolog.Level(params.Cmd.Int("log-level")))

	// initialize logger
	params.Logger.Info().Msg("zobtube starting")

	// create http server
	httpServer := http.New(params.WebFS, params.Cmd.Bool("gin-debug"), params.Logger)

	wg.Add(1)

	// channel for http server shutdown
	shutdownChannel = make(chan int)

	// handle shutdown
	go httpServer.WaitForStopSignal(shutdownChannel)

	// create controller
	c := controller.New(shutdownChannel)

	// register logger
	c.LoggerRegister(params.Logger)

	// setup configuration
	cfg, err := config.New(
		params.Logger,
		params.Cmd.String("server-bind"),
		params.Cmd.String("db-driver"),
		params.Cmd.String("db-connstring"),
		params.Cmd.String("media-path"),
	)
	if err != nil {
		startFailsafeWebServer(httpServer, err, c)
		return nil
	}

	params.Logger.Debug().Str("kind", "system").Msg("ensure library folders are present")
	err = cfg.EnsureTreePresent()
	if err != nil {
		startFailsafeWebServer(httpServer, err, c)
		return nil
	}
	c.ConfigurationRegister(cfg)

	params.Logger.Debug().Str("kind", "system").Msg("apply models on database")
	db, err := model.New(cfg)
	if err != nil {
		startFailsafeWebServer(httpServer, err, c)
		return nil
	}
	c.DatabaseRegister(db)

	params.Logger.Debug().Str("kind", "system").Msg("check if at least one user exists")
	var count int64
	db.Model(&model.User{}).Count(&count)
	if count == 0 {
		// instance first start, create a default user
		params.Logger.Warn().Str("kind", "system").Msg("no user setup, creating default admin")

		newUser := &model.User{
			Username: "admin",
			Admin:    true,
		}

		// save it
		tx := db.Begin()
		err = tx.Save(&newUser).Error
		if err != nil {
			params.Logger.Error().Str("kind", "system").Err(err).Msg("unable to create initial user")
			tx.Rollback()
			startFailsafeWebServer(httpServer, err, c)
			return err
		}

		// register the instance to be authentication-less
		config := &model.Configuration{
			ID:                 1,
			UserAuthentication: false,
		}

		// save it
		err = tx.Assign(&config).FirstOrCreate(&config).Error
		if err != nil {
			params.Logger.Error().Str("kind", "system").Err(err).Msg("unable to create initial user")
			tx.Rollback()
			startFailsafeWebServer(httpServer, err, c)
			return err
		}

		tx.Commit()
	} else {
		// at least one user present
		// now checking if configuration is set (allowing migration from previous versions)
		config := &model.Configuration{}
		result := db.First(config)

		// check result
		if result.RowsAffected < 1 {
			params.Logger.Warn().Str("kind", "system").Msg("configuration unset with existing users, enabling authentication")
			// register the instance to be authentication-less
			config := &model.Configuration{
				ID:                 1,
				UserAuthentication: true,
			}

			// save it
			err = db.Assign(&config).FirstOrCreate(&config).Error
			if err != nil {
				params.Logger.Error().Str("kind", "system").Err(err).Msg("unable to create initial configuration")
				startFailsafeWebServer(httpServer, err, c)
				return err
			}
		}
	}

	// loading configuration from database
	dbconfig := &model.Configuration{}
	result := db.First(dbconfig)

	// check result
	if result.RowsAffected < 1 {
		params.Logger.Fatal().Str("kind", "system").Msg("configuration should not be empty")
		return nil
	}
	c.ConfigurationFromDBApply(dbconfig)

	// external providers
	params.Logger.Debug().Str("kind", "system").Msg("register external providers")
	c.ProviderRegister(&provider.BabesDirectory{})
	c.ProviderRegister(&provider.Babepedia{})
	c.ProviderRegister(&provider.Boobpedia{})
	c.ProviderRegister(&provider.Pornhub{})

	go c.CleanupRoutine()

	runner := &runner.Runner{}
	runner.RegisterTask(video.NewVideoCreating())
	runner.RegisterTask(video.NewVideoDeleting())
	runner.RegisterTask(video.NewVideoGenerateThumbnail())
	runner.Start(cfg, db)
	c.RunnerRegister(runner)

	c.BuildDetailsRegister(params.Version, params.Commit, params.Date)

	// register controller
	httpServer.ControllerSetupDefault(&c)

	// serve content
	httpServer.Start(cfg.Server.Bind)

	return nil
}
