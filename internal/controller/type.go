package controller

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/zobtube/zobtube/internal/config"
	"gitlab.com/zobtube/zobtube/internal/provider"
	"gorm.io/gorm"
)

type AbtractController interface {
	// Back office
	AdmHome(c *gin.Context)

	// Home
	Home(c *gin.Context)

	// Actors
	ActorAjaxLinkThumbGet(c *gin.Context)
	ActorAjaxLinkThumbDelete(c *gin.Context)
	ActorAjaxNew(c *gin.Context)
	ActorAjaxProviderSearch(c *gin.Context)
	ActorAjaxThumb(c *gin.Context)
	ActorEdit(c *gin.Context)
	ActorList(c *gin.Context)
	ActorNew(c *gin.Context)
	ActorView(c *gin.Context)
	ActorThumb(c *gin.Context)

	// Generic Video, used for Clips, Movies and Videos
	GenericVideoAjaxActors(c *gin.Context)
	GenericVideoAjaxComputeDuration(c *gin.Context)
	GenericVideoAjaxGenerateThumbnail(c *gin.Context)
	GenericVideoAjaxGenerateThumbnailXS(c *gin.Context)
	GenericVideoAjaxImport(c *gin.Context)
	GenericVideoAjaxRename(c *gin.Context)
	GenericVideoAjaxUpload(c *gin.Context)
	GenericVideoAjaxUploadThumb(c *gin.Context)
	GenericVideoAjaxCreate(c *gin.Context)
	GenericVideoAjaxStreamInfo(c *gin.Context)
	GenericVideoEdit(c *gin.Context)
	GenericVideoList(vt string, c *gin.Context)
	GenericVideoStream(vt string, c *gin.Context)
	GenericVideoThumb(vt string, c *gin.Context)
	GenericVideoThumbXS(vt string, c *gin.Context)
	GenericVideoView(vt string, c *gin.Context)

	// Channels
	ChannelCreate(c *gin.Context)
	ChannelList(c *gin.Context)
	ChannelView(c *gin.Context)
	ChannelThumb(c *gin.Context)

	// Clips
	ClipList(c *gin.Context)
	ClipView(c *gin.Context)
	ClipStream(c *gin.Context)
	ClipThumb(c *gin.Context)
	ClipThumbXS(c *gin.Context)

	// Movies
	MovieList(c *gin.Context)
	MovieView(c *gin.Context)
	MovieStream(c *gin.Context)
	MovieThumb(c *gin.Context)
	MovieThumbXS(c *gin.Context)

	// Videos
	VideoList(c *gin.Context)
	VideoView(c *gin.Context)
	VideoStream(c *gin.Context)
	VideoEdit(c *gin.Context)
	VideoThumb(c *gin.Context)
	VideoThumbXS(c *gin.Context)

	// Uploads
	UploadHome(c *gin.Context)
	UploadTriage(c *gin.Context)
	UploadPreview(c *gin.Context)
	UploadImport(c *gin.Context)

	// Providers
	ProviderRegister(provider.Provider)
	ProviderGet(string) (provider.Provider, error)
}

type Controller struct {
	config    *config.Config
	datastore *gorm.DB
	providers map[string]provider.Provider
}

func New(cfg *config.Config, db *gorm.DB) AbtractController {
	return &Controller{
		config:    cfg,
		datastore: db,
		providers: make(map[string]provider.Provider),
	}
}
