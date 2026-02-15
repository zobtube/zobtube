package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const ErrConfigAbsent = "configuration not found, restarting the appliaction should fix the issue"

func (c *Controller) ErrNotFound(g *gin.Context) {
	g.JSON(http.StatusNotFound, gin.H{"error": "not found"})
}

func (c *Controller) ErrUnauthorized(g *gin.Context) {
	g.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
}
