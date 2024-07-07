package controller

import (
	"github.com/gin-gonic/gin"
)

func (c *Controller) ClipList(g *gin.Context) {
	/*
	   _clip := &model.Video{
	       Name: "test",
	       Filename: "test-filename",
	       Thumbnail: false,
	       ThumbnailMini: false,
	       Duration: 0,
	       Type: "c",
	   }
	   c.datastore.Create(_clip)
	   //fmt.Println(res.Error)
	*/

	c.GenericVideoList("clip", g)
}

func (c *Controller) ClipView(g *gin.Context) {
	c.GenericVideoView("clip", g)
}

func (c *Controller) ClipStream(g *gin.Context) {
	c.GenericVideoStream("clip", g)
}

func (c *Controller) ClipThumb(g *gin.Context) {
	c.GenericVideoThumb("clip", g)
}

func (c *Controller) ClipThumbXS(g *gin.Context) {
	c.GenericVideoThumbXS("clip", g)
}
