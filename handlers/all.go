package handlers

import (
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/websocket"
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

	if strings.HasPrefix(path, "/goexpose/ws") {
		websocket.Handler(WebSocket).ServeHTTP(c.Writer, c.Request)
		return
	}

	Tunnel(c)
}
