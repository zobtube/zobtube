package controller

import (
	"github.com/gin-gonic/gin"
)

const ErrConfigAbsent = "configuration not found, restarting the appliaction should fix the issue"

func (c *Controller) ErrFatal(g *gin.Context, err string) {
	g.HTML(404, "err/fatal.html", gin.H{
		"error": err,
	})
}

func (c *Controller) ErrNotFound(g *gin.Context) {
	g.HTML(404, "err/not-found.html", gin.H{})
}

func (c *Controller) ErrUnauthorized(g *gin.Context) {
	c.HTML(g, 401, "err/unauthorized.html", gin.H{})
}
