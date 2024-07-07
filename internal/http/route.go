package http

import (
	"gitlab.com/zobtube/zobtube/internal/controller"
)

func (s *Server) setupRoutes(c controller.AbtractController) {
	// load templates
	s.Server.LoadHTMLGlob("web/page/**/*")

	// load static
	s.Server.Static("/static", "web/static")
	s.Server.GET("/ping", livenessProbe)

	// home
	s.Server.GET("", c.Home)

	// actors
	s.Server.GET("/actors", c.ActorList)
	actors := s.Server.Group("/actor")
	actors.GET("/new", c.ActorNew)
	actors.POST("/new", c.ActorNew)
	actors.GET("/:id", c.ActorView)
	actors.GET("/:id/edit", c.ActorEdit)
	actors.GET("/:id/thumb", c.ActorThumb)

	actorAPI := s.Server.Group("/api/actor")
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
	s.Server.GET("/channels", c.ChannelList)
	channels := s.Server.Group("/channel")
	channels.GET("/new", c.ChannelCreate)
	channels.POST("/new", c.ChannelCreate)
	channels.GET("/:id", c.ChannelView)
	channels.GET("/:id/thumb", c.ChannelThumb)

	// clips
	s.Server.GET("/clips", c.ClipList)
	clips := s.Server.Group("/clip")
	clips.GET("/:id", c.ClipView)
	clips.GET("/:id/stream", c.ClipStream)
	clips.GET("/:id/thumb", c.ClipThumb)
	clips.GET("/:id/thumb_xs", c.ClipThumbXS)

	// movies
	s.Server.GET("/movies", c.MovieList)
	movies := s.Server.Group("/movie")
	movies.GET("/:id", c.MovieView)
	movies.GET("/:id/stream", c.MovieStream)
	movies.GET("/:id/thumb", c.MovieThumb)
	movies.GET("/:id/thumb_xs", c.MovieThumbXS)

	// videos
	s.Server.GET("/videos", c.VideoList)
	videos := s.Server.Group("/video")
	videos.GET("/:id", c.VideoView)
	videos.GET("/:id/edit", c.GenericVideoEdit)
	videos.GET("/:id/stream", c.VideoStream)
	videos.GET("/:id/thumb", c.VideoThumb)
	videos.GET("/:id/thumb_xs", c.VideoThumbXS)

	videoAPI := s.Server.Group("/api/video")
	videoAPI.POST("/", c.GenericVideoAjaxCreate)
	videoAPI.HEAD("/:id", c.GenericVideoAjaxStreamInfo)
	videoAPI.POST("/:id/upload", c.GenericVideoAjaxUpload)
	videoAPI.POST("/:id/thumb", c.GenericVideoAjaxUploadThumb)
	videoAPI.PUT("/:id/actor/:actor_id", c.GenericVideoAjaxActors)
	videoAPI.DELETE("/:id/actor/:actor_id", c.GenericVideoAjaxActors)
	videoAPI.POST("/:id/compute-duration", c.GenericVideoAjaxComputeDuration)
	videoAPI.POST("/:id/generate-thumbnail/:timing", c.GenericVideoAjaxGenerateThumbnail)
	videoAPI.POST("/:id/generate-thumbnail-xs", c.GenericVideoAjaxGenerateThumbnailXS)
	videoAPI.POST("/:id/import", c.GenericVideoAjaxImport)
	videoAPI.POST("/:id/rename", c.GenericVideoAjaxRename)

	// uploads
	uploads := s.Server.Group("/upload")
	uploads.GET("/", c.UploadHome)
	uploads.GET("/triage", c.UploadTriage)
	uploads.GET("/preview/:filepath", c.UploadPreview)
	uploads.POST("/import", c.UploadImport)

	// adm
	s.Server.GET("/adm", c.AdmHome)

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
