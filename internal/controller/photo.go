package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/zobtube/zobtube/internal/model"
)

func (c *Controller) PhotoStream(g *gin.Context) {
	photo, ps, ok := c.loadPhotoWithPhotoset(g)
	if !ok {
		return
	}
	store, path, _ := c.resolvePhotoFile(ps, photo)
	if store == nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": "storage not available"})
		return
	}
	c.serveFromStorage(g, store, path)
}

func (c *Controller) PhotoThumbMini(g *gin.Context) {
	photo, ps, ok := c.loadPhotoWithPhotoset(g)
	if !ok {
		return
	}
	c.PhotoThumbMiniServe(g, photo, ps)
}

func (c *Controller) PhotoThumbMiniServe(g *gin.Context, photo *model.Photo, ps *model.Photoset) {
	thumbPath := photo.ThumbnailMiniRelativePath(ps)
	store, err := c.metadataStore(true)
	if err == nil {
		if exists, exErr := store.Exists(thumbPath); exErr == nil && exists {
			c.serveFromStorage(g, store, thumbPath)
			return
		}
	}
	// Fallback: serve full image when mini thumb is missing or metadata storage is unavailable.
	libStore, path, ok := c.resolvePhotoFile(ps, photo)
	if ok && libStore != nil {
		c.serveFromStorage(g, libStore, path)
		return
	}
	g.JSON(http.StatusNotFound, gin.H{"error": "not found"})
}

func (c *Controller) PhotoDelete(g *gin.Context) {
	photo, ps, ok := c.loadPhotoWithPhotoset(g)
	if !ok {
		return
	}
	store, err := c.storageResolver.Storage(c.photosetLibraryID(ps))
	if err == nil {
		_ = store.Delete(photo.RelativePath(ps))
	}
	meta, err := c.metadataStore(true)
	if err == nil {
		_ = meta.Delete(photo.ThumbnailMiniRelativePath(ps))
	}
	if err := c.datastore.Delete(photo).Error; err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	g.JSON(http.StatusOK, gin.H{})
}

func (c *Controller) PhotoEditChannel(g *gin.Context) {
	photo, _, ok := c.loadPhotoWithPhotoset(g)
	if !ok {
		return
	}
	channelID := g.PostForm("channel_id")
	if channelID == "" {
		photo.ChannelID = nil
	} else {
		ch := &model.Channel{ID: channelID}
		if c.datastore.First(ch).RowsAffected < 1 {
			g.JSON(http.StatusNotFound, gin.H{"error": "channel not found"})
			return
		}
		photo.ChannelID = &channelID
	}
	if err := c.datastore.Save(photo).Error; err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	g.JSON(http.StatusOK, gin.H{})
}

func (c *Controller) PhotoActors(g *gin.Context) {
	c.photoAssociationMutate(g, "Actors", "actor_id")
}

func (c *Controller) PhotoCategories(g *gin.Context) {
	c.photoAssociationMutate(g, "Categories", "category_id")
}

func (c *Controller) photoAssociationMutate(g *gin.Context, assoc, param string) {
	photo, _, ok := c.loadPhotoWithPhotoset(g)
	if !ok {
		return
	}
	var target any
	switch assoc {
	case "Actors":
		target = &model.Actor{ID: g.Param(param)}
	case "Categories":
		target = &model.CategorySub{ID: g.Param(param)}
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
		res = c.datastore.Model(photo).Association(assoc).Append(target)
	} else {
		res = c.datastore.Model(photo).Association(assoc).Delete(target)
	}
	if res != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": res.Error()})
		return
	}
	g.JSON(http.StatusOK, gin.H{})
}

func (c *Controller) loadPhotoWithPhotoset(g *gin.Context) (*model.Photo, *model.Photoset, bool) {
	photoID := g.Param("photo_id")
	if photoID == "" {
		photoID = g.Param("id")
	}
	photo := &model.Photo{ID: photoID}
	if c.datastore.First(photo).RowsAffected < 1 {
		g.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return nil, nil, false
	}
	ps := &model.Photoset{ID: photo.PhotosetID}
	if c.datastore.First(ps).RowsAffected < 1 {
		g.JSON(http.StatusNotFound, gin.H{"error": "photoset not found"})
		return nil, nil, false
	}
	return photo, ps, true
}
