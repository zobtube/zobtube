package controller

import (
	"github.com/gin-gonic/gin"
)

func (c *Controller) ErrUnauthorized(g *gin.Context) {
	c.HTML(g, 401, "err/unauthorized.html", gin.H{})
}
