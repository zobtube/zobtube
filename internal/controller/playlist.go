package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/zobtube/zobtube/internal/model"
)

const (
	PlaylistUnseenVideosID = "unseen-v"
	PlaylistUnseenClipsID  = "unseen-c"
	PlaylistUnseenMoviesID = "unseen-m"
)

type virtualPlaylistDef struct {
	id, name, videoType string
}

var virtualPlaylists = []virtualPlaylistDef{
	{PlaylistUnseenVideosID, "Unseen videos", "v"},
	{PlaylistUnseenClipsID, "Unseen clips", "c"},
	{PlaylistUnseenMoviesID, "Unseen movies", "m"},
}

func playlistIsVirtual(id string) (virtual bool, name string, videoType string) {
	for _, vp := range virtualPlaylists {
		if vp.id == id {
			return true, vp.name, vp.videoType
		}
	}
	return false, "", ""
}

func unseenVideosQuery(db *gorm.DB, userID, videoType string) *gorm.DB {
	return db.Model(&model.Video{}).
		Joins("LEFT JOIN video_views ON video_views.video_id = videos.id AND video_views.user_id = ?", userID).
		Where("video_views.video_id IS NULL").
		Where("videos.type = ?", videoType).
		Where("videos.status = ?", model.VideoStatusReady)
}

func unseenVideosForUser(db *gorm.DB, userID, videoType string) []model.Video {
	var videos []model.Video
	unseenVideosQuery(db, userID, videoType).Order("videos.created_at desc").Find(&videos)
	return videos
}

func unseenVideoCount(db *gorm.DB, userID, videoType string) int64 {
	var count int64
	unseenVideosQuery(db, userID, videoType).Count(&count)
	return count
}

func playlistPlaybackContextFromVideos(playlistID, playlistName string, videos []model.Video, videoID string) gin.H {
	index := -1
	for i, v := range videos {
		if v.ID == videoID {
			index = i
			break
		}
	}
	if index < 0 {
		return nil
	}
	ids := make([]string, len(videos))
	for i, v := range videos {
		ids[i] = v.ID
	}
	var upNext []model.Video
	if index+1 < len(videos) {
		upNext = videos[index+1:]
	}
	items := make([]gin.H, len(videos))
	for i, v := range videos {
		items[i] = gin.H{"id": v.ID, "type": v.Type}
	}
	return gin.H{
		"playlist": gin.H{
			"id":   playlistID,
			"name": playlistName,
		},
		"playlist_index":     index,
		"playlist_video_ids": ids,
		"playlist_items":     items,
		"playlist_videos":    videos,
		"playlist_up_next":   upNext,
	}
}

func (c *Controller) virtualPlaylistListEntries(userID, videoID string) []gin.H {
	now := time.Now()
	var targetType string
	if videoID != "" {
		var v model.Video
		if c.datastore.First(&v, "id = ?", videoID).RowsAffected > 0 {
			targetType = v.Type
		}
	}
	out := make([]gin.H, 0, len(virtualPlaylists))
	for _, vp := range virtualPlaylists {
		count := unseenVideoCount(c.datastore, userID, vp.videoType)
		item := gin.H{
			"id":          vp.id,
			"name":        vp.name,
			"video_count": count,
			"virtual":     true,
			"deletable":   false,
			"updated_at":  now,
			"created_at":  now,
		}
		if videoID != "" {
			contains := false
			if targetType == vp.videoType {
				var n int64
				unseenVideosQuery(c.datastore, userID, vp.videoType).Where("videos.id = ?", videoID).Count(&n)
				contains = n > 0
			}
			item["contains"] = contains
		}
		out = append(out, item)
	}
	return out
}

func (c *Controller) playlistOwnedByUser(id, userID string) (*model.Playlist, bool) {
	var playlist model.Playlist
	if c.datastore.Where("id = ? AND user_id = ?", id, userID).First(&playlist).RowsAffected < 1 {
		return nil, false
	}
	return &playlist, true
}

func playlistVideoCount(db *gorm.DB, playlistID string) int64 {
	var count int64
	db.Model(&model.PlaylistVideo{}).Where("playlist_id = ?", playlistID).Count(&count)
	return count
}

func playlistOrderedVideos(db *gorm.DB, playlistID string) []model.Video {
	var entries []model.PlaylistVideo
	db.Where("playlist_id = ?", playlistID).Order("position asc, added_at asc").Find(&entries)
	if len(entries) == 0 {
		return nil
	}
	videoIDs := make([]string, 0, len(entries))
	for _, e := range entries {
		videoIDs = append(videoIDs, e.VideoID)
	}
	videosByID := map[string]model.Video{}
	var videos []model.Video
	db.Where("id IN ?", videoIDs).Find(&videos)
	for _, v := range videos {
		videosByID[v.ID] = v
	}
	ordered := make([]model.Video, 0, len(entries))
	for _, e := range entries {
		if v, found := videosByID[e.VideoID]; found {
			ordered = append(ordered, v)
		}
	}
	return ordered
}

