package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gitlab.com/zobtube/zobtube/internal/model"
)

func (c *Controller) Home(g *gin.Context) {
	// get clips
	var clips []model.Video
	c.datastore.Where("type = ?", "c").Find(&clips).Limit(12).Offset(0).Order("created_at")

	// get videos
	var videos []model.Video
	c.datastore.Where("type = ?", "v").Find(&videos).Limit(12).Offset(0).Order("created_at")

	// get movies
	var movies []model.Video
	c.datastore.Where("type = ?", "m").Find(&movies).Limit(12).Offset(0).Order("created_at")

	//TODO: get comics

	// get counts
	var (
		clipCount  int64
		movieCount int64
		videoCount int64
	)
	c.datastore.Table("videos").Where("type = ?", "c").Count(&clipCount)
	c.datastore.Table("videos").Where("type = ?", "m").Count(&movieCount)
	c.datastore.Table("videos").Where("type = ?", "v").Count(&videoCount)

	g.HTML(http.StatusOK, "home/home.html", gin.H{
		"User":   user,
		"Clips":  clips,
		"Movies": movies,
		"Videos": videos,
		"Counts": map[string]int64{
			"clips":  clipCount,
			"videos": videoCount,
			"movies": movieCount,
			"comics": 0,
		},
	})
}
