package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/zobtube/zobtube/internal/model"
)

func (c *Controller) HTML(g *gin.Context, httpCode int, webPage string, parameters gin.H) {
	parameters["User"] = g.MustGet("user").(*model.User)
	parameters["AuthenticationEnabled"] = c.config.Authentication
	g.HTML(httpCode, webPage, parameters)
}
