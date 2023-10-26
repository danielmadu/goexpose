package handlers

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func All(c *gin.Context) {

	path := c.Param("all")
	if strings.HasPrefix(path, "/goexpose/api/config") {
		UpdateConfig(c)
		return
	}

	if strings.HasPrefix(path, "/goexpose/api/ping") {
		UpdateClient(c)
		return
	}

	Tunnel(c)
}
