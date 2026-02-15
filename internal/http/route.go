package http

import (
	"github.com/zobtube/zobtube/internal/controller"
)

func (s *Server) setupRoutes(c controller.AbstractController) {
	// Auth (no auth required): SPA at /auth shows login form; POST /auth/login for login
	auth := s.Router.Group("/auth")
	auth.GET("", c.SPAApp)
	auth.POST("/login", c.AuthLogin)
	auth.GET("/logout", c.AuthLogoutRedirect)

	// SPA shell - single endpoint for the app
	s.Router.GET("/", c.SPAApp)

	// Bootstrap (unauthenticated) - returns auth_enabled and user for SPA init
	s.Router.GET("/api/bootstrap", c.Bootstrap)

	authGroup := s.Router.Group("")
	authGroup.Use(UserIsAuthenticated(c))

	admGroup := s.Router.Group("")
	admGroup.Use(UserIsAuthenticated(c))
	admGroup.Use(UserIsAdmin(c))

	// Logout (authenticated)
	authGroup.POST("/api/auth/logout", c.AuthLogout)
	authGroup.GET("/api/auth/me", c.AuthMe)

	// Home
	authGroup.GET("/api/home", c.Home)

	// Actors - list and get for all auth users; create/delete/mutate for admin
	authGroup.GET("/api/actor", c.ActorList)
	authGroup.GET("/api/actor/:id", c.ActorGet)
	authGroup.GET("/api/actor/:id/thumb", c.ActorThumb)

	actorGroup := admGroup.Group("/api/actor")
	{
		actorGroup.POST("/", c.ActorNew)
		actorGroup.DELETE("/:id", c.ActorDelete)
		actorGroup.POST("/:id/rename", c.ActorRename)
		actorGroup.POST("/:id/description", c.ActorDescription)
		actorGroup.POST("/:id/merge", c.ActorMerge)

		// providers
		actorGroup.GET("/:id/provider/:provider_slug", c.ActorProviderSearch)

		// links
		actorGroup.DELETE("/link/:id", c.ActorLinkThumbDelete)
		actorGroup.GET("/link/:id/thumb", c.ActorLinkThumbGet)
		actorGroup.POST("/:id/link", c.ActorLinkCreate)

		// thumb
		actorGroup.POST("/:id/thumb", c.ActorUploadThumb)

		// alias
		actorGroup.POST("/:id/alias", c.ActorAliasCreate)
		actorGroup.DELETE("/alias/:id", c.ActorAliasRemove)

		// categories
		actorGroup.PUT("/:id/category/:category_id", c.ActorCategories)
		actorGroup.DELETE("/:id/category/:category_id", c.ActorCategories)
	}

	// Categories
	authGroup.GET("/api/category", c.CategoryList)
	authGroup.GET("/api/category/:id", c.CategorySubGet)
	authGroup.GET("/api/category-sub/:id/thumb", c.CategorySubThumb)
	admGroup.POST("/api/category", c.CategoryAdd)
	admGroup.DELETE("/api/category/:id", c.CategoryDelete)
	admGroup.POST("/api/category-sub/:id/thumb", c.CategorySubThumbSet)
	admGroup.DELETE("/api/category-sub/:id/thumb", c.CategorySubThumbRemove)
	admGroup.POST("/api/category-sub", c.CategorySubAdd)
	admGroup.POST("/api/category-sub/:id/rename", c.CategorySubRename)

	// Channels
	authGroup.GET("/api/channel/map", c.ChannelMap)
	authGroup.GET("/api/channel", c.ChannelList)
	authGroup.GET("/api/channel/:id", c.ChannelGet)
	authGroup.GET("/api/channel/:id/thumb", c.ChannelThumb)
	admGroup.POST("/api/channel", c.ChannelCreate)
	admGroup.PUT("/api/channel/:id", c.ChannelUpdate)

	// Videos
	authGroup.GET("/api/clip", c.ClipList)
	authGroup.GET("/api/clip/:id", c.ClipView)
	authGroup.GET("/api/movie", c.MovieList)
	authGroup.GET("/api/video", c.VideoList)
	authGroup.GET("/api/video/:id", c.VideoView)
	admGroup.GET("/api/video/:id/edit", c.VideoEdit)
	authGroup.GET("/api/video/:id/summary", c.VideoGet)
	authGroup.GET("/api/video/:id/stream", c.VideoStream)
	authGroup.GET("/api/video/:id/thumb", c.VideoThumb)
	authGroup.GET("/api/video/:id/thumb_xs", c.VideoThumbXS)

	videoGroup := admGroup.Group("/api/video")
	{
		videoGroup.POST("", c.VideoCreate)
		videoGroup.HEAD("/:id", c.VideoStreamInfo)
		videoGroup.DELETE("/:id", c.VideoDelete)
		videoGroup.POST("/:id/upload", c.VideoUpload)
		videoGroup.POST("/:id/thumb", c.VideoUploadThumb)
		videoGroup.POST("/:id/migrate", c.VideoMigrate)
		videoGroup.PUT("/:id/actor/:actor_id", c.VideoActors)
		videoGroup.DELETE("/:id/actor/:actor_id", c.VideoActors)
		videoGroup.PUT("/:id/category/:category_id", c.VideoCategories)
		videoGroup.DELETE("/:id/category/:category_id", c.VideoCategories)
		videoGroup.POST("/:id/generate-thumbnail/:timing", c.VideoGenerateThumbnail)
		videoGroup.POST("/:id/rename", c.VideoRename)
		videoGroup.POST("/:id/count-view", c.VideoViewIncrement)
		videoGroup.POST("/:id/channel", c.VideoEditChannel)
	}

	// Uploads
	uploadGroup := admGroup.Group("/api/upload")
	{
		uploadGroup.POST("/import", c.UploadImport)
		uploadGroup.GET("/preview/:filepath", c.UploadPreview)
		uploadGroup.POST("/triage/folder", c.UploadTriageFolder)
		uploadGroup.POST("/triage/file", c.UploadTriageFile)
		uploadGroup.POST("/file", c.UploadFile)
		uploadGroup.DELETE("/file", c.UploadDeleteFile)
		uploadGroup.POST("/folder", c.UploadFolderCreate)
		uploadGroup.POST("/triage/mass-action", c.UploadMassImport)
		uploadGroup.DELETE("/triage/mass-action", c.UploadMassDelete)
	}

	// Adm
	admGroup.GET("/api/adm", c.AdmHome)
	admGroup.GET("/api/adm/video", c.AdmVideoList)
	admGroup.GET("/api/adm/actor", c.AdmActorList)
	admGroup.GET("/api/adm/channel", c.AdmChannelList)
	admGroup.GET("/api/adm/category", c.AdmCategory)
	admGroup.GET("/api/adm/config/auth", c.AdmConfigAuth)
	admGroup.GET("/api/adm/config/auth/:action", c.AdmConfigAuthUpdate)
	admGroup.GET("/api/adm/config/provider", c.AdmConfigProvider)
	admGroup.GET("/api/adm/config/provider/:id/switch", c.AdmConfigProviderSwitch)
	admGroup.GET("/api/adm/config/offline", c.AdmConfigOfflineMode)
	admGroup.GET("/api/adm/config/offline/:action", c.AdmConfigOfflineModeUpdate)
	admGroup.GET("/api/adm/task/home", c.AdmTaskHome)
	admGroup.GET("/api/adm/task", c.AdmTaskList)
	admGroup.GET("/api/adm/task/:id", c.AdmTaskView)
	admGroup.POST("/api/adm/task/:id/retry", c.AdmTaskRetry)
	admGroup.GET("/api/adm/user", c.AdmUserList)
	admGroup.POST("/api/adm/user", c.AdmUserNew)
	admGroup.DELETE("/api/adm/user/:id", c.AdmUserDelete)

	// Profile
	authGroup.GET("/api/profile", c.ProfileView)

	// Error
	authGroup.Any("/api/error/unauthorized", c.ErrUnauthorized)

	// NoRoute: serve SPA for GET (client-side routes) or JSON 404
	s.Router.NoRoute(c.NoRouteOrSPA)
}
