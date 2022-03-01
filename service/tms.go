package service

import (
	"github.com/gin-gonic/gin"
	"map-server/logger"
	"net/http"
)

type TMSRequest struct {
	X int `form:"x" json:"x" xml:"x"  binding:"required"`
	Y int `form:"y" json:"y" xml:"y"  binding:"required"`
	Z int `form:"z" json:"z" xml:"z"  binding:"required"`
}

// tmsHandler TMS service
func tmsHandler(c *gin.Context) {
	req := TMSRequest{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.Header("Content-Type", "image/jpeg")

	path := tilePathBuilder(req.X, req.Y, req.Z)
	logger.Logger().Debugw("tms", "path", path)
	c.File(path)
}
