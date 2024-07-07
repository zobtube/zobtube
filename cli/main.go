package main

import (
	"gitlab.com/zobtube/zobtube/internal/config"
	"gitlab.com/zobtube/zobtube/internal/controller"
	"gitlab.com/zobtube/zobtube/internal/http"
	"gitlab.com/zobtube/zobtube/internal/model"
	"gitlab.com/zobtube/zobtube/internal/provider"
)

func main() {
	cfg, err := config.New()
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

	httpServer, _ := http.New(&c)
	httpServer.Start(":8080")
}
