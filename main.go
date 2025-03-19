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

func main() {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	err = cfg.EnsureTreePresent()
	if err != nil {
		panic(err)
	}

	db, err := model.New(cfg)
	if err != nil {
		panic(err)
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
