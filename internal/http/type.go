package http

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/zobtube/zobtube/internal/controller"
)

type Server struct {
	Server *http.Server
	Router *gin.Engine
	FS     *embed.FS
}

// main http server setup
func New(c *controller.AbtractController, fs *embed.FS) (*Server, error) {
	server := &Server{
		Router: gin.Default(),
		FS:     fs,
	}

	server.setupRoutes(*c)
	// both next settings are needed for filepath used above
	server.Router.UseRawPath = true
	server.Router.UnescapePathValues = false
	server.Router.RemoveExtraSlash = false

	return server, nil
}

// failsafe http server setup - no valid config found
func NewFailsafeConfig(c controller.AbtractController, embedfs *embed.FS) (*Server, error) {
	server := &Server{
		Router: gin.Default(),
		FS:     embedfs,
	}

	// load templates
	server.LoadHTMLFromEmbedFS("web/page/**/*")

	// prepare subfs
	staticFS, _ := fs.Sub(server.FS, "web/static")

	// load static
	server.Router.StaticFS("/static", http.FS(staticFS))
	server.Router.GET("/ping", livenessProbe)

	// failsafe configuration route
	server.Router.GET("", c.FailsafeConfiguration)
	server.Router.POST("", c.FailsafeConfiguration)

	return server, nil
}

// failsafe http server setup - unexpected error
func NewUnexpectedError(c controller.AbtractController, embedfs *embed.FS, faultyError error) (*Server, error) {
	server := &Server{
		Router: gin.Default(),
		FS:     embedfs,
	}

	// load templates
	server.LoadHTMLFromEmbedFS("web/page/**/*")

	// prepare subfs
	staticFS, _ := fs.Sub(server.FS, "web/static")

	// load static
	server.Router.StaticFS("/static", http.FS(staticFS))
	server.Router.GET("/ping", livenessProbe)

	server.Router.GET("", func(g *gin.Context) {
		g.HTML(http.StatusOK, "failsafe/error.html", gin.H{
			"Error": faultyError,
		})
	})

	return server, nil
}

// failsafe http server setup - no valid config found
func NewFailsafeUser(c controller.AbtractController, embedfs *embed.FS) (*Server, error) {
	server := &Server{
		Router: gin.Default(),
		FS:     embedfs,
	}

	// load templates
	server.LoadHTMLFromEmbedFS("web/page/**/*")

	// prepare subfs
	staticFS, _ := fs.Sub(server.FS, "web/static")

	// load static
	server.Router.StaticFS("/static", http.FS(staticFS))
	server.Router.GET("/ping", livenessProbe)

	// failsafe configuration route
	server.Router.GET("", c.FailsafeUser)
	server.Router.POST("", c.FailsafeUser)

	return server, nil
}
