package http

import (
	"github.com/gin-gonic/gin"

	"gitlab.com/zobtube/zobtube/internal/controller"
)

type Server struct {
	Server *gin.Engine
}

func New(c *controller.AbtractController) (*Server, error) {
	server := &Server{
		Server: gin.Default(),
	}

	server.setupRoutes(*c)
	// both next settings are needed for filepath used above
	server.Server.UseRawPath = true
	server.Server.UnescapePathValues = false
	server.Server.RemoveExtraSlash = false

	return server, nil
}
