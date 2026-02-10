package controller

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zobtube/zobtube/internal/model"
)

func (c *Controller) AuthPage(g *gin.Context) {
	if !c.config.Authentication {
		g.Redirect(http.StatusFound, "/")
		return
	}

	_, err := g.Cookie(cookieName)
	if err != nil {
		c.createSession(g)
	}

	next := g.Query("next")
	if next != "" && (!strings.HasPrefix(next, "/") || strings.HasPrefix(next, "//")) {
		next = ""
	}

	g.HTML(http.StatusOK, "auth/login.html", gin.H{"next": next})
}

func (c *Controller) AuthLogout(g *gin.Context) {
	cookie, err := g.Cookie(cookieName)
	if err != nil {
		c.ErrFatal(g, err.Error())
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

	// clear the auth cookie so the next visit to /auth gets a fresh session
	g.SetCookie(cookieName, "", -1, "/", "127.0.0.1:8069", cookieSecure, cookieHttpOnly)
	g.Redirect(http.StatusFound, "/auth")
}
