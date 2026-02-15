package controller

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zobtube/zobtube/internal/model"
)

// Bootstrap returns auth_enabled and current user (or nil) without requiring authentication.
func (c *Controller) Bootstrap(g *gin.Context) {
	user := &model.User{}
	if c.config.Authentication {
		if _, err := g.Cookie(cookieName); err != nil {
			c.createSession(g)
		} else {
			cookie, _ := g.Cookie(cookieName)
			session := &model.UserSession{ID: cookie}
			if c.datastore.First(session).RowsAffected > 0 &&
				session.ValidUntil.After(time.Now()) &&
				session.UserID != nil && *session.UserID != "" {
				u := &model.User{ID: *session.UserID}
				if c.datastore.First(u).RowsAffected > 0 {
					user = u
				}
			}
		}
	} else {
		_ = c.datastore.Order("created_at").First(user)
	}
	resp := gin.H{"auth_enabled": c.config.Authentication}
	if user.ID != "" {
		resp["user"] = gin.H{"id": user.ID, "username": user.Username, "admin": user.Admin}
	} else {
		resp["user"] = nil
	}
	g.JSON(http.StatusOK, resp)
}

func (c *Controller) SPAApp(g *gin.Context) {
	// Ensure session cookie exists for bootstrap
	if c.config.Authentication {
		if _, err := g.Cookie(cookieName); err != nil {
			c.createSession(g)
		}
	}
	g.HTML(http.StatusOK, "web/page/app-static.html", gin.H{})
}

// NoRouteOrSPA serves SPA for GET requests to non-API, non-static paths; otherwise returns JSON 404.
func (c *Controller) NoRouteOrSPA(g *gin.Context) {
	if g.Request.Method == http.MethodGet {
		path := g.Request.URL.Path
		if !strings.HasPrefix(path, "/api") && !strings.HasPrefix(path, "/static") && path != "/ping" {
			c.SPAApp(g)
			return
		}
	}
	c.ErrNotFound(g)
}
