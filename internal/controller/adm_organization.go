package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/zobtube/zobtube/internal/model"
)

// AdmOrganizationList godoc
//
//	@Summary	List all organizations (admin)
//	@Tags		admin
//	@Produce	json
//	@Success	200	{object}	map[string]interface{}
//	@Router		/adm/organizations [get]
func (c *Controller) AdmOrganizationList(g *gin.Context) {
	var orgs []model.Organization
	c.datastore.Order("created_at").Find(&orgs)

	type withCount struct {
		model.Organization
		VideoCount int64 `json:"video_count"`
	}
	items := make([]withCount, 0, len(orgs))
	for i := range orgs {
		var count int64
		c.datastore.Model(&model.Video{}).Where("organization_id = ?", orgs[i].ID).Count(&count)
		items = append(items, withCount{Organization: orgs[i], VideoCount: count})
	}

	var cfg model.Configuration
	_ = c.datastore.First(&cfg).Error
	g.JSON(http.StatusOK, gin.H{
		"items":                items,
		"total":                len(items),
		"reorganize_on_import": cfg.ReorganizeOnImport,
	})
}

// AdmOrganizationCreate godoc
//
//	@Summary	Create an organization (admin)
//	@Tags		admin
//	@Accept		json
//	@Param		body	body	object	true	"JSON with name, template, active"
//	@Success	201	{object}	map[string]interface{}
//	@Failure	400	{object}	map[string]interface{}
//	@Router		/adm/organizations [post]
func (c *Controller) AdmOrganizationCreate(g *gin.Context) {
	var body struct {
		Name     string `json:"name" binding:"required"`
		Template string `json:"template" binding:"required"`
		Scope    string `json:"scope"`
		Active   bool   `json:"active"`
	}
	if err := g.ShouldBindJSON(&body); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := model.ValidateOrganizationTemplate(body.Template); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	scope := body.Scope
	if scope == "" {
		scope = model.OrganizationScopeVideo
	}
	org := model.Organization{
		ID:       uuid.NewString(),
		Name:     body.Name,
		Template: body.Template,
		Scope:    scope,
		Active:   body.Active,
	}
	if err := c.datastore.Create(&org).Error; err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if org.Active {
		scope := org.Scope
		if scope == "" {
			scope = model.OrganizationScopeVideo
		}
		c.datastore.Model(&model.Organization{}).Where("id != ? AND scope = ?", org.ID, scope).Update("active", false)
	}
	g.JSON(http.StatusCreated, gin.H{"id": org.ID, "organization": org})
}

