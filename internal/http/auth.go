package http

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/zobtube/zobtube/internal/controller"
	"github.com/zobtube/zobtube/internal/model"
)

const cookieName = "zt_auth"

func authRedirectURL(g *gin.Context) string {
	nextVal := g.Request.URL.Path
	if g.Request.URL.RawQuery != "" {
		nextVal += "?" + g.Request.URL.RawQuery
	}
	return "/auth?next=" + url.QueryEscape(nextVal)
}

func apiUnauthorized(g *gin.Context) {
	g.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	g.Abort()
}

func UserIsAuthenticated(c controller.AbstractController) gin.HandlerFunc {
	return func(g *gin.Context) {
		isAPI := strings.HasPrefix(g.Request.URL.Path, "/api")

		if !c.AuthenticationEnabled() {
			// get user
			user := &model.User{}
			result := c.GetFirstUser(user)
			if result.RowsAffected < 1 {
				if isAPI {
					g.JSON(http.StatusInternalServerError, gin.H{"error": "no user"})
					g.Abort()
					return
				}
				g.Redirect(http.StatusFound, "/")
				g.Abort()
				return
			}

			// set meta in context
			g.Set("user", user)

			// all good, exiting middleware
			g.Next()
			return
		}

		cookie, err := g.Cookie(cookieName)
		if err != nil {
			if isAPI {
				apiUnauthorized(g)
				return
			}
			g.Redirect(http.StatusFound, authRedirectURL(g))
			g.Abort()
			return
		}

		// get session
		session := &model.UserSession{
			ID: cookie,
		}
		result := c.GetSession(session)

		// check result
		if result.RowsAffected < 1 {
			g.Redirect(http.StatusFound, authRedirectURL(g))
			g.Abort()
			return
		}

		// check validity
		if session.ValidUntil.Before(time.Now()) {
			if isAPI {
				apiUnauthorized(g)
				return
			}
			g.Redirect(http.StatusFound, authRedirectURL(g))
			g.Abort()
			return
		}

		// check if user is authenticated
		if session.UserID == nil || *session.UserID == "" {
			if isAPI {
				apiUnauthorized(g)
				return
			}
			g.Redirect(http.StatusFound, authRedirectURL(g))
			g.Abort()
			return
		}

		// get user
		user := &model.User{
			ID: *session.UserID,
		}
		result = c.GetUser(user)
		if result.RowsAffected < 1 {
			if isAPI {
				apiUnauthorized(g)
				return
			}
			g.Redirect(http.StatusFound, authRedirectURL(g))
			g.Abort()
			return
		}

		// set meta in context
		g.Set("user", user)
		g.Next()
	}
}
