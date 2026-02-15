package controller

import (
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/zobtube/zobtube/internal/model"
)

const (
	cookieName           = "zt_auth"
	cookieSecure         = false
	cookieHttpOnly       = false
	sessionTimePending   = 10 * time.Minute
	sessionTimeValidated = 24 * time.Hour
)

func (c *Controller) createSession(g *gin.Context) {
	// create a short session
	session := &model.UserSession{
		ValidUntil: time.Now().Add(sessionTimePending),
	}
	c.datastore.Save(session)

	// cookie not set, creating it
	cookieMaxAge := int(sessionTimePending / time.Second)
	g.SetCookie(cookieName, session.ID, cookieMaxAge, "/", "", cookieSecure, cookieHttpOnly)
}

func (c *Controller) GetSession(session *model.UserSession) *gorm.DB {
	return c.datastore.First(session)
}

func (c *Controller) GetUser(user *model.User) *gorm.DB {
	return c.datastore.First(user)
}

func (c *Controller) GetFirstUser(user *model.User) *gorm.DB {
	return c.datastore.Order("created_at").First(user)
}

func (c *Controller) AuthenticationEnabled() bool {
	return c.config.Authentication
}
