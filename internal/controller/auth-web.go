package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zobtube/zobtube/internal/model"
)

func (c *Controller) AuthPage(g *gin.Context) {
	_, err := g.Cookie(cookieName)
	if err != nil {
		c.createSession(g)
	}

	g.HTML(http.StatusOK, "auth/login.html", gin.H{})
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
