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

	var tasks []model.Task
	c.datastore.Limit(5).Order("created_at DESC").Find(&tasks)

	g.HTML(http.StatusOK, "adm/home.html", gin.H{
		"User":         g.MustGet("user").(*model.User),
		"Version":      ZT_VERSION,
		"VideoCount":   videoCount,
		"ActorCount":   actorCount,
		"ChannelCount": channelCount,
		"Tasks":        tasks,
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

func (c *Controller) AdmTaskView(g *gin.Context) {
	// get id from path
	id := g.Param("id")

	// get item from ID
	task := &model.Task{
		ID: id,
	}
	result := c.datastore.First(task)

	// check result
	if result.RowsAffected < 1 {
		g.JSON(404, gin.H{})
		return
	}

	g.HTML(http.StatusOK, "adm/task-view.html", gin.H{
		"User": g.MustGet("user").(*model.User),
		"Task": task,
	})
}

func (c *Controller) AdmTaskList(g *gin.Context) {
	tasks := []model.Task{}
	result := c.datastore.Find(&tasks)

	// check result
	if result.RowsAffected < 1 {
		g.JSON(404, gin.H{})
		return
	}

	g.HTML(http.StatusOK, "adm/task-list.html", gin.H{
		"User":  g.MustGet("user").(*model.User),
		"Tasks": tasks,
	})
}
