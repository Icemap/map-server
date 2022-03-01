package service

import (
	"github.com/gin-gonic/gin"
	"map-server/config"
	"net/http"
)

// configHandle handle config request
func configHandle(c *gin.Context) {
	c.JSON(http.StatusOK, config.ReadConfig())
}
