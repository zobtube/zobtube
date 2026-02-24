package controller

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/zobtube/zobtube/internal/model"
)

type profileActorViewResult struct {
	Actor model.Actor `json:"actor"`
	Count int         `json:"count"`
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

	// get user views
	var videoViewsTop []model.VideoView
	c.datastore.Where("user_id = ?", user.ID).Order("count desc").Limit(8).Preload("Video").Find(&videoViewsTop)

	// count actors
	countPerActor := make(map[string]int)
	var videoViewsAll []model.VideoView
	c.datastore.Where("user_id = ?", user.ID).Find(&videoViewsAll)

	type ActorResult struct { // create temporary type to hold actor ids
		ActorID string
	}
	for _, videoView := range videoViewsAll {
		var actors []ActorResult
		c.datastore.Table("video_actors").Select("actor_id").Where("video_id = ?", videoView.VideoID).Scan(&actors)
		for _, actor := range actors {
			countPerActor[actor.ActorID] += videoView.Count
		}
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
