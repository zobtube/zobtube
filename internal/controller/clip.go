package controller

import (
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zobtube/zobtube/internal/model"
)

// ClipView godoc
//
//	@Summary	Get clip view data with video, actors, categories and clip list
//	@Tags		video
//	@Produce	json
//	@Param		id	path	string	true	"Clip ID"
//	@Success	200	{object}	map[string]interface{}
//	@Failure	404	{object}	map[string]interface{}
//	@Router		/clip/{id} [get]
func (c *Controller) ClipView(g *gin.Context) {
	id := g.Param("id")
	video := &model.Video{ID: id}
	result := c.datastore.Preload("Actors.Categories").Preload("Categories").First(video)
	if result.RowsAffected < 1 {
		g.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	if video.Type != "c" {
		g.JSON(http.StatusNotFound, gin.H{"error": "not a clip"})
		return
	}
	type clipID struct{ ID string }
	var clipIDs []clipID
	c.datastore.Model(&model.Video{}).Where("type = ?", "c").Find(&clipIDs)
	var clipList []string
	for _, cid := range clipIDs {
		if cid.ID != id {
			clipList = append(clipList, cid.ID)
		}
	}
	for i := range clipList {
		j := rand.Intn(i + 1)
		clipList[i], clipList[j] = clipList[j], clipList[i]
	}
	clipList = append([]string{id}, clipList...)
	g.JSON(http.StatusOK, gin.H{
		"video":    video,
		"clip_ids": clipList,
	})
}
