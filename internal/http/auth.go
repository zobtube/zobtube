package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"gitlab.com/zobtube/zobtube/internal/controller"
	"gitlab.com/zobtube/zobtube/internal/model"
)

const cookieName = "zt_auth"

func UserIsAuthenticated(c controller.AbtractController) gin.HandlerFunc {
	return func(g *gin.Context) {
		cookie, err := g.Cookie(cookieName)
		if err != nil {
			// cookie not set
			g.Redirect(http.StatusFound, "/auth")
			return
		}

		// get session
		session := &model.UserSession{
			ID: cookie,
		}
		result := c.GetSession(session)

		// check result
		if result.RowsAffected < 1 {
			g.Redirect(http.StatusFound, "/auth")
			return
		}

		// check validity
		if session.ValidUntil.Before(time.Now()) {
			g.Redirect(http.StatusFound, "/auth")
			return
		}

		// check if user is authenticated
		if session.UserID == "" {
			g.Redirect(http.StatusFound, "/auth")
			return
		}

		// get user
		user := &model.User{
			ID: session.UserID,
		}
		result = c.GetUser(user)
		if result.RowsAffected < 1 {
			g.Redirect(http.StatusFound, "/auth")
			return
		}

		// set meta in context
		g.Set("user", user)

		g.Next()
	}
}
