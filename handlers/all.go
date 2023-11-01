package handlers

import (
	"strings"

	"github.com/danielmadu/goexpose/config"
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
		config := config.GetConfig()
		requestToken := c.Query("token")
		if requestToken != config.Token {
			c.Status(403)
			return
		}
		websocket.Handler(WebSocket).ServeHTTP(c.Writer, c.Request)
		return
	}

	Tunnel(c)
}
