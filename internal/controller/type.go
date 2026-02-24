package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"github.com/zobtube/zobtube/internal/config"
	"github.com/zobtube/zobtube/internal/model"
	"github.com/zobtube/zobtube/internal/provider"
	"github.com/zobtube/zobtube/internal/runner"
	"github.com/zobtube/zobtube/internal/swagger"
)

type AbstractController interface {
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
	AuthLogin(*gin.Context)
	AuthLogout(*gin.Context)
	AuthLogoutRedirect(*gin.Context)
	AuthMe(*gin.Context)
	GetSession(*model.UserSession) *gorm.DB
	GetUser(*model.User) *gorm.DB
	GetFirstUser(*model.User) *gorm.DB

	// Actors
	ActorCategories(*gin.Context)
	ActorLinkThumbGet(*gin.Context)
	ActorLinkThumbDelete(*gin.Context)
	ActorNew(*gin.Context)
	ActorProviderSearch(*gin.Context)
	ActorRename(*gin.Context)
	ActorDescription(*gin.Context)
	ActorUploadThumb(*gin.Context)
	ActorLinkCreate(*gin.Context)
	ActorAliasCreate(*gin.Context)
	ActorAliasRemove(*gin.Context)
	ActorMerge(*gin.Context)
	ActorList(*gin.Context)
	ActorGet(*gin.Context)
	ActorDelete(*gin.Context)
	ActorThumb(*gin.Context)

	// Categories
	CategoryAdd(*gin.Context)
	CategoryDelete(*gin.Context)
	CategoryList(*gin.Context)
	CategorySubGet(*gin.Context)

	// Sub categories
	CategorySubAdd(*gin.Context)
	CategorySubRename(*gin.Context)
	CategorySubThumbSet(*gin.Context)
	CategorySubThumbRemove(*gin.Context)
	CategorySubThumb(*gin.Context)

	// Video, used for Clips, Movies and Videos
	VideoGet(*gin.Context)
	VideoActors(*gin.Context)
	VideoCategories(*gin.Context)
	VideoRename(*gin.Context)
	VideoUpload(*gin.Context)
	VideoUploadThumb(*gin.Context)
	VideoCreate(*gin.Context)
	VideoStreamInfo(*gin.Context)
	VideoDelete(*gin.Context)
	VideoMigrate(*gin.Context)
	VideoGenerateThumbnail(*gin.Context)
	VideoEditChannel(*gin.Context)
	VideoStream(*gin.Context)
	VideoThumb(*gin.Context)
	VideoThumbXS(*gin.Context)

	// Video Views
	VideoViewIncrement(*gin.Context)

	ClipList(*gin.Context)
	ClipView(*gin.Context)
	MovieList(*gin.Context)
	VideoList(*gin.Context)
	VideoView(*gin.Context)
	VideoEdit(*gin.Context)

	// Channels
	ChannelList(*gin.Context)
	ChannelGet(*gin.Context)
	ChannelCreate(*gin.Context)
	ChannelUpdate(*gin.Context)
	ChannelThumb(*gin.Context)
	ChannelMap(*gin.Context)

	// Uploads
	UploadImport(*gin.Context)
	UploadPreview(*gin.Context)
	UploadTriageFolder(*gin.Context)
	UploadTriageFile(*gin.Context)
	UploadFile(*gin.Context)
	UploadDeleteFile(*gin.Context)
	UploadFolderCreate(*gin.Context)
	UploadMassDelete(*gin.Context)
	UploadMassImport(*gin.Context)

	// Providers
	ProviderRegister(provider.Provider) error
	ProviderGet(string) (provider.Provider, error)

	// Profile
	ProfileView(*gin.Context)

	// Error pages
	ErrNotFound(*gin.Context)
	ErrUnauthorized(*gin.Context)

	// SPA
	Bootstrap(*gin.Context)
	SPAApp(*gin.Context)
	NoRouteOrSPA(*gin.Context)

	// Init
	LoggerRegister(*zerolog.Logger)
	ConfigurationRegister(*config.Config)
	DatabaseRegister(*gorm.DB)
	RunnerRegister(*runner.Runner)
	ConfigurationFromDBApply(*model.Configuration)
	BuildDetailsRegister(string, string, string)
	RegisterError(string)

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
	healthError     []string
}

func New(shutdownChannel chan int) AbstractController {
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

	swagger.SwaggerInfo.Version = version
}

func (c *Controller) RegisterError(err string) {
	c.healthError = append(c.healthError, err)
}
