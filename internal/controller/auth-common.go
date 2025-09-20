package controller

import (
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/zobtube/zobtube/internal/model"
)

const cookieName = "zt_auth"
const cookieSecure = false
const cookieHttpOnly = false
const sessionTimePending = 10 * time.Minute
const sessionTimeValidated = 24 * time.Hour

func (c *Controller) createSession(g *gin.Context) {
	// create a short session
	session := &model.UserSession{
		ValidUntil: time.Now().Add(sessionTimePending),
	}
	c.datastore.Save(session)

	// cookie not set, creating it
	cookieMaxAge := int(sessionTimePending / time.Second)
	g.SetCookie(cookieName, session.ID, cookieMaxAge, "/", "127.0.0.1:8069", cookieSecure, cookieHttpOnly)
}

func (c *Controller) GetSession(session *model.UserSession) *gorm.DB {
	return c.datastore.First(session)
}

func (c *Controller) GetUser(user *model.User) *gorm.DB {
	return c.datastore.First(user)
}
