package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func livenessProbe(c *gin.Context) {
	c.String(http.StatusOK, "alive")
}
