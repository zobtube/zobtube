package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gitlab.com/zobtube/zobtube/internal/model"
)

func (c *Controller) Home(g *gin.Context) {
	// get clips
	var clips []model.Video
	c.datastore.Where("type = ?", "c").Limit(12).Offset(0).Order("created_at desc").Find(&clips)

	// get videos
	var videos []model.Video
	c.datastore.Where("type = ?", "v").Limit(12).Offset(0).Order("created_at desc").Find(&videos)

	// get movies
	var movies []model.Video
	c.datastore.Where("type = ?", "m").Limit(12).Offset(0).Order("created_at desc").Find(&movies)

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
		"User":   g.MustGet("user").(*model.User),
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
