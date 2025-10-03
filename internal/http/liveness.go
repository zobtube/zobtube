package http

import (
	"github.com/gin-gonic/gin"
)

func (s *Server) livenessProbe(c *gin.Context) {
	if s.healthy {
		c.String(200, "alive")
	} else {
		c.String(500, "ko")
	}
}
