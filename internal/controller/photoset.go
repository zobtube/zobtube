package controller

import (
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/zobtube/zobtube/internal/model"
	"github.com/zobtube/zobtube/internal/storage"
)

func (c *Controller) PhotosetList(g *gin.Context) {
	var items []model.Photoset
	c.datastore.Order("created_at desc").Find(&items)
	g.JSON(http.StatusOK, gin.H{"items": items, "total": len(items)})
}

func (c *Controller) PhotosetView(g *gin.Context) {
	id := g.Param("id")
	ps := &model.Photoset{ID: id}
	if c.datastore.Preload("Actors.Categories").Preload("Channel").Preload("Categories").
		Preload("Photos", func(db *gorm.DB) *gorm.DB {
			return db.Order("position asc, filename asc")
		}).
		Preload("Photos.Actors.Categories").Preload("Photos.Categories").Preload("Photos.Channel").
		First(ps).RowsAffected < 1 {
		g.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	categories := mergeCategoryMap(ps.Categories, nil)
	for _, photo := range ps.Photos {
		for id, name := range mergeCategoryMap(nil, photo.Categories) {
			categories[id] = name
		}
		for _, actor := range photo.Actors {
			for _, cat := range actor.Categories {
				categories[cat.ID] = cat.Name
			}
		}
	}
	for _, actor := range ps.Actors {
		for _, cat := range actor.Categories {
			categories[cat.ID] = cat.Name
		}
	}
	type photoView struct {
		model.Photo
		EffectiveActors     []model.Actor `json:"effective_actors"`
		EffectiveCategories []model.CategorySub `json:"effective_categories"`
		EffectiveChannel    *model.Channel `json:"effective_channel"`
	}
	photos := make([]photoView, 0, len(ps.Photos))
	for i := range ps.Photos {
		p := ps.Photos[i]
		photos = append(photos, photoView{
			Photo:               p,
			EffectiveActors:     mergeActors(ps.Actors, p.Actors),
			EffectiveCategories: mergeActorsAsCategories(ps.Categories, p.Categories),
			EffectiveChannel:    effectiveChannel(ps, &p),
		})
	}
	g.JSON(http.StatusOK, gin.H{
		"photoset":   ps,
		"photos":     photos,
		"categories": categories,
	})
}

func mergeActorsAsCategories(base, extra []model.CategorySub) []model.CategorySub {
	seen := make(map[string]struct{})
	out := make([]model.CategorySub, 0, len(base)+len(extra))
	for _, cat := range base {
		if _, ok := seen[cat.ID]; ok {
			continue
		}
		seen[cat.ID] = struct{}{}
		out = append(out, cat)
	}
	for _, cat := range extra {
		if _, ok := seen[cat.ID]; ok {
			continue
		}
		seen[cat.ID] = struct{}{}
		out = append(out, cat)
	}
	return out
}

func (c *Controller) PhotosetEdit(g *gin.Context) {
	id := g.Param("id")
	ps := &model.Photoset{ID: id}
	if c.datastore.Preload("Actors").Preload("Channel").Preload("Categories").Preload("Organization").
		Preload("Photos", func(db *gorm.DB) *gorm.DB {
			return db.Order("position asc, filename asc")
		}).First(ps).RowsAffected < 1 {
		g.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	var actors []model.Actor
	c.datastore.Find(&actors)
	var categories []model.Category
	c.datastore.Preload("Sub").Find(&categories)
	var libraries []model.Library
	c.datastore.Order("created_at").Find(&libraries)
	activeOrg, _ := model.ActiveOrganizationForScope(c.datastore, model.OrganizationScopePhotoset)
	resp := gin.H{
		"photoset":   ps,
		"actors":     actors,
		"categories": categories,
		"libraries":  libraries,
		"organized":  ps.IsOrganizedWith(activeOrg),
	}
	if activeOrg != nil {
		resp["active_organization"] = activeOrg
	}
	if ps.NeedsReorganization(activeOrg) {
		resp["needs_organize"] = true
	}
	g.JSON(http.StatusOK, resp)
}

func (c *Controller) PhotosetCreate(g *gin.Context) {
	form := struct {
		Name      string `form:"name" json:"name"`
		LibraryID string `form:"library_id" json:"library_id"`
	}{}
	if err := g.ShouldBind(&form); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if strings.TrimSpace(form.Name) == "" {
		g.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}
	libID := c.uploadLibraryID(form.LibraryID)
	ps := &model.Photoset{
		Name:      strings.TrimSpace(form.Name),
		LibraryID: &libID,
		Status:    model.PhotosetStatusCreating,
	}
	if err := c.datastore.Create(ps).Error; err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	g.JSON(http.StatusOK, gin.H{"photoset_id": ps.ID, "id": ps.ID})
}

func (c *Controller) PhotosetUploadFiles(g *gin.Context) {
	id := g.Param("id")
	ps := &model.Photoset{ID: id}
	if c.datastore.First(ps).RowsAffected < 1 {
		g.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	store, err := c.storageResolver.Storage(c.photosetLibraryID(ps))
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := store.MkdirAll(ps.FolderRelativePath()); err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if g.Query("done") == "1" {
		ps.Imported = true
		_ = c.datastore.Save(ps).Error
		if c.runner != nil {
			_ = c.runner.NewTask("photoset/finalize", map[string]string{"photosetID": ps.ID})
		}
		g.JSON(http.StatusOK, gin.H{"done": true})
		return
	}

	form, err := g.MultipartForm()
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	files := form.File["file"]
	if len(files) == 0 {
		g.JSON(http.StatusBadRequest, gin.H{"error": "no files uploaded"})
		return
	}

	imported := 0
	for _, fh := range files {
		base := sanitizePhotoBasename(fh.Filename)
		if base == "" {
			continue
		}
		mimeType, ok := photoMimeForFilename(base)
		if !ok {
			continue
		}
		rel := filepath.Join(ps.FolderRelativePath(), base)
		src, err := fh.Open()
		if err != nil {
			g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		dst, err := store.Create(rel)
		if err != nil {
			src.Close()
			g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		_, copyErr := io.Copy(dst, src)
		src.Close()
		_ = dst.Close()
		if copyErr != nil {
			g.JSON(http.StatusInternalServerError, gin.H{"error": copyErr.Error()})
			return
		}

		var existing model.Photo
		if c.datastore.Where("photoset_id = ? AND filename = ?", ps.ID, base).First(&existing).RowsAffected > 0 {
			path := rel
			existing.Path = &path
			existing.Mime = mimeType
			existing.SizeBytes = fh.Size
			_ = c.datastore.Save(&existing).Error
			imported++
			continue
		}
		path := rel
		photo := &model.Photo{
			PhotosetID: ps.ID,
			Filename:   base,
			Path:       &path,
			Mime:       mimeType,
			SizeBytes:  fh.Size,
		}
		if err := c.datastore.Create(photo).Error; err != nil {
			g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		imported++
	}

	if g.Query("done") == "1" {
		ps.Imported = true
		_ = c.datastore.Save(ps).Error
		if c.runner != nil {
			_ = c.runner.NewTask("photoset/finalize", map[string]string{"photosetID": ps.ID})
		}
	}
	g.JSON(http.StatusOK, gin.H{"imported": imported, "done": g.Query("done") == "1"})
}

func (c *Controller) PhotosetUploadArchive(g *gin.Context) {
	id := g.Param("id")
	ps := &model.Photoset{ID: id}
	if c.datastore.First(ps).RowsAffected < 1 {
		g.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	store, err := c.storageResolver.Storage(c.photosetLibraryID(ps))
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	archivePath := filepath.Join(ps.FolderRelativePath(), "_upload.zip")
	if err := store.MkdirAll(ps.FolderRelativePath()); err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fh, err := g.FormFile("file")
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	src, err := fh.Open()
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer src.Close()
	dst, err := store.Create(archivePath)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if _, err := io.Copy(dst, src); err != nil {
		dst.Close()
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	_ = dst.Close()

	if c.runner == nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": "task runner not available"})
		return
	}
	if err := c.runner.NewTask("photoset/unzip", map[string]string{
		"photosetID":  ps.ID,
		"archivePath": archivePath,
	}); err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	g.JSON(http.StatusAccepted, gin.H{
		"message":  "archive upload queued for extraction",
		"redirect": "/adm/tasks",
	})
}

func (c *Controller) PhotosetDelete(g *gin.Context) {
	id := g.Param("id")
	ps := &model.Photoset{ID: id}
	if c.datastore.First(ps).RowsAffected < 1 {
		g.JSON(http.StatusNotFound, gin.H{})
		return
	}
	ps.Status = model.PhotosetStatusDeleting
	if err := c.datastore.Save(ps).Error; err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if c.runner == nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": "task runner not available"})
		return
	}
	if err := c.runner.NewTask("photoset/delete", map[string]string{"photosetID": ps.ID}); err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	g.JSON(http.StatusOK, gin.H{})
}

type PhotosetRenameForm struct {
	Name string `form:"name" json:"name"`
}

func (c *Controller) PhotosetRename(g *gin.Context) {
	var form PhotosetRenameForm
	if err := g.ShouldBind(&form); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ps := &model.Photoset{ID: g.Param("id")}
	if c.datastore.First(ps).RowsAffected < 1 {
		g.JSON(http.StatusNotFound, gin.H{})
		return
	}
	ps.Name = strings.TrimSpace(form.Name)
	if ps.Name == "" {
		g.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}
	if err := c.datastore.Save(ps).Error; err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	g.JSON(http.StatusOK, gin.H{})
}

func (c *Controller) PhotosetEditChannel(g *gin.Context) {
	ps := &model.Photoset{ID: g.Param("id")}
	if c.datastore.First(ps).RowsAffected < 1 {
		g.JSON(http.StatusNotFound, gin.H{})
		return
	}
	channelID := g.PostForm("channel_id")
	if channelID == "" {
		ps.ChannelID = nil
	} else {
		ch := &model.Channel{ID: channelID}
		if c.datastore.First(ch).RowsAffected < 1 {
			g.JSON(http.StatusNotFound, gin.H{"error": "channel not found"})
			return
		}
		ps.ChannelID = &channelID
	}
	if err := c.datastore.Save(ps).Error; err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	g.JSON(http.StatusOK, gin.H{})
}

func (c *Controller) PhotosetActors(g *gin.Context) {
	c.photosetAssociationMutate(g, "Actors")
}

func (c *Controller) PhotosetCategories(g *gin.Context) {
	c.photosetAssociationMutate(g, "Categories")
}

func (c *Controller) photosetAssociationMutate(g *gin.Context, assoc string) {
	ps := &model.Photoset{ID: g.Param("id")}
	if c.datastore.First(ps).RowsAffected < 1 {
		g.JSON(http.StatusNotFound, gin.H{})
		return
	}
	var target any
	switch assoc {
	case "Actors":
		target = &model.Actor{ID: g.Param("actor_id")}
	case "Categories":
		target = &model.CategorySub{ID: g.Param("category_id")}
	default:
		g.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	if c.datastore.First(target).RowsAffected < 1 {
		g.JSON(http.StatusNotFound, gin.H{})
		return
	}
	var res error
	if g.Request.Method == "PUT" {
		res = c.datastore.Model(ps).Association(assoc).Append(target)
	} else {
		res = c.datastore.Model(ps).Association(assoc).Delete(target)
	}
	if res != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": res.Error()})
		return
	}
	g.JSON(http.StatusOK, gin.H{})
}

func (c *Controller) PhotosetSetCover(g *gin.Context) {
	ps := &model.Photoset{ID: g.Param("id")}
	if c.datastore.First(ps).RowsAffected < 1 {
		g.JSON(http.StatusNotFound, gin.H{})
		return
	}
	photoID := g.Param("photo_id")
	photo := &model.Photo{ID: photoID}
	if c.datastore.Where("photoset_id = ?", ps.ID).First(photo).RowsAffected < 1 {
		g.JSON(http.StatusNotFound, gin.H{"error": "photo not found"})
		return
	}
	ps.CoverPhotoID = &photoID
	if err := c.datastore.Save(ps).Error; err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	g.JSON(http.StatusOK, gin.H{})
}

func (c *Controller) PhotosetReorganize(g *gin.Context) {
	if c.runner == nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": "task runner not available"})
		return
	}
	ps := &model.Photoset{ID: g.Param("id")}
	if c.datastore.Preload("Photos").First(ps).RowsAffected < 1 {
		g.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	if !ps.Imported {
		g.JSON(http.StatusBadRequest, gin.H{"error": "photoset has not been imported yet"})
		return
	}
	activeOrg, err := model.ActiveOrganizationForScope(c.datastore, model.OrganizationScopePhotoset)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": "no active organization"})
		return
	}
	if !ps.NeedsReorganization(activeOrg) {
		g.JSON(http.StatusOK, gin.H{"message": "photoset already follows the active organization"})
		return
	}
	if err := c.runner.NewTask("photoset/reorganize", map[string]string{
		"photosetID":           ps.ID,
		"targetOrganizationID": activeOrg.ID,
	}); err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	g.JSON(http.StatusAccepted, gin.H{
		"message":  "reorganize task queued",
		"redirect": "/adm/tasks",
	})
}

func (c *Controller) PhotosetCover(g *gin.Context) {
	ps := &model.Photoset{ID: g.Param("id")}
	if c.datastore.First(ps).RowsAffected < 1 {
		g.JSON(http.StatusNotFound, gin.H{})
		return
	}
	var photo model.Photo
	if ps.CoverPhotoID != nil && *ps.CoverPhotoID != "" {
		if c.datastore.First(&photo, "id = ?", *ps.CoverPhotoID).RowsAffected < 1 {
			photo = model.Photo{}
		}
	}
	if photo.ID == "" {
		c.datastore.Where("photoset_id = ?", ps.ID).Order("position asc, filename asc").First(&photo)
	}
	if photo.ID == "" {
		g.JSON(http.StatusNotFound, gin.H{})
		return
	}
	c.PhotoThumbMiniServe(g, &photo, ps)
}

func (c *Controller) resolvePhotoFile(ps *model.Photoset, photo *model.Photo) (storage.Storage, string, bool) {
	candidates := photo.StoragePathCandidates(ps)
	libID := c.photosetLibraryID(ps)
	store, err := c.storageResolver.Storage(libID)
	if err == nil {
		if p, ok, _ := storage.FirstExistingPath(store, candidates); ok {
			return store, p, true
		}
	}
	if store != nil {
		return store, photo.RelativePath(ps), false
	}
	return nil, "", false
}
