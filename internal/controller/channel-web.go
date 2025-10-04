package controller

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"

	"github.com/zobtube/zobtube/internal/model"
)

func (c *Controller) ChannelList(g *gin.Context) {
	var channels []model.Channel
	c.datastore.Find(&channels)
	c.HTML(g, http.StatusOK, "channel/list.html", gin.H{
		"Channels": channels,
	})
}

type ChannelNewForm struct {
	Name string `form:"name"`
}

func (c *Controller) ChannelCreate(g *gin.Context) {
	var err error
	if g.Request.Method == "POST" {
		var form ActorNewForm
		err = g.ShouldBind(&form)
		if err == nil {
			channel := &model.Channel{
				Name: form.Name,
			}
			err = c.datastore.Create(&channel).Error
			if err == nil {
				g.Redirect(http.StatusFound, "/channel/"+channel.ID)
				return
			}
		}
	}
	c.HTML(g, http.StatusOK, "channel/create.html", gin.H{
		"Error": err,
	})
}

func (c *Controller) ChannelView(g *gin.Context) {
	// get id from path
	id := g.Param("id")

	// get item from ID
	channel := &model.Channel{
		ID: id,
	}
	result := c.datastore.Preload(clause.Associations).First(channel)

	// check result
	if result.RowsAffected < 1 {
		c.ErrNotFound(g)
		return
	}

	var videos []model.Video
	c.datastore.Where("channel_id = ?", channel.ID).Find(&videos)

	c.HTML(g, http.StatusOK, "channel/view.html", gin.H{
		"Channel": channel,
		"Videos":  videos,
	})
}

func (c *Controller) ChannelThumb(g *gin.Context) {
	// get id from path
	id := g.Param("id")

	// get item from ID
	channel := &model.Channel{
		ID: id,
	}
	result := c.datastore.First(channel)

	// check result
	if result.RowsAffected < 1 {
		c.ErrNotFound(g)
		return
	}

	// check if thumbnail exists
	if !channel.Thumbnail {
		g.Redirect(http.StatusFound, ACTOR_PROFILE_PICTURE_MISSING)
		return
	}

	// construct file path
	targetPath := filepath.Join(c.config.Media.Path, CHANNEL_FILEPATH, id, "thumb.jpg")

	// give file path
	g.File(targetPath)
}

func (c *Controller) ChannelEdit(g *gin.Context) {
	// get id from path
	id := g.Param("id")

	// get item from ID
	channel := &model.Channel{
		ID: id,
	}
	result := c.datastore.First(channel)

	// check result
	if result.RowsAffected < 1 {
		c.ErrNotFound(g)
		return
	}

	// construct file path
	targetPath := filepath.Join(c.config.Media.Path, CHANNEL_FILEPATH, id, "thumb.jpg")

	if g.Request.Method == "POST" {
		file, err := g.FormFile("profile")
		if err != nil {
			c.ErrFatal(g, err.Error())
			return
		}

		//save thumb on disk
		err = g.SaveUploadedFile(file, targetPath)
		if err != nil {
			c.ErrFatal(g, err.Error())
			return
		}

		channel.Thumbnail = true
		err = c.datastore.Save(channel).Error
		if err != nil {
			c.ErrFatal(g, err.Error())
			return
		}
	}

	c.HTML(g, http.StatusOK, "channel/edit.html", gin.H{
		"Channel": channel,
	})
}
