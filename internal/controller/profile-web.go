package controller

import (
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/zobtube/zobtube/internal/model"
)

func (c *Controller) ProfileView(g *gin.Context) {
	// get user
	user := g.MustGet("user").(*model.User)

	// get user views
	var videoViewsTop []model.VideoView
	c.datastore.Where("user_id = ?", user.ID).Order("count desc").Limit(8).Preload("Video").Find(&videoViewsTop)

	// count actors
	var countPerActor = make(map[string]int)
	var videoViewsAll []model.VideoView
	c.datastore.Where("user_id = ?", user.ID).Find(&videoViewsAll) // get all views for user

	type ActorResult struct { // create temporary type to hold actor ids
		ActorID string
	}
	for _, videoView := range videoViewsAll {
		var actors []ActorResult
		c.datastore.Table("video_actors").Select("actor_id").Where("video_id = ?", videoView.VideoID).Scan(&actors)

		for _, actor := range actors {
			_, ok := countPerActor[actor.ActorID]
			if ok {
				countPerActor[actor.ActorID] = countPerActor[actor.ActorID] + videoView.Count
			} else {
				countPerActor[actor.ActorID] = videoView.Count
			}
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

	// create intermediate structure to hold count of actors
	type ActorView struct {
		Actor model.Actor
		Count int
	}

	var actorViews []ActorView
	actorLimit := 12

	for _, k := range keys {
		actorLimit--
		if actorLimit < 0 {
			break
		}

		actor := &model.Actor{
			ID: k,
		}
		c.datastore.First(actor)

		actorViews = append(actorViews, ActorView{
			Actor: *actor,
			Count: countPerActor[k],
		})
	}

	// render page
	g.HTML(http.StatusOK, "profile/view.html", gin.H{
		"User":       user,
		"VideoViews": videoViewsTop,
		"ActorViews": actorViews,
	})
}
