package http

import (
	"net/http"
	"net/url"
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

func UserIsAuthenticated(c controller.AbstractController) gin.HandlerFunc {
	return func(g *gin.Context) {
		if !c.AuthenticationEnabled() {
			// get user
			user := &model.User{}
			result := c.GetFirstUser(user)
			if result.RowsAffected < 1 {
				// something bad happened
				// most likely someone deleted the last user
				// restart the application to re-create the default user
				c.Restart()
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
			// cookie not set
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
			g.Redirect(http.StatusFound, authRedirectURL(g))
			g.Abort()
			return
		}

		// check if user is authenticated
		if session.UserID == nil || *session.UserID == "" {
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
			g.Redirect(http.StatusFound, authRedirectURL(g))
			g.Abort()
			return
		}

		// set meta in context
		g.Set("user", user)

		g.Next()
	}
}
