package http

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/zobtube/zobtube/internal/controller"
)

type Server struct {
	Server *gin.Engine
	FS     *embed.FS
}

// main http server setup
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

// failsafe http server setup - no valid config found
func NewFailsafeConfig(c controller.AbtractController, embedfs *embed.FS) (*Server, error) {
	server := &Server{
		Server: gin.Default(),
		FS:     embedfs,
	}

	// load templates
	server.LoadHTMLFromEmbedFS("web/page/**/*")

	// prepare subfs
	staticFS, _ := fs.Sub(server.FS, "web/static")

	// load static
	server.Server.StaticFS("/static", http.FS(staticFS))
	server.Server.GET("/ping", livenessProbe)

	// failsafe configuration route
	server.Server.GET("", c.FailsafeConfiguration)
	server.Server.POST("", c.FailsafeConfiguration)

	return server, nil
}

// failsafe http server setup - unexpected error
func NewUnexpectedError(c controller.AbtractController, embedfs *embed.FS, faultyError error) (*Server, error) {
	server := &Server{
		Server: gin.Default(),
		FS:     embedfs,
	}

	// load templates
	server.LoadHTMLFromEmbedFS("web/page/**/*")

	// prepare subfs
	staticFS, _ := fs.Sub(server.FS, "web/static")

	// load static
	server.Server.StaticFS("/static", http.FS(staticFS))
	server.Server.GET("/ping", livenessProbe)

	server.Server.GET("", func(g *gin.Context) {
		g.HTML(http.StatusOK, "failsafe/error.html", gin.H{
			"Error": faultyError,
		})
	})

	return server, nil
}

// failsafe http server setup - no valid config found
func NewFailsafeUser(c controller.AbtractController, embedfs *embed.FS) (*Server, error) {
	server := &Server{
		Server: gin.Default(),
		FS:     embedfs,
	}

	// load templates
	server.LoadHTMLFromEmbedFS("web/page/**/*")

	// prepare subfs
	staticFS, _ := fs.Sub(server.FS, "web/static")

	// load static
	server.Server.StaticFS("/static", http.FS(staticFS))
	server.Server.GET("/ping", livenessProbe)

	// failsafe configuration route
	server.Server.GET("", c.FailsafeUser)
	server.Server.POST("", c.FailsafeUser)

	return server, nil
}
