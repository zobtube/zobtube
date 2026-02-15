package http

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/zobtube/zobtube/internal/controller"
	"github.com/zobtube/zobtube/internal/model"
)

func UserIsAdmin(c controller.AbstractController) gin.HandlerFunc {
	return func(g *gin.Context) {
		user := g.MustGet("user").(*model.User)
		if !user.Admin {
			if strings.HasPrefix(g.Request.URL.Path, "/api") {
				g.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
				g.Abort()
				return
			}
			g.Redirect(http.StatusTemporaryRedirect, "/api/error/unauthorized")
			g.Abort()
			return
		}
		g.Next()
	}
}
