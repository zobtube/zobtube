package controller

import (
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/zobtube/zobtube/internal/model"
)

type profileActorViewResult struct {
	Actor model.Actor `json:"actor"`
	Count int         `json:"count"`
}

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
