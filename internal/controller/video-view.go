package controller

import (
	"github.com/gin-gonic/gin"

	"gitlab.com/zobtube/zobtube/internal/model"
)

func (c *Controller) VideoViewAjaxIncrement(g *gin.Context) {
	// get id from path
	id := g.Param("id")

	// get item from ID
	video := &model.Video{
		ID: id,
	}
	result := c.datastore.First(video)

	// check result
	if result.RowsAffected < 1 {
		g.JSON(404, gin.H{})
		return
	}

	user := g.MustGet("user").(*model.User)

	count := &model.VideoView{}
	result = c.datastore.Debug().First(&count, "video_id = ? AND user_id = ?", video.ID, user.ID)

	// check result
	if result.RowsAffected > 0 {
		// already exists, increment
		count.Count++
		err := c.datastore.Save(count).Error
		if err != nil {
			g.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		g.JSON(200, gin.H{"view-count": count.Count})
		return
	}

	// new view, create item
	count = &model.VideoView{
		User:  *user,
		Video: *video,
		Count: 1,
	}

	err := c.datastore.Debug().Create(count).Error
	if err != nil {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	g.JSON(200, gin.H{"view-count": count.Count})
}
