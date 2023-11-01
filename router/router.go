package router

import (
	"github.com/danielmadu/goexpose/handlers"
	"github.com/gin-gonic/gin"
)

func Init() error {

	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	router.Any("/*all", handlers.All)

	return router.Run(":3000")
}
