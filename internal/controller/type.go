package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"github.com/zobtube/zobtube/internal/config"
	"github.com/zobtube/zobtube/internal/model"
	"github.com/zobtube/zobtube/internal/provider"
	"github.com/zobtube/zobtube/internal/runner"
)

type AbtractController interface {
	// Back office
	AdmHome(*gin.Context)
	AdmVideoList(*gin.Context)
	AdmActorList(*gin.Context)
	AdmChannelList(*gin.Context)
	AdmCategory(*gin.Context)
	AdmConfigAuth(*gin.Context)
	AdmConfigAuthUpdate(*gin.Context)
	AdmConfigProvider(*gin.Context)
	AdmConfigProviderSwitch(*gin.Context)
	AdmConfigOfflineMode(*gin.Context)
	AdmConfigOfflineModeUpdate(*gin.Context)
	AdmTaskHome(*gin.Context)
	AdmTaskList(*gin.Context)
	AdmTaskRetry(*gin.Context)
	AdmTaskView(*gin.Context)
	AdmUserList(*gin.Context)
	AdmUserNew(*gin.Context)
	AdmUserDelete(*gin.Context)

	// Home
	Home(*gin.Context)

	// Auth
	AuthenticationEnabled() bool
	AuthPage(*gin.Context)
	AuthLogin(*gin.Context)
	AuthLogout(*gin.Context)
	GetSession(*model.UserSession) *gorm.DB
	GetUser(*model.User) *gorm.DB
	GetFirstUser(*model.User) *gorm.DB

	// Actors
	ActorAjaxCategories(*gin.Context)
	ActorAjaxLinkThumbGet(*gin.Context)
	ActorAjaxLinkThumbDelete(*gin.Context)
	ActorAjaxNew(*gin.Context)
	ActorAjaxProviderSearch(*gin.Context)
	ActorAjaxRename(*gin.Context)
	ActorAjaxThumb(*gin.Context)
	ActorAjaxLinkCreate(*gin.Context)
	ActorAjaxAliasCreate(*gin.Context)
	ActorAjaxAliasRemove(*gin.Context)
	ActorEdit(*gin.Context)
	ActorList(*gin.Context)
	ActorNew(*gin.Context)
	ActorView(*gin.Context)
	ActorThumb(*gin.Context)
	ActorDelete(*gin.Context)

	// Categories
	CategoryAjaxAdd(*gin.Context)
	CategoryAjaxDelete(*gin.Context)
	CategoryList(*gin.Context)

	// Sub categories
	CategorySubAjaxAdd(*gin.Context)
	CategorySubAjaxRename(*gin.Context)
	CategorySubAjaxThumbSet(*gin.Context)
	CategorySubAjaxThumbRemove(*gin.Context)
	CategorySubThumb(*gin.Context)
	CategorySubView(*gin.Context)

	// Video, used for Clips, Movies and Videos
	VideoAjaxGet(*gin.Context)
	VideoAjaxActors(*gin.Context)
	VideoAjaxCategories(*gin.Context)
	VideoAjaxRename(*gin.Context)
	VideoAjaxUpload(*gin.Context)
	VideoAjaxUploadThumb(*gin.Context)
	VideoAjaxCreate(*gin.Context)
	VideoAjaxStreamInfo(*gin.Context)
	VideoAjaxDelete(*gin.Context)
	VideoAjaxMigrate(*gin.Context)
	VideoAjaxGenerateThumbnail(*gin.Context)
	VideoEdit(*gin.Context)
	VideoAjaxEditChannel(*gin.Context)
	VideoStream(*gin.Context)
	VideoThumb(*gin.Context)
	VideoThumbXS(*gin.Context)
	VideoView(*gin.Context)

	// Video Views
	VideoViewAjaxIncrement(*gin.Context)

	ClipList(*gin.Context)
	ClipView(*gin.Context)
	MovieList(*gin.Context)
	VideoList(*gin.Context)
	GenericVideoList(string, *gin.Context)

	// Channels
	ChannelCreate(*gin.Context)
	ChannelList(*gin.Context)
	ChannelView(*gin.Context)
	ChannelThumb(*gin.Context)
	ChannelAjaxList(*gin.Context)
	ChannelEdit(*gin.Context)

	// Uploads
	UploadTriage(*gin.Context)
	UploadPreview(*gin.Context)
	UploadImport(*gin.Context)
	UploadAjaxTriageFolder(*gin.Context)
	UploadAjaxTriageFile(*gin.Context)
	UploadAjaxUploadFile(*gin.Context)
	UploadAjaxDeleteFile(*gin.Context)
	UploadAjaxFolderCreate(*gin.Context)
	UploadAjaxMassDelete(*gin.Context)
	UploadAjaxMassImport(*gin.Context)

	// Providers
	ProviderRegister(provider.Provider) error
	ProviderGet(string) (provider.Provider, error)

	// Profile
	ProfileView(*gin.Context)

	// Error pages
	ErrUnauthorized(*gin.Context)

	// Init
	LoggerRegister(*zerolog.Logger)
	ConfigurationRegister(*config.Config)
	DatabaseRegister(*gorm.DB)
	RunnerRegister(*runner.Runner)
	ConfigurationFromDBApply(*model.Configuration)
	BuildDetailsRegister(string, string, string)

	// Cleanup
	CleanupRoutine()

	// Maintenance
	Restart()
	Shutdown()
}

type buildDetails struct {
	Version   string
	Commit    string
	BuildDate string
}

type Controller struct {
	config          *config.Config
	datastore       *gorm.DB
	providers       map[string]provider.Provider
	shutdownChannel chan<- int
	runner          *runner.Runner
	build           *buildDetails
	logger          *zerolog.Logger
}

func New(shutdownChannel chan int) AbtractController {
	return &Controller{
		providers:       make(map[string]provider.Provider),
		shutdownChannel: shutdownChannel,
	}
}

func (c *Controller) ConfigurationRegister(cfg *config.Config) {
	c.config = cfg
}

func (c *Controller) DatabaseRegister(db *gorm.DB) {
	c.datastore = db
}

func (c *Controller) RunnerRegister(r *runner.Runner) {
	c.runner = r
}

func (c *Controller) ConfigurationFromDBApply(db *model.Configuration) {
	c.logger.Info().Str("kind", "system").Bool("authentication", db.UserAuthentication).Send()
	c.config.Authentication = db.UserAuthentication
}
func (c *Controller) LoggerRegister(logger *zerolog.Logger) {
	c.logger = logger
}

func (c *Controller) BuildDetailsRegister(version, commit, buildDate string) {
	c.build = &buildDetails{
		Version:   version,
		Commit:    commit,
		BuildDate: buildDate,
	}
}
