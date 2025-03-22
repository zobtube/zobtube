package main

import (
	"embed"
	"errors"

	"github.com/zobtube/zobtube/internal/config"
	"github.com/zobtube/zobtube/internal/controller"
	"github.com/zobtube/zobtube/internal/http"
	"github.com/zobtube/zobtube/internal/model"
	"github.com/zobtube/zobtube/internal/provider"
)

//go:embed web
var webFS embed.FS

var ErrNoUser = errors.New("database does not have any account")

func startFailsafeWebServer(err error, c controller.AbtractController) {
	var httpServer *http.Server

	if err == config.ErrNoDbDriverSet || err == config.ErrNoDbConnStringSet || err == config.ErrNoMediaPathSet {
		httpServer, _ = http.NewFailsafeConfig(c, &webFS)
	} else if err == ErrNoUser {
		httpServer, _ = http.NewFailsafeUser(c, &webFS)
	} else {
		httpServer, _ = http.NewUnexpectedError(c, &webFS, err)
	}

	httpServer.Start("0.0.0.0:8080")
}

func main() {
	cfg, err := config.New()
	if err != nil {
		startFailsafeWebServer(err, controller.New(nil, nil))
		return
	}

	err = cfg.EnsureTreePresent()
	if err != nil {
		startFailsafeWebServer(err, controller.New(nil, nil))
		return
	}

	db, err := model.New(cfg)
	if err != nil {
		startFailsafeWebServer(err, controller.New(cfg, nil))
		return
	}

	// check if at least one user exists
	var count int64
	db.Model(&model.User{}).Count(&count)
	if count == 0 {
		startFailsafeWebServer(ErrNoUser, controller.New(cfg, db))
		return
	}

	c := controller.New(cfg, db)
	c.ProviderRegister(&provider.BabesDirectory{})
	c.ProviderRegister(&provider.Babepedia{})
	c.ProviderRegister(&provider.Boobpedia{})
	c.ProviderRegister(&provider.Pornhub{})

	go c.CleanupRoutine()

	httpServer, _ := http.New(&c, &webFS)
	httpServer.Start(cfg.Server.Bind)
}
