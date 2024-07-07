package controller

import (
	"github.com/gin-gonic/gin"
)

const MEDIA_VIDEO_FILEPATH = "/clips"

func (c *Controller) VideoList(g *gin.Context) {
	c.GenericVideoList("video", g)
}

func (c *Controller) VideoView(g *gin.Context) {
	c.GenericVideoView("video", g)
}

func (c *Controller) VideoStream(g *gin.Context) {
	c.GenericVideoStream("video", g)
}

func (c *Controller) VideoThumb(g *gin.Context) {
	c.GenericVideoThumb("video", g)
}

func (c *Controller) VideoThumbXS(g *gin.Context) {
	c.GenericVideoThumbXS("video", g)
}

func (c *Controller) VideoEdit(g *gin.Context) {
}
