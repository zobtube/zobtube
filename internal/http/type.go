package http

import (
	"embed"
	"io/fs"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Server struct {
	Server         *http.Server
	Router         *gin.Engine
	FS             *embed.FS
	Logger         *zerolog.Logger
	authentication bool
	healthy        bool
}

func New(embedfs *embed.FS, ginDebug bool, logger *zerolog.Logger) *Server {
	// only show gin debugging if ginDebug is set to true
	if ginDebug {
		logger.Warn().Msg("gin debugging mode activated")
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// create common server
	server := &Server{
		Router:         gin.New(),
		FS:             embedfs,
		Logger:         logger,
		authentication: false,
		healthy:        true,
	}

	// add recovery middleware
	server.Router.Use(gin.Recovery())

	// setup logger
	server.Router.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		end := time.Now()
		latency := end.Sub(start)

		logger.Info().
			Str("method", c.Request.Method).
			Str("path", c.Request.RequestURI).
			Int("status", c.Writer.Status()).
			Str("ip", c.ClientIP()).
			Dur("latency", latency).
			Msg("http request")
	})

	// template error handling
	server.Router.Use(func(c *gin.Context) {
		c.Next()
		err := c.Errors.Last()
		if err != nil {
			logger.Error().Err(err).Msg("html template rendering failed")
		}
	})

	// load templates (web/page/** matches files in page and in subdirs e.g. app.html, actor/list.html)
	server.LoadHTMLFromEmbedFS("web/page/**")

	// prepare subfs
	staticFS, _ := fs.Sub(server.FS, "web/static")

	// load static
	server.Router.StaticFS("/static", http.FS(staticFS))
	server.Router.GET("/ping", server.livenessProbe)

	return server
}
