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
	g.HTML(http.StatusOK, "actor/list.html", gin.H{
		"Actors": actors,
		"User":   g.MustGet("user").(*model.User),
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
				g.Redirect(http.StatusFound, "/actor/"+actor.ID)
				return
			}
		}
	}
	g.HTML(http.StatusOK, "actor/create.html", gin.H{
		"User":  g.MustGet("user").(*model.User),
		"Error": err,
	})
}

func (c *Controller) ActorAjaxNew(g *gin.Context) {
	var err error
	form := struct {
		ID      string `form:"id"`
		Name    string `form:"name"`
		SexEnum string `form:"sex"`
	}{}
	err = g.ShouldBind(&form)
	if err != nil {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	actor := &model.Actor{
		ID:   form.ID,
		Name: form.Name,
		Sex:  form.SexEnum,
	}
	err = c.datastore.Create(&actor).Error
	if err != nil {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	g.JSON(200, gin.H{
		"result": actor.ID,
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

	g.HTML(http.StatusOK, "actor/view.html", gin.H{
		"User":  g.MustGet("user").(*model.User),
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
	result := c.datastore.Preload("Links").First(actor)

	// check result
	if result.RowsAffected < 1 {
		//TODO: return to homepage
		g.JSON(404, gin.H{})
		return
	}

	g.HTML(http.StatusOK, "actor/edit.html", gin.H{
		"User":      g.MustGet("user").(*model.User),
		"Actor":     actor,
		"Providers": c.providers,
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

func (c *Controller) ActorAjaxProviderSearch(g *gin.Context) {
	// get actor id from path
	id := g.Param("id")

	// get actor from ID
	actor := &model.Actor{
		ID: id,
	}
	result := c.datastore.First(actor)

	// check result
	if result.RowsAffected < 1 {
		g.JSON(404, gin.H{})
		return
	}

	// get provider slug from path
	provider_slug := g.Param("provider_slug")
	provider, err := c.ProviderGet(provider_slug)
	if err != nil {
		g.JSON(404, gin.H{
			"error":       err.Error(),
			"error_human": "Unable to retrieve provider",
		})
		return
	}

	url, err := provider.ActorSearch(actor.Name)
	if err != nil {
		g.JSON(404, gin.H{
			"error":       err.Error(),
			"error_human": "Provider did not found a result",
		})
		return
	}

	// url found, storing it
	link := &model.ActorLink{
		Actor:    *actor,
		Provider: provider_slug,
		URL:      url,
	}

	err = c.datastore.Create(link).Error
	if err != nil {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	g.JSON(200, gin.H{
		"link_id":  link.ID,
		"link_url": url,
	})
}

func (c *Controller) ActorAjaxLinkThumbGet(g *gin.Context) {
	// get actor id from path
	id := g.Param("id")

	// get actor from ID
	link := &model.ActorLink{
		ID: id,
	}
	result := c.datastore.Preload("Actor").First(link)

	// check result
	if result.RowsAffected < 1 {
		g.JSON(404, gin.H{})
		return
	}

	// get provider slug from path
	provider, err := c.ProviderGet(link.Provider)
	if err != nil {
		g.JSON(404, gin.H{
			"error":       err.Error(),
			"error_human": "Unable to retrieve provider",
		})
		return
	}

	thumb, err := provider.ActorGetThumb(link.Actor.Name, link.URL)
	if err != nil {
		g.JSON(404, gin.H{
			"error":       err.Error(),
			"error_human": "Provider did not found a result",
		})
		return
	}

	g.Data(200, "image/png", thumb)
}

func (c *Controller) ActorAjaxLinkThumbDelete(g *gin.Context) {
	// get actor id from path
	id := g.Param("id")

	// get actor from ID
	link := &model.ActorLink{
		ID: id,
	}
	result := c.datastore.Preload("Actor").First(link)

	// check result
	if result.RowsAffected < 1 {
		g.JSON(404, gin.H{})
		return
	}

	err := c.datastore.Delete(&link).Error
	if err != nil {
		g.JSON(500, gin.H{
			"error":       err.Error(),
			"human_error": "unable to delete actor link",
		})
		return
	}

	g.JSON(200, gin.H{})
}

func (c *Controller) ActorAjaxThumb(g *gin.Context) {
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

	file, err := g.FormFile("pp")
	if err != nil {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	// construct file path
	targetPath := filepath.Join(c.config.Media.Path, ACTOR_FILEPATH, id, "thumb.jpg")

	//save thumb on disk
	err = g.SaveUploadedFile(file, targetPath)
	if err != nil {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	// check if thumbnail exists
	if !actor.Thumbnail {
		actor.Thumbnail = true
		err = c.datastore.Save(actor).Error
		if err != nil {
			g.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	// all good
	g.JSON(200, gin.H{})
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

func (c *Controller) ActorAjaxLinkCreate(g *gin.Context) {
	var err error

	form := struct {
		URL      string `form:"url"`
		Provider string `form:"provider"`
	}{}
	err = g.ShouldBind(&form)
	if err != nil {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	// get actor id from path
	id := g.Param("id")

	// get actor from ID
	actor := &model.Actor{
		ID: id,
	}
	result := c.datastore.First(actor)

	// check result
	if result.RowsAffected < 1 {
		g.JSON(404, gin.H{})
		return
	}

	// get provider slug from path
	_, err = c.ProviderGet(form.Provider)
	if err != nil {
		g.JSON(404, gin.H{
			"error":       err.Error(),
			"error_human": "Unable to retrieve provider",
		})
		return
	}

	// url found, storing it
	link := &model.ActorLink{
		Actor:    *actor,
		Provider: form.Provider,
		URL:      form.URL,
	}

	err = c.datastore.Create(link).Error
	if err != nil {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	g.JSON(200, gin.H{
		"link_id":  link.ID,
		"link_url": link.URL,
	})
}
