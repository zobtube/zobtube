package controller

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/zobtube/zobtube/internal/model"
)

// ChannelList godoc
//
//	@Summary	List all channels
//	@Tags		channel
//	@Produce	json
//	@Success	200	{object}	map[string]interface{}
//	@Router		/channel [get]
func (c *Controller) ChannelList(g *gin.Context) {
	var channels []model.Channel
	c.datastore.Find(&channels)
	g.JSON(http.StatusOK, gin.H{
		"items": channels,
		"total": len(channels),
	})
}

// ChannelGet godoc
//
//	@Summary	Get channel by ID with videos
//	@Tags		channel
//	@Produce	json
//	@Param		id	path	string	true	"Channel ID"
//	@Success	200	{object}	map[string]interface{}
//	@Failure	404	{object}	map[string]interface{}
//	@Router		/channel/{id} [get]
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

// ChannelCreate godoc
//
//	@Summary	Create a new channel
//	@Tags		channel
//	@Accept		json
//	@Param		body	body	object	true	"JSON with name"
//	@Success	201	{object}	map[string]interface{}
//	@Failure	400	{object}	map[string]interface{}
//	@Router		/channel [post]
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

// ChannelUpdate godoc
//
//	@Summary	Update channel
//	@Tags		channel
//	@Accept		json
//	@Param		id	path	string	true	"Channel ID"
//	@Param		body	body	object	true	"JSON with optional name"
//	@Success	200	{object}	map[string]interface{}
//	@Failure	400	{object}	map[string]interface{}
//	@Failure	404	{object}	map[string]interface{}
//	@Router		/channel/{id} [put]
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

// ChannelMap godoc
//
//	@Summary	Get channel ID to name map
//	@Tags		channel
//	@Produce	json
//	@Success	200	{object}	map[string]interface{}
//	@Router		/channel/map [get]
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

// ChannelThumb godoc
//
//	@Summary	Get channel thumbnail image
//	@Tags		channel
//	@Param		id	path	string	true	"Channel ID"
//	@Success	200	file	bytes
//	@Failure	404
//	@Router		/channel/{id}/thumb [get]
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
	store, err := c.storageResolver.Storage(c.config.DefaultLibraryID)
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}
	path := filepath.Join("channels", id, "thumb.jpg")
	c.serveFromStorage(g, store, path)
}
