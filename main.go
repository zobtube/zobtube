package main

import (
	"embed"

	"github.com/zobtube/zobtube/internal/config"
	"github.com/zobtube/zobtube/internal/controller"
	"github.com/zobtube/zobtube/internal/http"
	"github.com/zobtube/zobtube/internal/model"
	"github.com/zobtube/zobtube/internal/provider"
)

//go:embed web
var webFS embed.FS

func startFailsafeWebServer(err error) {
	var httpServer *http.Server

	if err == config.ErrNoDbDriverSet || err == config.ErrNoDbConnStringSet || err == config.ErrNoMediaPathSet {
		c := controller.New(nil, nil)
		httpServer, _ = http.NewFailsafeConfig(c, &webFS)
	} else {
		c := controller.New(nil, nil)
		httpServer, _ = http.NewUnexpectedError(c, &webFS, err)
		//httpServer.registerError(err)
	}

	httpServer.Start("0.0.0.0:8080")
}

func main() {
	cfg, err := config.New()
	if err != nil {
		startFailsafeWebServer(err)
		return
	}

	err = cfg.EnsureTreePresent()
	if err != nil {
		startFailsafeWebServer(err)
		return
	}

	db, err := model.New(cfg)
	if err != nil {
		startFailsafeWebServer(err)
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
