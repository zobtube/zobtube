package controller

import (
	"log"
	"time"

	"github.com/zobtube/zobtube/internal/model"
)

func (c *Controller) CleanupRoutine() {
	for {
		// wait loop
		time.Sleep(time.Minute)

		c.sessionCleanup()
	}
}

func (c *Controller) sessionCleanup() {
	var sessions []model.UserSession
	result := c.datastore.Find(&sessions)
	err := result.Error
	if err != nil {
		log.Println("error while querying sessions")
		log.Println(err.Error())
	}

	for _, session := range sessions {
		if session.ValidUntil.Before(time.Now()) {
			c.datastore.Delete(&session)
		}
	}
}
