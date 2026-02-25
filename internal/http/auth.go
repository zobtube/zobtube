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

// tryCookieAuth attempts to authenticate via session cookie. Returns the user if successful, nil otherwise.
func tryCookieAuth(g *gin.Context, c controller.AbstractController) *model.User {
	cookie, err := g.Cookie(cookieName)
	if err != nil {
		return nil
	}
	session := &model.UserSession{ID: cookie}
	result := c.GetSession(session)
	if result.RowsAffected < 1 {
		return nil
	}
	if session.ValidUntil.Before(time.Now()) {
		return nil
	}
	if session.UserID == nil || *session.UserID == "" {
		return nil
	}
	user := &model.User{ID: *session.UserID}
	result = c.GetUser(user)
	if result.RowsAffected < 1 {
		return nil
	}
	return user
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

		user := tryCookieAuth(g, c)
		if user != nil {
			g.Set("user", user)
			g.Next()
			return
		}

		// For API requests, try Bearer token
		if isAPI {
			authz := g.GetHeader("Authorization")
			if strings.HasPrefix(authz, "Bearer ") {
				token := strings.TrimSpace(authz[7:])
				if token != "" {
					if bearerUser, ok := c.ResolveUserByApiTokenHash(token); ok {
						g.Set("user", bearerUser)
						g.Next()
						return
					}
				}
			}
		}

		if isAPI {
			apiUnauthorized(g)
			return
		}
		g.Redirect(http.StatusFound, authRedirectURL(g))
		g.Abort()
	}
}