// playlistPlaybackContext returns playback queue fields when videoID is in the user's playlist.
// Returns nil if the playlist is not owned or the video is not in the playlist.
func (c *Controller) playlistPlaybackContext(userID, playlistID, videoID string) gin.H {
	if virtual, name, videoType := playlistIsVirtual(playlistID); virtual {
		videos := unseenVideosForUser(c.datastore, userID, videoType)
		return playlistPlaybackContextFromVideos(playlistID, name, videos, videoID)
	}
	playlist, ok := c.playlistOwnedByUser(playlistID, userID)
	if !ok {
		return nil
	}
	videos := playlistOrderedVideos(c.datastore, playlist.ID)
	return playlistPlaybackContextFromVideos(playlist.ID, playlist.Name, videos, videoID)
}

// PlaylistList returns playlists for the current user.
// Optional query ?video_id= adds contains: bool per playlist for picker UI.
func (c *Controller) PlaylistList(g *gin.Context) {
	user := g.MustGet("user").(*model.User)
	videoID := g.Query("video_id")

	var playlists []model.Playlist
	c.datastore.Where("user_id = ?", user.ID).Order("updated_at desc").Find(&playlists)

	containsSet := map[string]bool{}
	if videoID != "" {
		var rows []model.PlaylistVideo
		c.datastore.Where("video_id = ?", videoID).Find(&rows)
		for _, r := range rows {
			containsSet[r.PlaylistID] = true
		}
	}

	list := c.virtualPlaylistListEntries(user.ID, videoID)
	for _, p := range playlists {
		item := gin.H{
			"id":          p.ID,
			"name":        p.Name,
			"video_count": playlistVideoCount(c.datastore, p.ID),
			"updated_at":  p.UpdatedAt,
			"created_at":  p.CreatedAt,
		}
		if videoID != "" {
			item["contains"] = containsSet[p.ID]
		}
		list = append(list, item)
	}
	g.JSON(http.StatusOK, gin.H{"playlists": list})
}

// PlaylistCreate creates a playlist for the current user.
func (c *Controller) PlaylistCreate(g *gin.Context) {
	user := g.MustGet("user").(*model.User)
	var body struct {
		Name string `json:"name"`
	}
	if err := g.ShouldBindJSON(&body); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	name := body.Name
	if name == "" {
		g.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}
	now := time.Now()
	playlist := &model.Playlist{
		UserID:    user.ID,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := c.datastore.Create(playlist).Error; err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save"})
		return
	}
	g.JSON(http.StatusCreated, gin.H{
		"id":          playlist.ID,
		"name":        playlist.Name,
		"video_count": 0,
		"created_at":  playlist.CreatedAt,
		"updated_at":  playlist.UpdatedAt,
	})
}

// PlaylistView returns a playlist with ordered videos if owned by the current user.
func (c *Controller) PlaylistView(g *gin.Context) {
	user := g.MustGet("user").(*model.User)
	id := g.Param("id")
	if virtual, name, videoType := playlistIsVirtual(id); virtual {
		videos := unseenVideosForUser(c.datastore, user.ID, videoType)
		now := time.Now()
		g.JSON(http.StatusOK, gin.H{
			"playlist": gin.H{
				"id":          id,
				"name":        name,
				"video_count": len(videos),
				"virtual":     true,
				"deletable":   false,
				"created_at":  now,
				"updated_at":  now,
			},
			"videos": videos,
		})
		return
	}
	playlist, ok := c.playlistOwnedByUser(id, user.ID)
	if !ok {
		g.JSON(http.StatusNotFound, gin.H{"error": "playlist not found"})
		return
	}

	videos := playlistOrderedVideos(c.datastore, playlist.ID)

	g.JSON(http.StatusOK, gin.H{
		"playlist": gin.H{
			"id":          playlist.ID,
			"name":        playlist.Name,
			"video_count": len(videos),
			"created_at":  playlist.CreatedAt,
			"updated_at":  playlist.UpdatedAt,
		},
		"videos": videos,
	})
}

