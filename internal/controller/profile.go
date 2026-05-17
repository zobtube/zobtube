package controller

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/zobtube/zobtube/internal/model"
)

type profileActorViewResult struct {
	Actor model.Actor `json:"actor"`
	Count int         `json:"count"`
}

type profileStatsResult struct {
	VideosUnique    int   `json:"videos_unique"`
	VideosTotal     int   `json:"videos_total"`
	ActorsUnique    int   `json:"actors_unique"`
	ActorsTotal     int   `json:"actors_total"`
	TotalViewTimeNs int64 `json:"total_view_time_ns"`
}

// migrateVideoViewsFromDeletedVideos moves view counts from soft-deleted videos onto
// a still-active video with the same filename in the same library (e.g. after
// re-importing triage content as a different type).
func (c *Controller) migrateVideoViewsFromDeletedVideos(userID string) {
	var views []model.VideoView
	c.datastore.Where("user_id = ?", userID).Find(&views)
	for _, vv := range views {
		var active model.Video
		if c.datastore.First(&active, "id = ?", vv.VideoID).RowsAffected > 0 {
			continue
		}
		var deleted model.Video
		if c.datastore.Unscoped().First(&deleted, "id = ?", vv.VideoID).Error != nil || deleted.ID == "" {
			c.datastore.Delete(&model.VideoView{}, "video_id = ? AND user_id = ?", vv.VideoID, userID)
			continue
		}
		replacement := c.findReplacementVideo(&deleted)
		if replacement == nil {
			continue
		}
		c.mergeVideoView(vv.VideoID, replacement.ID, userID, vv.Count)
	}
}

func (c *Controller) findReplacementVideo(deleted *model.Video) *model.Video {
	if deleted.Filename == "" {
		return nil
	}
	q := c.datastore.Where("filename = ? AND deleted_at IS NULL AND id != ?", deleted.Filename, deleted.ID)
	if deleted.LibraryID != nil && *deleted.LibraryID != "" {
		q = q.Where("library_id = ?", *deleted.LibraryID)
	}
	var replacement model.Video
	if q.Order("created_at desc").First(&replacement).Error != nil || replacement.ID == "" {
		return nil
	}
	return &replacement
}

func (c *Controller) mergeVideoView(fromVideoID, toVideoID, userID string, count int) {
	if fromVideoID == "" || toVideoID == "" || fromVideoID == toVideoID || count <= 0 {
		return
	}
	var target model.VideoView
	if c.datastore.First(&target, "video_id = ? AND user_id = ?", toVideoID, userID).RowsAffected > 0 {
		target.Count += count
		_ = c.datastore.Save(&target).Error
	} else {
		_ = c.datastore.Create(&model.VideoView{VideoID: toVideoID, UserID: userID, Count: count}).Error
	}
	c.datastore.Delete(&model.VideoView{}, "video_id = ? AND user_id = ?", fromVideoID, userID)
}

func profileVideoPreload(db *gorm.DB) *gorm.DB {
	return db.Unscoped()
}

// resolveProfileVideo returns the video to attribute views to, or nil if the view row is orphaned.
func (c *Controller) resolveProfileVideo(vv *model.VideoView, userID string) *model.Video {
	if vv.Video.ID == "" {
		c.datastore.Delete(&model.VideoView{}, "video_id = ? AND user_id = ?", vv.VideoID, userID)
		return nil
	}
	video := vv.Video
	if video.DeletedAt.Valid {
		if replacement := c.findReplacementVideo(&video); replacement != nil {
			return replacement
		}
	}
	return &video
}

// ProfileView godoc
//
//	@Summary	Get user profile with top video views and actor stats
//	@Tags		profile
//	@Produce	json
//	@Success	200	{object}	map[string]interface{}
//	@Router		/profile [get]
func (c *Controller) ProfileView(g *gin.Context) {
	// get user
	user := g.MustGet("user").(*model.User)

	c.migrateVideoViewsFromDeletedVideos(user.ID)

	var videoViewsAll []model.VideoView
	c.datastore.Where("user_id = ?", user.ID).
		Preload("Video", profileVideoPreload).Find(&videoViewsAll)

	type actorIDResult struct {
		ActorID string
	}
	countPerActor := make(map[string]int)
	actorsSeen := make(map[string]struct{})
	var stats profileStatsResult
	validTop := make([]model.VideoView, 0, len(videoViewsAll))

	for _, vv := range videoViewsAll {
		if vv.Count <= 0 {
			continue
		}
		video := c.resolveProfileVideo(&vv, user.ID)
		if video == nil {
			continue
		}
		stats.VideosUnique++
		stats.VideosTotal += vv.Count
		stats.TotalViewTimeNs += int64(vv.Count) * int64(video.Duration)

		display := vv
		display.Video = *video
		validTop = append(validTop, display)

		var actors []actorIDResult
		c.datastore.Table("video_actors").Select("actor_id").Where("video_id = ?", vv.VideoID).Scan(&actors)
		for _, actor := range actors {
			if actor.ActorID == "" {
				continue
			}
			countPerActor[actor.ActorID] += vv.Count
			actorsSeen[actor.ActorID] = struct{}{}
		}
	}

	stats.ActorsUnique = len(actorsSeen)
	for _, n := range countPerActor {
		stats.ActorsTotal += n
	}

	sort.Slice(validTop, func(i, j int) bool {
		return validTop[i].Count > validTop[j].Count
	})
	videoViewsTop := validTop
	if len(videoViewsTop) > 8 {
		videoViewsTop = videoViewsTop[:8]
	}

	// sort actors
	keys := make([]string, 0, len(countPerActor))
	for k := range countPerActor {
		keys = append(keys, k)
	}
	sort.SliceStable(keys, func(i, j int) bool {
		return countPerActor[keys[i]] > countPerActor[keys[j]]
	})

	var actorViews []profileActorViewResult
	actorLimit := 12
	for _, k := range keys {
		if actorLimit <= 0 {
			break
		}
		actorLimit--
		actor := &model.Actor{ID: k}
		c.datastore.First(actor)
		actorViews = append(actorViews, profileActorViewResult{Actor: *actor, Count: countPerActor[k]})
	}

	g.JSON(http.StatusOK, gin.H{
		"video_views": videoViewsTop,
		"actor_views": actorViews,
		"stats":       stats,
	})
}

// ProfileChangePassword godoc
//
//	@Summary	Change password for the authenticated user
//	@Tags		profile
//	@Accept		json
//	@Param		body	body	object	true	"JSON with current_password, new_password"
//	@Success	200	{object}	map[string]interface{}
//	@Failure	400	{object}	map[string]interface{}
//	@Router		/profile/password [post]
func (c *Controller) ProfileChangePassword(g *gin.Context) {
	user := g.MustGet("user").(*model.User)
	var body struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}
	if err := g.ShouldBindJSON(&body); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	currentSum := sha256.Sum256([]byte(body.CurrentPassword))
	currentHash := hex.EncodeToString(currentSum[:])
	if currentHash != user.Password {
		g.JSON(http.StatusBadRequest, gin.H{"error": "wrong current password"})
		return
	}
	if body.NewPassword == "" {
		g.JSON(http.StatusBadRequest, gin.H{"error": "new password cannot be empty"})
		return
	}
	newSum := sha256.Sum256([]byte(body.NewPassword))
	user.Password = hex.EncodeToString(newSum[:])
	if err := c.datastore.Save(user).Error; err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save"})
		return
	}
	g.JSON(http.StatusOK, gin.H{"ok": true})
}
