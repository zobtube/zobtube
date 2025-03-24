package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zobtube/zobtube/internal/model"
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
		"User":         g.MustGet("user").(*model.User),
		"Version":      ZT_VERSION,
		"VideoCount":   videoCount,
		"ActorCount":   actorCount,
		"ChannelCount": channelCount,
	})
}

func (c *Controller) AdmVideoList(g *gin.Context) {
	var videos []model.Video

	c.datastore.Find(&videos)

	g.HTML(http.StatusOK, "adm/object-list.html", gin.H{
		"User":       g.MustGet("user").(*model.User),
		"ObjectName": "Video",
		"Objects":    videos,
	})
}

func (c *Controller) AdmActorList(g *gin.Context) {
	var actors []model.Actor

	c.datastore.Find(&actors)

	g.HTML(http.StatusOK, "adm/object-list.html", gin.H{
		"User":       g.MustGet("user").(*model.User),
		"ObjectName": "Actor",
		"Objects":    actors,
	})
}

func (c *Controller) AdmChannelList(g *gin.Context) {
	var channels []model.Channel

	c.datastore.Find(&channels)

	g.HTML(http.StatusOK, "adm/object-list.html", gin.H{
		"User":       g.MustGet("user").(*model.User),
		"ObjectName": "Channel",
		"Objects":    channels,
	})
}
