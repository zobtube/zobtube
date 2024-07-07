package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (c *Controller) AdmHome(g *gin.Context) {
	// get counts
	var (
		videoCount   int64
		actorCount   int64
		channelCount int64
	)

	c.datastore.Table("videos").Count(&videoCount)
	c.datastore.Table("actors").Count(&actorCount)
	c.datastore.Table("channels").Count(&channelCount)

	g.HTML(http.StatusOK, "adm/home.html", gin.H{
		"User":         user,
		"Version":      "0.0.1",
		"VideoCount":   videoCount,
		"ActorCount":   actorCount,
		"ChannelCount": channelCount,
	})
}
