package router

import (
	"fmt"

	"github.com/danielmadu/goexpose/config"
	"github.com/danielmadu/goexpose/handlers"
	"github.com/danielmadu/goexpose/utilities"
	"github.com/gin-gonic/gin"
)

func Init() error {

	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	router.Any("/*all", handlers.All)

	conf := config.GetConfig()

	if utilities.IsValid(conf.CertFile) && utilities.IsValid(conf.KeyFile) {
		fmt.Println("SSL Enabled")
		return router.RunTLS(":3000", conf.CertFile, conf.KeyFile)
	}

	return router.Run(":3000")

}
