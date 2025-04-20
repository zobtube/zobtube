package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/zobtube/zobtube/internal/model"
)

func (c *Controller) ErrUnauthorized(g *gin.Context) {
	g.HTML(401, "err/unauthorized.html", gin.H{
		"User":       g.MustGet("user").(*model.User),
	})
}
