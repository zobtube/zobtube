package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/zobtube/zobtube/internal/model"
)

func (c *Controller) Home(g *gin.Context) {
	var videos []model.Video
	c.datastore.Where("type = ?", "v").Order("created_at desc").Find(&videos)
	g.JSON(http.StatusOK, gin.H{
		"items": videos,
		"total": len(videos),
	})
}