// PlaylistUpdate renames a playlist owned by the current user.
func (c *Controller) PlaylistUpdate(g *gin.Context) {
	user := g.MustGet("user").(*model.User)
	id := g.Param("id")
	if virtual, _, _ := playlistIsVirtual(id); virtual {
		g.JSON(http.StatusForbidden, gin.H{"error": "cannot modify automatic playlist"})
		return
	}
	playlist, ok := c.playlistOwnedByUser(id, user.ID)
	if !ok {
		g.JSON(http.StatusNotFound, gin.H{"error": "playlist not found"})
		return
	}
	var body struct {
		Name string `json:"name"`
	}
	if err := g.ShouldBindJSON(&body); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if body.Name == "" {
		g.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}
	playlist.Name = body.Name
	playlist.UpdatedAt = time.Now()
	if err := c.datastore.Save(playlist).Error; err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save"})
		return
	}
	g.JSON(http.StatusOK, gin.H{
		"id":          playlist.ID,
		"name":        playlist.Name,
		"video_count": playlistVideoCount(c.datastore, playlist.ID),
		"updated_at":  playlist.UpdatedAt,
	})
}

// PlaylistDelete deletes a playlist and its videos if owned by the current user.
func (c *Controller) PlaylistDelete(g *gin.Context) {
	user := g.MustGet("user").(*model.User)
	id := g.Param("id")
	if virtual, _, _ := playlistIsVirtual(id); virtual {
		g.JSON(http.StatusForbidden, gin.H{"error": "cannot delete automatic playlist"})
		return
	}
	playlist, ok := c.playlistOwnedByUser(id, user.ID)
	if !ok {
		g.JSON(http.StatusNotFound, gin.H{"error": "playlist not found"})
		return
	}
	if err := c.datastore.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("playlist_id = ?", playlist.ID).Delete(&model.PlaylistVideo{}).Error; err != nil {
			return err
		}
		return tx.Delete(playlist).Error
	}); err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete"})
		return
	}
	g.AbortWithStatus(http.StatusNoContent)
}

// PlaylistVideoAdd adds a video to a playlist owned by the current user.
func (c *Controller) PlaylistVideoAdd(g *gin.Context) {
	user := g.MustGet("user").(*model.User)
	id := g.Param("id")
	if virtual, _, _ := playlistIsVirtual(id); virtual {
		g.JSON(http.StatusForbidden, gin.H{"error": "cannot modify automatic playlist"})
		return
	}
	playlist, ok := c.playlistOwnedByUser(id, user.ID)
	if !ok {
		g.JSON(http.StatusNotFound, gin.H{"error": "playlist not found"})
		return
	}
	var body struct {
		VideoID string `json:"video_id"`
	}
	if err := g.ShouldBindJSON(&body); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if body.VideoID == "" {
		g.JSON(http.StatusBadRequest, gin.H{"error": "video_id is required"})
		return
	}
	var video model.Video
	if c.datastore.First(&video, "id = ?", body.VideoID).RowsAffected < 1 {
		g.JSON(http.StatusNotFound, gin.H{"error": "video not found"})
		return
	}

	var existing model.PlaylistVideo
	if c.datastore.Where("playlist_id = ? AND video_id = ?", playlist.ID, body.VideoID).First(&existing).RowsAffected > 0 {
		g.JSON(http.StatusOK, gin.H{"added": false, "already_present": true})
		return
	}

	var maxPos int
	c.datastore.Model(&model.PlaylistVideo{}).Where("playlist_id = ?", playlist.ID).
		Select("COALESCE(MAX(position), -1)").Scan(&maxPos)

	now := time.Now()
	entry := &model.PlaylistVideo{
		PlaylistID: playlist.ID,
		VideoID:    body.VideoID,
		Position:   maxPos + 1,
		AddedAt:    now,
	}
	if err := c.datastore.Create(entry).Error; err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add video"})
		return
	}
	playlist.UpdatedAt = now
	c.datastore.Save(playlist)

	g.JSON(http.StatusOK, gin.H{"added": true})
}

// PlaylistVideoRemove removes a video from a playlist owned by the current user.
func (c *Controller) PlaylistVideoRemove(g *gin.Context) {
	user := g.MustGet("user").(*model.User)
	id := g.Param("id")
	if virtual, _, _ := playlistIsVirtual(id); virtual {
		g.JSON(http.StatusForbidden, gin.H{"error": "cannot modify automatic playlist"})
		return
	}
	videoID := g.Param("video_id")
	playlist, ok := c.playlistOwnedByUser(id, user.ID)
	if !ok {
		g.JSON(http.StatusNotFound, gin.H{"error": "playlist not found"})
		return
	}
	if videoID == "" {
		g.JSON(http.StatusBadRequest, gin.H{"error": "video_id required"})
		return
	}
	result := c.datastore.Where("playlist_id = ? AND video_id = ?", playlist.ID, videoID).Delete(&model.PlaylistVideo{})
	if result.RowsAffected < 1 {
		g.JSON(http.StatusNotFound, gin.H{"error": "video not in playlist"})
		return
	}
	playlist.UpdatedAt = time.Now()
	c.datastore.Save(playlist)
	g.AbortWithStatus(http.StatusNoContent)
}
