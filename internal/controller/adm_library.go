package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/zobtube/zobtube/internal/model"
)

// AdmLibraryList godoc
//
//	@Summary	List all libraries (admin)
//	@Tags		admin
//	@Produce	json
//	@Success	200	{object}	map[string]interface{}
//	@Router		/adm/libraries [get]
func (c *Controller) AdmLibraryList(g *gin.Context) {
	var libs []model.Library
	c.datastore.Order("created_at").Find(&libs)
	g.JSON(http.StatusOK, gin.H{"items": libs, "total": len(libs)})
}

// AdmLibraryCreate godoc
//
//	@Summary	Create a library (admin)
//	@Tags		admin
//	@Accept		json
//	@Param		body	body	object	true	"JSON with name, type, config"
//	@Success	201	{object}	map[string]interface{}
//	@Failure	400	{object}	map[string]interface{}
//	@Router		/adm/libraries [post]
func (c *Controller) AdmLibraryCreate(g *gin.Context) {
	var body struct {
		Name   string                 `json:"name" binding:"required"`
		Type   model.LibraryType      `json:"type" binding:"required"`
		Config model.LibraryConfig    `json:"config"`
		Default bool                  `json:"default"`
	}
	if err := g.ShouldBindJSON(&body); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	lib := model.Library{
		ID:        uuid.NewString(),
		Name:      body.Name,
		Type:      body.Type,
		Config:    body.Config,
		IsDefault: body.Default,
	}
	if body.Type == model.LibraryTypeFilesystem && body.Config.Filesystem != nil && body.Config.Filesystem.Path != "" {
		// no-op
	} else if body.Type == model.LibraryTypeS3 && body.Config.S3 != nil && body.Config.S3.Bucket != "" {
		// no-op
	} else {
		g.JSON(http.StatusBadRequest, gin.H{"error": "invalid config for library type"})
		return
	}
	if err := c.datastore.Create(&lib).Error; err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if body.Default {
		c.datastore.Model(&model.Library{}).Where("id != ?", lib.ID).Update("is_default", false)
	}
	if c.storageResolver != nil {
		c.storageResolver.Invalidate(lib.ID)
	}
	g.JSON(http.StatusCreated, gin.H{"id": lib.ID, "library": lib})
}

// AdmLibraryUpdate godoc
//
//	@Summary	Update a library (admin)
//	@Tags		admin
//	@Accept		json
//	@Param		id	path	string	true	"Library ID"
//	@Param		body	body	object	true	"JSON with name, config, default"
//	@Success	200	{object}	map[string]interface{}
//	@Failure	400	{object}	map[string]interface{}
//	@Router		/adm/libraries/{id} [put]
func (c *Controller) AdmLibraryUpdate(g *gin.Context) {
	id := g.Param("id")
	var lib model.Library
	if err := c.datastore.First(&lib, "id = ?", id).Error; err != nil {
		g.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	var body struct {
		Name    *string                `json:"name"`
		Type    *model.LibraryType     `json:"type"`
		Config  *model.LibraryConfig   `json:"config"`
		Default *bool                  `json:"default"`
	}
	if err := g.ShouldBindJSON(&body); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if body.Name != nil {
		lib.Name = *body.Name
	}
	if body.Type != nil {
		t := *body.Type
		if t != model.LibraryTypeFilesystem && t != model.LibraryTypeS3 {
			g.JSON(http.StatusBadRequest, gin.H{"error": "invalid library type"})
			return
		}
		lib.Type = t
	}
	if body.Config != nil {
		newConfig := *body.Config
		// Preserve S3 secret when edit form omits it (empty = keep existing)
		if newConfig.S3 != nil && lib.Config.S3 != nil && newConfig.S3.SecretAccessKey == "" && lib.Config.S3.SecretAccessKey != "" {
			newConfig.S3.SecretAccessKey = lib.Config.S3.SecretAccessKey
		}
		lib.Config = newConfig
	}
	// Validate config matches type
	if lib.Type == model.LibraryTypeFilesystem && (lib.Config.Filesystem == nil || lib.Config.Filesystem.Path == "") {
		g.JSON(http.StatusBadRequest, gin.H{"error": "filesystem library requires config.filesystem.path"})
		return
	}
	if lib.Type == model.LibraryTypeS3 && (lib.Config.S3 == nil || lib.Config.S3.Bucket == "") {
		g.JSON(http.StatusBadRequest, gin.H{"error": "s3 library requires config.s3 with bucket"})
		return
	}
	if body.Default != nil {
		lib.IsDefault = *body.Default
		if lib.IsDefault {
			c.datastore.Model(&model.Library{}).Where("id != ?", lib.ID).Update("is_default", false)
		}
	}
	if err := c.datastore.Save(&lib).Error; err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if c.storageResolver != nil {
		c.storageResolver.Invalidate(lib.ID)
	}
	g.JSON(http.StatusOK, gin.H{"library": lib})
}

// AdmLibraryDelete godoc
//
//	@Summary	Delete a library (admin)
//	@Tags		admin
//	@Param		id	path	string	true	"Library ID"
//	@Success	204
//	@Failure	400	{object}	map[string]interface{}
//	@Router		/adm/libraries/{id} [delete]
func (c *Controller) AdmLibraryDelete(g *gin.Context) {
	id := g.Param("id")
	var lib model.Library
	if err := c.datastore.First(&lib, "id = ?", id).Error; err != nil {
		g.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	if lib.IsDefault {
		g.JSON(http.StatusBadRequest, gin.H{"error": "cannot delete the default library"})
		return
	}
	var videoCount int64
	c.datastore.Model(&model.Video{}).Where("library_id = ?", id).Count(&videoCount)
	if videoCount > 0 {
		g.JSON(http.StatusBadRequest, gin.H{"error": "library has videos, move or delete them first"})
		return
	}
	if err := c.datastore.Delete(&lib).Error; err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if c.storageResolver != nil {
		c.storageResolver.Invalidate(id)
	}
	g.Status(http.StatusNoContent)
}
