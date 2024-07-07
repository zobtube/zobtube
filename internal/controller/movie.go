package controller

import (
	"github.com/gin-gonic/gin"
)

func (c *Controller) MovieList(g *gin.Context) {
	c.GenericVideoList("movie", g)
}

func (c *Controller) MovieView(g *gin.Context) {
	c.GenericVideoView("movie", g)
}

func (c *Controller) MovieStream(g *gin.Context) {
	c.GenericVideoStream("movie", g)
}

func (c *Controller) MovieThumb(g *gin.Context) {
	c.GenericVideoThumb("movie", g)
}

func (c *Controller) MovieThumbXS(g *gin.Context) {
	c.GenericVideoThumbXS("movie", g)
}
