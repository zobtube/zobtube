package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/zobtube/zobtube/internal/model"
)

func (c *Controller) Home(g *gin.Context) {
	// get videos
	var videos []model.Video
	c.datastore.Where("type = ?", "v").Order("created_at desc").Find(&videos)

	c.HTML(g, http.StatusOK, "home/home.html", gin.H{
		"Videos":    videos,
		"VideoType": "video",
	})
}
