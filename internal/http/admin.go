package http

import (
	"github.com/gin-gonic/gin"

	"github.com/zobtube/zobtube/internal/controller"
	"github.com/zobtube/zobtube/internal/model"
)

func UserIsAdmin(c controller.AbstractController) gin.HandlerFunc {
	return func(g *gin.Context) {
		// get user
		user := g.MustGet("user").(*model.User)

		// check if admin
		if !user.Admin {
			g.Redirect(307, "/error/unauthorized")
			g.Abort()
			return
		}

		// all good
		g.Next()
	}
}
