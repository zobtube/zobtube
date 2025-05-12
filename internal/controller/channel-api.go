package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/zobtube/zobtube/internal/model"
)

func (c *Controller) ChannelAjaxList(g *gin.Context) {
	// get item from ID
	channels := []model.Channel{}
	c.datastore.Find(&channels)

	channelsJSON := make(map[string]string)
	for _, channel := range channels {
		channelsJSON[channel.ID] = channel.Name
	}

	g.JSON(200, gin.H{
		"channels": channelsJSON,
	})
}
