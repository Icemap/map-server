package service

import (
	"github.com/gin-gonic/gin"
	"map-server/config"
)

// configHandle handle config request
func configHandle(c *gin.Context) {
	c.JSON(200, config.ReadConfig())
}
