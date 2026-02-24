package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	swagv2 "github.com/swaggo/swag/v2"

	"github.com/zobtube/zobtube/internal/controller"
	_ "github.com/zobtube/zobtube/internal/swagger"
)

// main http server setup
func (server *Server) ControllerSetupDefault(c *controller.AbstractController) {
	server.setupRoutes(*c)
	// both next settings are needed for filepath used above
	server.Router.UseRawPath = true
	server.Router.UnescapePathValues = false
	server.Router.RemoveExtraSlash = false

	// serve swagger documentation from embedded internal/swagger docs (swag v2).
	// gin-swagger uses swag v1, so we serve doc.json ourselves via swag v2 and point the UI at it.
	// /swagger redirects to /swagger/ so relative URLs (./swagger-ui.css) resolve correctly.
	// Register /swagger/*any first to avoid Gin route conflict; handle doc.json inside the handler.
	swaggerHandler := ginSwagger.WrapHandler(swaggerfiles.Handler, ginSwagger.URL("/swagger/doc.json"))
	server.Router.GET("/swagger", func(g *gin.Context) {
		g.Redirect(http.StatusFound, "/swagger/")
	})
	server.Router.GET("/swagger/*any", func(g *gin.Context) {
		any := g.Param("any")
		if any == "/doc.json" {
			doc, err := swagv2.ReadDoc(swagv2.Name)
			if err != nil {
				g.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			g.Data(http.StatusOK, "application/json", []byte(doc))
			return
		}
		if any == "/" || any == "" {
			g.Request.URL.Path = "/swagger/index.html"
			g.Request.RequestURI = "/swagger/index.html"
		}
		swaggerHandler(g)
	})
}

// failsafe http server setup - unexpected error
func (server *Server) ControllerSetupFailsafeError(c controller.AbstractController, faultyError error) {
	// server is not healthy
	server.healthy = false

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
