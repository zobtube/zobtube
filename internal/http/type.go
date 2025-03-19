package http

import (
	"embed"

	"github.com/gin-gonic/gin"

	"github.com/zobtube/zobtube/internal/controller"
)

type Server struct {
	Server *gin.Engine
	FS     *embed.FS
}

func New(c *controller.AbtractController, fs *embed.FS) (*Server, error) {
	server := &Server{
		Server: gin.Default(),
		FS:     fs,
	}

	server.setupRoutes(*c)
	// both next settings are needed for filepath used above
	server.Server.UseRawPath = true
	server.Server.UnescapePathValues = false
	server.Server.RemoveExtraSlash = false

	return server, nil
}
