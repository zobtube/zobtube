package http

import (
	"io/fs"
	"net/http"

	"github.com/zobtube/zobtube/internal/controller"
)

func (s *Server) setupRoutes(c controller.AbtractController) {
	// load templates
	s.LoadHTMLFromEmbedFS("web/page/**/*")

	// prepare subfs
	staticFS, _ := fs.Sub(s.FS, "web/static")

	// load static
	s.Router.StaticFS("/static", http.FS(staticFS))
	s.Router.GET("/ping", livenessProbe)

	// authentication
	auth := s.Router.Group("/auth")
	auth.GET("", c.AuthPage)
	auth.POST("/login", c.AuthLogin)
	auth.GET("/logout", c.AuthLogout)

	authGroup := s.Router.Group("")
	authGroup.Use(UserIsAuthenticated(c))

	// home
	authGroup.GET("", c.Home)

	// actors
	authGroup.GET("/actors", c.ActorList)
	actors := authGroup.Group("/actor")
	actors.Use(UserIsAuthenticated(c))
	actors.GET("/new", c.ActorNew)
	actors.POST("/new", c.ActorNew)
	actors.GET("/:id", c.ActorView)
	actors.GET("/:id/edit", c.ActorEdit)
	actors.GET("/:id/thumb", c.ActorThumb)

	actorAPI := authGroup.Group("/api/actor")
	{
		actorAPI.POST("/", c.ActorAjaxNew)

		// providers
		actorAPI.GET("/:id/provider/:provider_slug", c.ActorAjaxProviderSearch)

		// links
		actorAPI.DELETE("/link/:id", c.ActorAjaxLinkThumbDelete)
		actorAPI.GET("/link/:id/thumb", c.ActorAjaxLinkThumbGet)

		// thumb
		actorAPI.POST("/:id/thumb", c.ActorAjaxThumb)
	}

	// channels
	authGroup.GET("/channels", c.ChannelList)
	channels := authGroup.Group("/channel")
	channels.GET("/new", c.ChannelCreate)
	channels.POST("/new", c.ChannelCreate)
	channels.GET("/:id", c.ChannelView)
	channels.GET("/:id/thumb", c.ChannelThumb)

	// videos
	authGroup.GET("/clips", c.ClipList)
	authGroup.GET("/movies", c.MovieList)
	authGroup.GET("/videos", c.VideoList)
	videos := authGroup.Group("/video")
	videos.GET("/:id", c.VideoView)
	videos.GET("/:id/edit", c.VideoEdit)
	videos.GET("/:id/stream", c.VideoStream)
	videos.GET("/:id/thumb", c.VideoThumb)
	videos.GET("/:id/thumb_xs", c.VideoThumbXS)

	videoAPI := authGroup.Group("/api/video")
	videoAPI.POST("", c.VideoAjaxCreate)
	videoAPI.HEAD("/:id", c.VideoAjaxStreamInfo)
	videoAPI.DELETE("/:id", c.VideoAjaxDelete)
	videoAPI.POST("/:id/upload", c.VideoAjaxUpload)
	videoAPI.POST("/:id/thumb", c.VideoAjaxUploadThumb)
	videoAPI.POST("/:id/migrate", c.VideoAjaxMigrate)
	videoAPI.PUT("/:id/actor/:actor_id", c.VideoAjaxActors)
	videoAPI.DELETE("/:id/actor/:actor_id", c.VideoAjaxActors)
	videoAPI.POST("/:id/compute-duration", c.VideoAjaxComputeDuration)
	videoAPI.POST("/:id/generate-thumbnail/:timing", c.VideoAjaxGenerateThumbnail)
	videoAPI.POST("/:id/generate-thumbnail-xs", c.VideoAjaxGenerateThumbnailXS)
	videoAPI.POST("/:id/import", c.VideoAjaxImport)
	videoAPI.POST("/:id/rename", c.VideoAjaxRename)
	videoAPI.POST("/:id/count-view", c.VideoViewAjaxIncrement)

	// uploads
	uploads := authGroup.Group("/upload")
	uploads.GET("/", c.UploadTriage)
	uploads.GET("/preview/:filepath", c.UploadPreview)
	uploads.POST("/import", c.UploadImport)
	uploadAPI := authGroup.Group("/api/upload")
	uploadAPI.POST("/triage/folder", c.UploadAjaxTriageFolder)
	uploadAPI.POST("/triage/file", c.UploadAjaxTriageFile)
	uploadAPI.POST("/file", c.UploadAjaxUploadFile)

	// adm
	authGroup.GET("/adm", c.AdmHome)

	// profile
	authGroup.GET("/profile", c.ProfileView)

	// remainings routes to implement
	/*
	   path('actor/<uuid:id>/edit/first-time', views.actor_edit, name='actor_edit_first_time', kwargs={'first_time': True}),
	   path('actor/<uuid:id>/remove', views.actor_remove, name='actor_remove'),

	   path('channels', views.ChannelListView.as_view(), name='channel_list'),
	   path('channel/new', views.channel_new, name='channel_new'),
	   path('channel/<uuid:id>', views.channel_view, name='channel_view'),
	   path('channel/<uuid:id>/thumb', views.channel_thumb, name='channel_thumb'),

	   path('profile', views.profile_view, name='profile_view'),

	   path('triage/delete/<path:name>', views.triage_delete, name='triage_delete'),
	   path('uploads', views.upload_home, name='upload_home'),
	   path('upload/list', views.upload_list, name='upload_list'),
	   path('upload/<uuid:pk>/stream', views.upload_stream, name='upload_stream'),
	   path('upload/<uuid:pk>/delete', views.upload_delete, name='upload_delete'),
	   path('upload/<uuid:pk>/import/<str:import_as>', views.upload_import, name='upload_import'),
	   path('upload/new', views.upload_new, name='upload_new'),
	   path('upload/file', views.ChunkedUploadView.as_view(), name='upload_file'),
	   path('upload/file/<uuid:pk>', views.ChunkedUploadView.as_view(), name='upload_file_view'),
	   path('adm/actor/list', views.adm_actor_list, name='adm_actor_list'),
	   path('adm/actor/fix-thumb', views.adm_actor_fix_missing_thumb, name='adm_actor_fix_missing_thumb'),
	   path('adm/actor/<uuid:id>/fix-thumb', views.adm_actor_fix_thumb, name='adm_actor_fix_thumb'),
	   path('adm/actor/<uuid:id>/gen-thumb', views.adm_actor_gen_thumb, name='adm_actor_gen_thumb'),
	*/
}
