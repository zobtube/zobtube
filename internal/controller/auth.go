package controller

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zobtube/zobtube/internal/model"
	"gorm.io/gorm"
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
	g.SetCookie(cookieName, session.ID, cookieMaxAge, "/", "127.0.0.1:8080", cookieSecure, cookieHttpOnly)
}

func (c *Controller) AuthPage(g *gin.Context) {
	_, err := g.Cookie(cookieName)
	if err != nil {
		c.createSession(g)
	}

	g.HTML(http.StatusOK, "auth/login.html", gin.H{})
}
func (c *Controller) AuthLogin(g *gin.Context) {
	// validate authentication
	cookie, err := g.Cookie(cookieName)
	if err != nil {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
	}

	// get session
	session := &model.UserSession{
		ID: cookie,
	}
	result := c.datastore.First(session)

	// check result
	if result.RowsAffected < 1 {
		c.createSession(g)

		g.JSON(401, gin.H{
			"error": "invalid session",
		})
		return
	}

	// check validity
	if session.ValidUntil.Before(time.Now()) {
		// session expired, creating a new one
		c.createSession(g)

		g.JSON(401, gin.H{
			"error": "session expired",
		})
		return
	}

	// retrieve user
	username := g.PostForm("username")
	user := &model.User{}
	result = c.datastore.First(&user, "username = ?", username)
	if result.RowsAffected < 1 {
		g.JSON(401, gin.H{
			"error": "auth failed - user not found",
		})
		return
	}

	// validate authentication
	challengeHex := sha256.Sum256([]byte(session.ID + user.Password))
	challenge := hex.EncodeToString(challengeHex[:])
	if g.PostForm("password") != challenge {
		g.JSON(401, gin.H{
			"error": "auth failed - password",
		})
		return
	}

	// extend expiration
	session.ValidUntil = time.Now().Add(sessionTimeValidated)
	session.UserID = &user.ID
	c.datastore.Save(session)

	// set auth cookie
	cookieMaxAge := int(sessionTimeValidated / time.Second)
	g.SetCookie(cookieName, session.ID, cookieMaxAge, "/", "127.0.0.1:8080", cookieSecure, cookieHttpOnly)

	g.JSON(200, gin.H{})
}

func (c *Controller) AuthLogout(g *gin.Context) {
	cookie, err := g.Cookie(cookieName)
	if err != nil {
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
	}

	// get session
	session := &model.UserSession{
		ID: cookie,
	}
	result := c.datastore.First(session)

	// check result
	if result.RowsAffected > 0 {
		// if found, delete it
		c.datastore.Delete(&session)
	}

	g.Redirect(http.StatusFound, "/auth")
}

func (c *Controller) GetSession(session *model.UserSession) *gorm.DB {
	return c.datastore.First(session)
}

func (c *Controller) GetUser(user *model.User) *gorm.DB {
	return c.datastore.First(user)
}
