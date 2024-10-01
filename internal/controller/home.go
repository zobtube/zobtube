package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gitlab.com/zobtube/zobtube/internal/model"
)

func (c *Controller) Home(g *gin.Context) {
	// get videos
	var videos []model.Video
	c.datastore.Where("type = ?", "v").Order("created_at desc").Find(&videos)

	g.HTML(http.StatusOK, "home/home.html", gin.H{
		"User":      g.MustGet("user").(*model.User),
		"Videos":    videos,
		"VideoType": "video",
	})
}
