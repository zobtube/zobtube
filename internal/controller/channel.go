package controller

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/zobtube/zobtube/internal/model"
)

func (c *Controller) ChannelList(g *gin.Context) {
	var channels []model.Channel
	c.datastore.Find(&channels)
	g.JSON(http.StatusOK, gin.H{
		"items": channels,
		"total": len(channels),
	})
}

func (c *Controller) ChannelGet(g *gin.Context) {
	id := g.Param("id")
	channel := &model.Channel{ID: id}
	result := c.datastore.First(channel)
	if result.RowsAffected < 1 {
		g.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	var videos []model.Video
	c.datastore.Where("channel_id = ?", channel.ID).Find(&videos)
	g.JSON(http.StatusOK, gin.H{
		"channel": channel,
		"videos":  videos,
	})
}

func (c *Controller) ChannelCreate(g *gin.Context) {
	var body struct {
		Name string `json:"name"`
	}
	if err := g.ShouldBindJSON(&body); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	channel := &model.Channel{Name: body.Name}
	if err := c.datastore.Create(channel).Error; err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	g.JSON(http.StatusCreated, gin.H{"id": channel.ID, "redirect": "/channel/" + channel.ID})
}

func (c *Controller) ChannelUpdate(g *gin.Context) {
	id := g.Param("id")
	channel := &model.Channel{ID: id}
	if result := c.datastore.First(channel); result.RowsAffected < 1 {
		g.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	var body struct {
		Name *string `json:"name"`
	}
	if err := g.ShouldBindJSON(&body); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if body.Name != nil {
		channel.Name = *body.Name
	}
	if err := c.datastore.Save(channel).Error; err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	g.JSON(http.StatusOK, channel)
}

func (c *Controller) ChannelMap(g *gin.Context) {
	channels := []model.Channel{}
	c.datastore.Find(&channels)
	channelsJSON := make(map[string]string)
	for _, channel := range channels {
		channelsJSON[channel.ID] = channel.Name
	}
	g.JSON(http.StatusOK, gin.H{
		"channels": channelsJSON,
	})
}

func (c *Controller) ChannelThumb(g *gin.Context) {
	id := g.Param("id")
	channel := &model.Channel{ID: id}
	result := c.datastore.First(channel)
	if result.RowsAffected < 1 {
		c.ErrNotFound(g)
		return
	}
	if !channel.Thumbnail {
		g.Redirect(http.StatusFound, ACTOR_PROFILE_PICTURE_MISSING)
		return
	}
	targetPath := filepath.Join(c.config.Media.Path, CHANNEL_FILEPATH, id, "thumb.jpg")
	g.File(targetPath)
}
