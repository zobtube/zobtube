package controller

import (
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"github.com/zobtube/zobtube/internal/model"
)

const errHumanProviderNotFound = "Unable to retrieve provider"

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
			"error_human": errHumanProviderNotFound,
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

	url, err := provider.ActorSearch(dbconfig.OfflineMode, actor.Name)
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
			"error_human": errHumanProviderNotFound,
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

	thumb, err := provider.ActorGetThumb(dbconfig.OfflineMode, link.Actor.Name, link.URL)
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

	// save thumb on disk
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
			"error_human": errHumanProviderNotFound,
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

func (c *Controller) ActorAjaxDescription(g *gin.Context) {
	var err error

	form := struct {
		Description string `form:"description"`
	}{}
	err = g.ShouldBind(&form)
	if err != nil {
		g.JSON(400, gin.H{
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

	actor.Description = form.Description
	err = c.datastore.Save(actor).Error
	if err != nil {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	g.JSON(200, gin.H{})
}

func (c *Controller) ActorAjaxMerge(g *gin.Context) {
	sourceID := g.Param("id")
	form := struct {
		TargetID string `json:"target_id"`
	}{}
	if err := g.ShouldBindJSON(&form); err != nil {
		g.JSON(400, gin.H{"error": err.Error()})
		return
	}
	targetID := form.TargetID
	if targetID == "" {
		g.JSON(400, gin.H{"error": "target_id is required"})
		return
	}
	if sourceID == targetID {
		g.JSON(400, gin.H{"error": "source and target must be different"})
		return
	}

	source := &model.Actor{ID: sourceID}
	if res := c.datastore.Preload("Videos").Preload("Aliases").Preload("Links").Preload("Categories").First(source); res.RowsAffected < 1 {
		g.JSON(404, gin.H{"error": "source actor not found"})
		return
	}
	target := &model.Actor{ID: targetID}
	if res := c.datastore.Preload("Aliases").Preload("Links").First(target); res.RowsAffected < 1 {
		g.JSON(404, gin.H{"error": "target actor not found"})
		return
	}

	tx := c.datastore.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for i := range source.Videos {
		video := &source.Videos[i]
		if err := tx.Model(video).Association("Actors").Delete(source); err != nil {
			tx.Rollback()
			g.JSON(500, gin.H{"error": err.Error()})
			return
		}
		if err := tx.Model(video).Association("Actors").Append(target); err != nil {
			tx.Rollback()
			g.JSON(500, gin.H{"error": err.Error()})
			return
		}
	}

	targetAliasNames := make(map[string]struct{})
	for _, a := range target.Aliases {
		targetAliasNames[a.Name] = struct{}{}
	}
	for i := range source.Aliases {
		a := &source.Aliases[i]
		if _, exists := targetAliasNames[a.Name]; exists {
			if err := tx.Delete(a).Error; err != nil {
				tx.Rollback()
				g.JSON(500, gin.H{"error": err.Error()})
				return
			}
		} else {
			a.ActorID = target.ID
			if err := tx.Save(a).Error; err != nil {
				tx.Rollback()
				g.JSON(500, gin.H{"error": err.Error()})
				return
			}
			targetAliasNames[a.Name] = struct{}{}
		}
	}

	type linkKey struct{ Provider, URL string }
	targetLinks := make(map[linkKey]struct{})
	for _, l := range target.Links {
		targetLinks[linkKey{l.Provider, l.URL}] = struct{}{}
	}
	for i := range source.Links {
		l := &source.Links[i]
		k := linkKey{l.Provider, l.URL}
		if _, exists := targetLinks[k]; exists {
			if err := tx.Delete(l).Error; err != nil {
				tx.Rollback()
				g.JSON(500, gin.H{"error": err.Error()})
				return
			}
		} else {
			l.ActorID = target.ID
			if err := tx.Save(l).Error; err != nil {
				tx.Rollback()
				g.JSON(500, gin.H{"error": err.Error()})
				return
			}
			targetLinks[k] = struct{}{}
		}
	}

	for i := range source.Categories {
		cat := &source.Categories[i]
		if err := tx.Model(target).Association("Categories").Append(cat); err != nil {
			tx.Rollback()
			g.JSON(500, gin.H{"error": err.Error()})
			return
		}
	}
	if err := tx.Model(source).Association("Categories").Clear(); err != nil {
		tx.Rollback()
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if err := tx.Delete(source).Error; err != nil {
		tx.Rollback()
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if err := tx.Commit().Error; err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}

	thumbPath := filepath.Join(c.config.Media.Path, ACTOR_FILEPATH, source.ID, "thumb.jpg")
	if _, err := os.Stat(thumbPath); err == nil {
		_ = os.Remove(thumbPath)
	}
	folderPath := filepath.Join(c.config.Media.Path, ACTOR_FILEPATH, source.ID)
	if _, err := os.Stat(folderPath); err == nil {
		_ = os.RemoveAll(folderPath)
	}

	g.JSON(200, gin.H{"redirect": "/actor/" + target.ID + "/edit"})
}
