package service

import (
	"fmt"
	"github.com/Icemap/coordinate"
	"github.com/gin-gonic/gin"
	"io"
	"map-server/config"
	"map-server/logger"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
)

type MapDownloadRequest struct {
	MapType     string                `form:"mapType" json:"mapType" xml:"mapType"  binding:"required"`
	Level       int                   `form:"level" json:"level" xml:"level"  binding:"required"`
	LeftTop     coordinate.Coordinate `form:"leftTop" json:"leftTop" xml:"leftTop"  binding:"required"`
	RightBottom coordinate.Coordinate `form:"rightBottom" json:"rightBottom" xml:"rightBottom"  binding:"required"`
}

var validMapTypeSet = map[string]interface{}{
	coordinate.GoogleSatellite: nil,
	coordinate.GoogleImage:     nil,
	coordinate.GoogleTerrain:   nil,
	coordinate.AMapSatellite:   nil,
	coordinate.AMapCover:       nil,
	coordinate.AMapImage:       nil,
}

// mapDownloadHandler download map by request
func mapDownloadHandler(c *gin.Context) {
	req := MapDownloadRequest{}
	err := c.BindJSON(&req)
	if err != nil {
		logger.Logger().Errorw("[Post] /map: request error", "err", err)
		return
	}

	if _, exist := validMapTypeSet[req.MapType]; !exist {
		err = fmt.Errorf("map type [%s] error", req.MapType)
		c.JSON(http.StatusBadRequest, err)
		logger.Logger().Errorw("[Post] /map: map type error", "err", err)
		return
	}

	tileRange, err := InitLevelTileRange(req.Level, req.LeftTop, req.RightBottom)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		logger.Logger().Errorw("[Post] /map: level init error", "err", err)
		return
	}

	go tileRange.EffectAll(func(level, x, y int) {
		// download
		url := coordinate.WebMercatorTileToURL(req.MapType, x, y, level)
		tilePath := strings.Join([]string{config.ReadConfig().Service.MapPath,
			strconv.Itoa(level), strconv.Itoa(x), strconv.Itoa(y), "pic.jpg"}, string(os.PathSeparator))

		err = download(url, tilePath)
		for i := 0; i < config.ReadConfig().Service.DownloadRetry; i++ {
			if err != nil {
				logger.Logger().Errorw("download pic error",
					"retry num", i, "level", level, "x", x, "y", y)
				err = download(url, tilePath)
			}
		}

		logger.Logger().Debugw("[download] download success", "level", level, "x", x, "y", y)
	})
}

// download file downloader
func download(url, filePath string) error {
	if fileExist(filePath) {
		return nil
	}

	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("get file error: %+v", err)
	}
	defer res.Body.Close()

	err = os.MkdirAll(path.Dir(filePath), os.FileMode(0777))
	if err != nil {
		return fmt.Errorf("create dir error: %+v", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("create file error: %+v", err)
	}
	defer file.Close()

	_, err = io.Copy(file, res.Body)
	if err != nil {
		return fmt.Errorf("save file error: %+v", err)
	}

	return nil
}

func fileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}