// AdmOrganizationUpdate godoc
//
//	@Summary	Update an organization (admin)
//	@Tags		admin
//	@Accept		json
//	@Param		id	path	string	true	"Organization ID"
//	@Param		body	body	object	true	"JSON with name, template, active"
//	@Success	200	{object}	map[string]interface{}
//	@Failure	400	{object}	map[string]interface{}
//	@Router		/adm/organizations/{id} [put]
func (c *Controller) AdmOrganizationUpdate(g *gin.Context) {
	id := g.Param("id")
	var org model.Organization
	if err := c.datastore.First(&org, "id = ?", id).Error; err != nil {
		g.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	var body struct {
		Name     *string `json:"name"`
		Template *string `json:"template"`
		Active   *bool   `json:"active"`
	}
	if err := g.ShouldBindJSON(&body); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if body.Name != nil {
		org.Name = *body.Name
	}
	if body.Template != nil {
		if err := model.ValidateOrganizationTemplate(*body.Template); err != nil {
			g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var videoCount int64
		c.datastore.Model(&model.Video{}).Where("organization_id = ?", org.ID).Count(&videoCount)
		if videoCount > 0 && *body.Template != org.Template {
			g.JSON(http.StatusBadRequest, gin.H{
				"error": "cannot change template of an organization used by existing videos; create a new organization instead",
			})
			return
		}
		org.Template = *body.Template
	}
	if body.Active != nil {
		org.Active = *body.Active
	}
	if err := c.datastore.Save(&org).Error; err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if org.Active {
		scope := org.Scope
		if scope == "" {
			scope = model.OrganizationScopeVideo
		}
		c.datastore.Model(&model.Organization{}).Where("id != ? AND scope = ?", org.ID, scope).Update("active", false)
	}
	g.JSON(http.StatusOK, gin.H{"organization": org})
}

// AdmOrganizationDelete godoc
//
//	@Summary	Delete an organization (admin)
//	@Tags		admin
//	@Param		id	path	string	true	"Organization ID"
//	@Success	204
//	@Failure	400	{object}	map[string]interface{}
//	@Router		/adm/organizations/{id} [delete]
func (c *Controller) AdmOrganizationDelete(g *gin.Context) {
	id := g.Param("id")
	var org model.Organization
	if err := c.datastore.First(&org, "id = ?", id).Error; err != nil {
		g.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	if org.Active {
		g.JSON(http.StatusBadRequest, gin.H{"error": "cannot delete the active organization; activate another one first"})
		return
	}
	var videoCount int64
	c.datastore.Model(&model.Video{}).Where("organization_id = ?", id).Count(&videoCount)
	if videoCount > 0 {
		g.JSON(http.StatusBadRequest, gin.H{"error": "organization is still used by videos; reorganize them first"})
		return
	}
	if err := c.datastore.Delete(&org).Error; err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	g.Status(http.StatusNoContent)
}

// AdmOrganizationActivate godoc
//
//	@Summary	Mark an organization as the active one (admin)
//	@Tags		admin
//	@Param		id	path	string	true	"Organization ID"
//	@Success	200	{object}	map[string]interface{}
//	@Failure	404	{object}	map[string]interface{}
//	@Router		/adm/organizations/{id}/activate [post]
func (c *Controller) AdmOrganizationActivate(g *gin.Context) {
	id := g.Param("id")
	var org model.Organization
	if err := c.datastore.First(&org, "id = ?", id).Error; err != nil {
		g.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	scope := org.Scope
	if scope == "" {
		scope = model.OrganizationScopeVideo
	}
	err := c.datastore.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Organization{}).Where("id != ? AND scope = ?", org.ID, scope).Update("active", false).Error; err != nil {
			return err
		}
		return tx.Model(&org).Update("active", true).Error
	})
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	g.JSON(http.StatusOK, gin.H{"organization": org})
}

// AdmOrganizationReorganize godoc
//
//	@Summary	Enqueue a reorganize task to move videos onto the target organization layout (admin)
//	@Tags		admin
//	@Param		id	path	string	true	"Target Organization ID"
//	@Success	202	{object}	map[string]interface{}
//	@Failure	404	{object}	map[string]interface{}
//	@Router		/adm/organizations/{id}/reorganize [post]
func (c *Controller) AdmOrganizationReorganize(g *gin.Context) {
	if c.runner == nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": "task runner not available"})
		return
	}
	id := g.Param("id")
	var org model.Organization
	if err := c.datastore.First(&org, "id = ?", id).Error; err != nil {
		g.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	var videos []model.Video
	if err := c.datastore.Where("imported = ? AND (organization_id IS NULL OR organization_id != ?)", true, org.ID).Find(&videos).Error; err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(videos) == 0 {
		g.JSON(http.StatusOK, gin.H{"message": "no videos to reorganize", "queued": 0})
		return
	}
	queued := 0
	var lastErr error
	for i := range videos {
		err := c.runner.NewTask("video/reorganize", map[string]string{
			"videoID":              videos[i].ID,
			"targetOrganizationID": org.ID,
			"sourcePath":           videos[i].RelativePath(),
		})
		if err != nil {
			lastErr = err
			continue
		}
		queued++
	}
	if queued == 0 && lastErr != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": lastErr.Error()})
		return
	}
	g.JSON(http.StatusAccepted, gin.H{
		"message":  "reorganize tasks queued",
		"task":     "video/reorganize",
		"queued":   queued,
		"redirect": "/adm/tasks",
	})
}

// AdmConfigReorganizeOnImportUpdate toggles the global default for the
// "reorganize on import" behaviour. The per-import override (in the upload
// UI) takes precedence over this default.
//
//	@Summary	Toggle reorganize-on-import default (admin)
//	@Tags		admin
//	@Param		action	path	string	true	"enable or disable"
//	@Success	200	{object}	map[string]interface{}
//	@Failure	400	{object}	map[string]interface{}
//	@Router		/adm/config/reorganize-on-import/{action} [get]
func (c *Controller) AdmConfigReorganizeOnImportUpdate(g *gin.Context) {
	action := g.Param("action")
	if action != "enable" && action != "disable" {
		g.JSON(http.StatusBadRequest, gin.H{"error": "action must be enable or disable"})
		return
	}
	var cfg model.Configuration
	err := c.datastore.First(&cfg).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		cfg = model.Configuration{ID: 1}
	}
	cfg.ReorganizeOnImport = action == "enable"
	if err := c.datastore.Save(&cfg).Error; err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	g.JSON(http.StatusOK, gin.H{"reorganize_on_import": cfg.ReorganizeOnImport})
}
