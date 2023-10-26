package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/danielmadu/goexpose/config"
	"github.com/gin-gonic/gin"
)

func UpdateClient(c *gin.Context) {
	ping := config.Ping{}
	config := config.GetConfig()
	remoteAddr := c.Request.RemoteAddr

	fmt.Println(remoteAddr)

	c.BindJSON(&ping)

	res := strings.Split(remoteAddr, ":")

	config.Shared = "http://" + res[0] + ":" + ping.Port

	fmt.Println(config.Shared)

	c.JSON(http.StatusOK, config)
}
