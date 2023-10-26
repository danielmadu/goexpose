package handlers

import (
	"net/http"

	"github.com/danielmadu/goexpose/config"
	"github.com/gin-gonic/gin"
)

func UpdateConfig(c *gin.Context) {
	c.Header("Content-type", "application/json")

	config := config.GetConfig()

	c.BindJSON(&config)

	c.JSON(http.StatusOK, config)
}
