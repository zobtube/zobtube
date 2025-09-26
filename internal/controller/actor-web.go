package controller

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"

	"github.com/zobtube/zobtube/internal/model"
)

func (c *Controller) ActorList(g *gin.Context) {
	var actors []model.Actor
	c.datastore.Find(&actors).Limit(30).Offset(0).Order("created_at")
	c.HTML(g, http.StatusOK, "actor/list.html", gin.H{
		"Actors": actors,
	})
}

type ActorNewForm struct {
	Name    string `form:"name"`
	SexEnum string `form:"sex"`
}

func (c *Controller) ActorNew(g *gin.Context) {
	var err error
	if g.Request.Method == "POST" {
		var form ActorNewForm
		err = g.ShouldBind(&form)
		if err == nil {
			actor := &model.Actor{
				Name: form.Name,
				Sex:  form.SexEnum,
			}
			err = c.datastore.Create(&actor).Error
			if err == nil {
				g.Redirect(http.StatusFound, "/actor/"+actor.ID+"/edit")
				return
			}
		}
	}

	c.HTML(g, http.StatusOK, "actor/create.html", gin.H{
		"Error": err,
	})
}

func (c *Controller) ActorView(g *gin.Context) {
	// get id from path
	id := g.Param("id")

	// get item from ID
	actor := &model.Actor{
		ID: id,
	}
	result := c.datastore.Preload(clause.Associations).First(actor)

	// check result
	if result.RowsAffected < 1 {
		//TODO: return to homepage
		g.JSON(404, gin.H{})
		return
	}

	c.HTML(g, http.StatusOK, "actor/view.html", gin.H{
		"Actor": actor,
	})
}

func (c *Controller) ActorEdit(g *gin.Context) {
	// get id from path
	id := g.Param("id")

	// get item from ID
	actor := &model.Actor{
		ID: id,
	}
	result := c.datastore.Preload("Links").Preload("Aliases").Preload("Categories").First(actor)

	// check result
	if result.RowsAffected < 1 {
		//TODO: return to homepage
		g.JSON(404, gin.H{})
		return
	}

	// get categories
	categories := []model.Category{}
	err := c.datastore.Preload("Sub").Find(&categories).Error
	if err != nil {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	// loading configuration from database
	dbconfig := &model.Configuration{}
	result = c.datastore.First(dbconfig)

	// check result
	if result.RowsAffected < 1 {
		g.JSON(500, gin.H{
			"error": "configuration not found, restarting the appliaction should fix the issue",
		})
		return
	}

	// get all providers
	var providers []model.Provider
	c.datastore.Find(&providers)

	c.HTML(g, http.StatusOK, "actor/edit.html", gin.H{
		"Actor":       actor,
		"Providers":   providers,
		"Categories":  categories,
		"OfflineMode": dbconfig.OfflineMode,
	})
}

func (c *Controller) ActorThumb(g *gin.Context) {
	// get id from path
	id := g.Param("id")

	// get item from ID
	actor := &model.Actor{
		ID: id,
	}
	result := c.datastore.First(actor)

	// check result
	if result.RowsAffected < 1 {
		g.JSON(404, gin.H{})
		return
	}

	// check if thumbnail exists
	if !actor.Thumbnail {
		g.Redirect(http.StatusFound, ACTOR_PROFILE_PICTURE_MISSING)
		return
	}

	// construct file path
	targetPath := filepath.Join(c.config.Media.Path, ACTOR_FILEPATH, id, "thumb.jpg")

	// give file path
	g.File(targetPath)
}

func (c *Controller) ActorDelete(g *gin.Context) {
	// get id from path
	id := g.Param("id")

	// get item from ID
	actor := &model.Actor{
		ID: id,
	}
	result := c.datastore.First(actor)

	// check result
	if result.RowsAffected < 1 {
		g.JSON(404, gin.H{})
		return
	}

	// delete thumb
	thumbPath := filepath.Join(c.config.Media.Path, ACTOR_FILEPATH, id, "thumb.jpg")
	_, err := os.Stat(thumbPath)
	if err != nil && !os.IsNotExist(err) {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	if !os.IsNotExist(err) {
		// exist, deleting it
		err = os.Remove(thumbPath)
		if err != nil {
			g.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	// delete folder
	folderPath := filepath.Join(c.config.Media.Path, ACTOR_FILEPATH, id)
	_, err = os.Stat(folderPath)
	if err != nil && !os.IsNotExist(err) {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	if !os.IsNotExist(err) {
		// exist, deleting it
		err = os.Remove(folderPath)
		if err != nil {
			g.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	// delete object
	err = c.datastore.Delete(actor).Error
	if err != nil {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	// all good
	g.Redirect(http.StatusFound, "/actors")
}
