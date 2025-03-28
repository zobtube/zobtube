package main

import (
	"embed"
	"errors"
	"fmt"
	"sync"

	"github.com/zobtube/zobtube/internal/config"
	"github.com/zobtube/zobtube/internal/controller"
	"github.com/zobtube/zobtube/internal/http"
	"github.com/zobtube/zobtube/internal/model"
	"github.com/zobtube/zobtube/internal/provider"
	"github.com/zobtube/zobtube/internal/runner"
	"github.com/zobtube/zobtube/internal/task/video"
)

//go:embed web
var webFS embed.FS

// const
const cfgPath = "config.yml"

// errors
var ErrNoUser = errors.New("database does not have any account")

// channel for http server shutdown
var wg sync.WaitGroup
var shutdownChannel chan int

func startFailsafeWebServer(err error, c controller.AbtractController) {
	// http server
	httpServer := &http.Server{}

	if err == config.ErrNoDbDriverSet || err == config.ErrNoDbConnStringSet || err == config.ErrNoMediaPathSet {
		httpServer, _ = http.NewFailsafeConfig(c, &webFS)
	} else if err == ErrNoUser {
		httpServer, _ = http.NewFailsafeUser(c, &webFS)
	} else {
		httpServer, _ = http.NewUnexpectedError(c, &webFS, err)
	}

	// handle shutdown
	go httpServer.WaitForStopSignal(shutdownChannel)

	httpServer.Start("0.0.0.0:8080")

	// Wait for all HTTP fetches to complete.
	wg.Wait()

	fmt.Println("exiting")
}

func main() {
	wg.Add(1)

	// http server
	httpServer := &http.Server{}

	// channel for http server shutdown
	shutdownChannel = make(chan int)

	// create controller
	c := controller.New(shutdownChannel)
	cfg, err := config.New(cfgPath)
	if err != nil {
		startFailsafeWebServer(err, c)
		return
	}

	err = cfg.EnsureTreePresent()
	if err != nil {
		startFailsafeWebServer(err, c)
		return
	}
	c.ConfigurationRegister(cfg)

	db, err := model.New(cfg)
	if err != nil {
		startFailsafeWebServer(err, c)
		return
	}
	c.DatabaseRegister(db)

	// check if at least one user exists
	var count int64
	db.Model(&model.User{}).Count(&count)
	if count == 0 {
		startFailsafeWebServer(ErrNoUser, c)
		return
	}

	c.ProviderRegister(&provider.BabesDirectory{})
	c.ProviderRegister(&provider.Babepedia{})
	c.ProviderRegister(&provider.Boobpedia{})
	c.ProviderRegister(&provider.Pornhub{})

	go c.CleanupRoutine()

	runner := &runner.Runner{}
	runner.RegisterTask(video.NewVideoCreating())
	runner.RegisterTask(video.NewVideoDeleting())
	runner.Start(cfg, db)
	c.RunnerRegister(runner)

	// create http server
	httpServer, _ = http.New(&c, &webFS)

	// serve content
	httpServer.Start(cfg.Server.Bind)
}
