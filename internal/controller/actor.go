package controller

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"

	"github.com/zobtube/zobtube/internal/model"
)

const errHumanProviderNotFound = "Unable to retrieve provider"

// ActorList godoc
//
//	@Summary	List all actors
//	@Tags		actor
//	@Produce	json
//	@Success	200	{object}	map[string]model.Actor[]
//	@Router		/actor [get]
func (c *Controller) ActorList(g *gin.Context) {
	var actors []model.Actor
	c.datastore.Preload("Videos").Preload("Links").Order("name").Find(&actors)
	g.JSON(http.StatusOK, gin.H{
		"items": actors,
		"total": len(actors),
	})
}

// ActorGet godoc
//
//	@Summary	Get actor by ID
//	@Tags		actor
//	@Produce	json
//	@Param		id	path		string	true	"Actor ID"
//	@Success	200	{object}	model.Actor
//	@Failure	404
//	@Router		/actor/{id} [get]
func (c *Controller) ActorGet(g *gin.Context) {
	id := g.Param("id")
	actor := &model.Actor{ID: id}
	result := c.datastore.Preload(clause.Associations).First(actor)
	if result.RowsAffected < 1 {
		g.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	g.JSON(http.StatusOK, actor)
}

// ActorDelete godoc
//
//	@Summary	Delete an actor
//	@Tags		actor
//	@Param		id	path	string	true	"Actor ID"
//	@Success	204
//	@Failure	404
//	@Router		/actor/{id} [delete]
func (c *Controller) ActorDelete(g *gin.Context) {
	id := g.Param("id")
	actor := &model.Actor{ID: id}
	result := c.datastore.First(actor)
	if result.RowsAffected < 1 {
		g.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	thumbPath := filepath.Join(c.config.Media.Path, ACTOR_FILEPATH, actor.ID, "thumb.jpg")
	if _, err := os.Stat(thumbPath); err == nil {
		_ = os.Remove(thumbPath)
	}
	folderPath := filepath.Join(c.config.Media.Path, ACTOR_FILEPATH, actor.ID)
	if _, err := os.Stat(folderPath); err == nil {
		_ = os.Remove(folderPath)
	}
	if err := c.datastore.Delete(actor).Error; err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	g.JSON(http.StatusNoContent, gin.H{})
}

// ActorThumb godoc
//
//	@Summary	Get actor thumbnail image
//	@Tags		actor
//	@Param		id	path	string	true	"Actor ID"
//	@Success	200	file	bytes
//	@Failure	404
//	@Router		/actor/{id}/thumb [get]
func (c *Controller) ActorThumb(g *gin.Context) {
	id := g.Param("id")
	actor := &model.Actor{ID: id}
	result := c.datastore.First(actor)
	if result.RowsAffected < 1 {
		g.JSON(404, gin.H{})
		return
	}
	if !actor.Thumbnail {
		g.Redirect(http.StatusFound, ACTOR_PROFILE_PICTURE_MISSING)
		return
	}
	targetPath := filepath.Join(c.config.Media.Path, ACTOR_FILEPATH, id, "thumb.jpg")
	g.File(targetPath)
}

// ActorNew godoc
//
//	@Summary	Create a new actor
//	@Tags		actor
//	@Accept		multipart/form-data
//	@Param		id		formData	string	true	"Actor ID"
//	@Param		name	formData	string	true	"Actor name"
//	@Param		sex		formData	string	false	"Sex (m/f)"
//	@Success	200	{object}	map[string]interface{}
//	@Failure	500	{object}	map[string]interface{}
//	@Router		/actor [post]
func (c *Controller) ActorNew(g *gin.Context) {
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

// ActorProviderSearch godoc
//
//	@Summary	Search actor in external provider and create link
//	@Tags		actor
//	@Param		id	path	string	true	"Actor ID"
//	@Param		provider_slug	path	string	true	"Provider slug (babesdirectory, babepedia, etc.)"
//	@Success	200	{object}	map[string]interface{}
//	@Failure	404	{object}	map[string]interface{}
//	@Router		/actor/{id}/provider/{provider_slug} [get]
func (c *Controller) ActorProviderSearch(g *gin.Context) {
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

// ActorLinkThumbGet godoc
//
//	@Summary	Get thumbnail from actor link (provider)
//	@Tags		actor
//	@Param		id	path	string	true	"Link ID"
//	@Success	200	file	bytes
//	@Failure	404	{object}	map[string]interface{}
//	@Router		/actor/link/{id}/thumb [get]
func (c *Controller) ActorLinkThumbGet(g *gin.Context) {
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

// ActorLinkThumbDelete godoc
//
//	@Summary	Delete actor link
//	@Tags		actor
//	@Param		id	path	string	true	"Link ID"
//	@Success	200
//	@Failure	404	{object}	map[string]interface{}
//	@Router		/actor/link/{id} [delete]
func (c *Controller) ActorLinkThumbDelete(g *gin.Context) {
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

// ActorUploadThumb godoc
//
//	@Summary	Upload actor thumbnail
//	@Tags		actor
//	@Accept		multipart/form-data
//	@Param		id	path	string	true	"Actor ID"
//	@Param		pp	formData	file	true	"Thumbnail image"
//	@Success	200
//	@Failure	404	{object}	map[string]interface{}
//	@Router		/actor/{id}/thumb [post]
func (c *Controller) ActorUploadThumb(g *gin.Context) {
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

// ActorLinkCreate godoc
//
//	@Summary	Create actor link to provider profile
//	@Tags		actor
//	@Accept		x-www-form-urlencoded
//	@Param		id	path	string	true	"Actor ID"
//	@Param		url	formData	string	true	"Provider profile URL"
//	@Param		provider	formData	string	true	"Provider slug"
//	@Success	200	{object}	map[string]interface{}
//	@Failure	404	{object}	map[string]interface{}
//	@Router		/actor/{id}/link [post]
func (c *Controller) ActorLinkCreate(g *gin.Context) {
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

// ActorAliasCreate godoc
//
//	@Summary	Add alias to actor
//	@Tags		actor
//	@Accept		x-www-form-urlencoded
//	@Param		id	path	string	true	"Actor ID"
//	@Param		alias	formData	string	true	"Alias name"
//	@Success	200	{object}	map[string]interface{}
//	@Failure	400	{object}	map[string]interface{}
//	@Router		/actor/{id}/alias [post]
func (c *Controller) ActorAliasCreate(g *gin.Context) {
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

// ActorAliasRemove godoc
//
//	@Summary	Remove actor alias
//	@Tags		actor
//	@Param		id	path	string	true	"Alias ID"
//	@Success	200
//	@Failure	404	{object}	map[string]interface{}
//	@Router		/actor/alias/{id} [delete]
func (c *Controller) ActorAliasRemove(g *gin.Context) {
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

// ActorCategories godoc
//
//	@Summary	Add or remove category from actor (PUT=add, DELETE=remove)
//	@Tags		actor
//	@Param		id	path	string	true	"Actor ID"
//	@Param		category_id	path	string	true	"Category (sub) ID"
//	@Success	200
//	@Failure	404	{object}	map[string]interface{}
//	@Router		/actor/{id}/category/{category_id} [put]
//	@Router		/actor/{id}/category/{category_id} [delete]
func (c *Controller) ActorCategories(g *gin.Context) {
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

// ActorRename godoc
//
//	@Summary	Rename actor
//	@Tags		actor
//	@Accept		x-www-form-urlencoded
//	@Param		id	path	string	true	"Actor ID"
//	@Param		name	formData	string	true	"New name"
//	@Success	200
//	@Failure	400	{object}	map[string]interface{}
//	@Failure	404	{object}	map[string]interface{}
//	@Router		/actor/{id}/rename [post]
func (c *Controller) ActorRename(g *gin.Context) {
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

// ActorDescription godoc
//
//	@Summary	Set actor description
//	@Tags		actor
//	@Accept		x-www-form-urlencoded
//	@Param		id	path	string	true	"Actor ID"
//	@Param		description	formData	string	false	"Description text"
//	@Success	200
//	@Failure	400	{object}	map[string]interface{}
//	@Failure	404	{object}	map[string]interface{}
//	@Router		/actor/{id}/description [post]
func (c *Controller) ActorDescription(g *gin.Context) {
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

// ActorMerge godoc
//
//	@Summary	Merge source actor into target actor
//	@Tags		actor
//	@Accept		json
//	@Param		id	path	string	true	"Source actor ID"
//	@Param		body	body	object	true	"JSON with target_id"
//	@Success	200	{object}	map[string]interface{}
//	@Failure	400	{object}	map[string]interface{}
//	@Failure	404	{object}	map[string]interface{}
//	@Router		/actor/{id}/merge [post]
func (c *Controller) ActorMerge(g *gin.Context) {
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
