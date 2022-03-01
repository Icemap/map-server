package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"map-server/config"
)

func Serve() {
	r := gin.New()
	r.Use(GinLogger(), GinRecovery(true))

	r.GET("/config", configHandle)
	r.POST("/map", mapDownloadHandler)

	err := r.Run(fmt.Sprintf(":%d", config.ReadConfig().Service.Port))
	if err != nil {
		panic(err)
	}
}
