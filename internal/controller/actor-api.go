package controller

import (
	"path/filepath"

	"github.com/gin-gonic/gin"

	"github.com/zobtube/zobtube/internal/model"
)

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

func (c *Controller) ActorAjaxAliasCreate(g *gin.Context) {
	var err error

	form := struct {
		Alias string `form:"alias"`
	}{}
	err = g.ShouldBind(&form)
	if err != nil {
		g.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	// get actor id from path
	actorID := g.Param("id")

	alias := model.ActorAlias{
		Name:    form.Alias,
		ActorID: actorID,
	}

	err = c.datastore.Create(&alias).Error
	if err != nil {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	g.JSON(200, gin.H{
		"id": alias.ID,
	})
}

func (c *Controller) ActorAjaxAliasRemove(g *gin.Context) {
	var err error

	// get alias id from path
	aliasID := g.Param("id")

	alias := model.ActorAlias{
		ID: aliasID,
	}
	result := c.datastore.First(&alias)

	// check result
	if result.RowsAffected < 1 {
		g.JSON(404, gin.H{})
		return
	}

	err = c.datastore.Delete(&alias).Error
	if err != nil {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	g.JSON(200, gin.H{})
}

func (c *Controller) ActorAjaxCategories(g *gin.Context) {
	// get id from path
	id := g.Param("id")
	category_id := g.Param("category_id")

	// get item from ID
	actor := &model.Actor{
		ID: id,
	}
	result := c.datastore.First(actor)

	// check result
	if result.RowsAffected < 1 {
		g.JSON(404, gin.H{
			"error": "actor not found",
		})
		return
	}

	subCategory := &model.CategorySub{
		ID: category_id,
	}
	result = c.datastore.First(&subCategory)

	// check result
	if result.RowsAffected < 1 {
		g.JSON(404, gin.H{
			"error": "sub-category not found",
		})
		return
	}

	var res error
	if g.Request.Method == "PUT" {
		res = c.datastore.Model(actor).Association("Categories").Append(subCategory)
	} else {
		res = c.datastore.Model(actor).Association("Categories").Delete(subCategory)
	}

	if res != nil {
		g.JSON(500, gin.H{
			"error": res.Error(),
		})
		return
	}

	g.JSON(200, gin.H{})
}

func (c *Controller) ActorAjaxRename(g *gin.Context) {
	var err error

	form := struct {
		Name string `form:"name"`
	}{}
	err = g.ShouldBind(&form)
	if err != nil {
		g.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	if form.Name == "" {
		g.JSON(400, gin.H{
			"error": "actor name cannot be empty",
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

	actor.Name = form.Name
	err = c.datastore.Save(actor).Error
	if err != nil {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	g.JSON(200, gin.H{})
}
