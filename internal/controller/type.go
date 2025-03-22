package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/zobtube/zobtube/internal/config"
	"github.com/zobtube/zobtube/internal/model"
	"github.com/zobtube/zobtube/internal/provider"

	"gorm.io/gorm"
)

type AbtractController interface {
	// Back office
	AdmHome(c *gin.Context)

	// Home
	Home(c *gin.Context)

	// Auth
	AuthPage(*gin.Context)
	AuthLogin(*gin.Context)
	AuthLogout(*gin.Context)
	GetSession(*model.UserSession) *gorm.DB
	GetUser(*model.User) *gorm.DB

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

	// Video, used for Clips, Movies and Videos
	VideoAjaxActors(c *gin.Context)
	VideoAjaxComputeDuration(c *gin.Context)
	VideoAjaxGenerateThumbnail(c *gin.Context)
	VideoAjaxGenerateThumbnailXS(c *gin.Context)
	VideoAjaxImport(c *gin.Context)
	VideoAjaxRename(c *gin.Context)
	VideoAjaxUpload(c *gin.Context)
	VideoAjaxUploadThumb(c *gin.Context)
	VideoAjaxCreate(c *gin.Context)
	VideoAjaxStreamInfo(c *gin.Context)
	VideoAjaxDelete(c *gin.Context)
	VideoAjaxMigrate(c *gin.Context)
	VideoEdit(c *gin.Context)
	VideoStream(c *gin.Context)
	VideoThumb(c *gin.Context)
	VideoThumbXS(c *gin.Context)
	VideoView(c *gin.Context)

	// Video Views
	VideoViewAjaxIncrement(g *gin.Context)

	ClipList(c *gin.Context)
	MovieList(c *gin.Context)
	VideoList(c *gin.Context)
	GenericVideoList(vt string, c *gin.Context)

	// Channels
	ChannelCreate(c *gin.Context)
	ChannelList(c *gin.Context)
	ChannelView(c *gin.Context)
	ChannelThumb(c *gin.Context)

	// Uploads
	UploadTriage(c *gin.Context)
	UploadPreview(c *gin.Context)
	UploadImport(c *gin.Context)
	UploadAjaxTriageFolder(*gin.Context)
	UploadAjaxTriageFile(*gin.Context)
	UploadAjaxUploadFile(*gin.Context)

	// Providers
	ProviderRegister(provider.Provider)
	ProviderGet(string) (provider.Provider, error)

	// Cleanup
	CleanupRoutine()

	// Profile
	ProfileView(c *gin.Context)

	// Failsafe
	FailsafeConfiguration(c *gin.Context)
	FailsafeUser(c *gin.Context)
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
