package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/zobtube/zobtube/internal/controller"
)

// main http server setup
func (server *Server) ControllerSetupDefault(c *controller.AbtractController) {
	server.setupRoutes(*c)
	// both next settings are needed for filepath used above
	server.Router.UseRawPath = true
	server.Router.UnescapePathValues = false
	server.Router.RemoveExtraSlash = false
}

// failsafe http server setup - unexpected error
func (server *Server) ControllerSetupFailsafeError(c controller.AbtractController, faultyError error) {
	// redirect everything on '/'
	server.Router.NoRoute(func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/")
	})

	// server the error page
	server.Router.GET("", func(g *gin.Context) {
		g.HTML(http.StatusOK, "failsafe/error.html", gin.H{
			"Error": faultyError,
		})
	})
}
