package service

import (
	"fmt"
	"github.com/Icemap/coordinate"
	"map-server/logger"
	"map-server/pool"
)

type (
	LevelTileRange struct {
		level        int
		levelTileMax int

		minX int
		minY int
		maxX int
		maxY int
	}

	LevelTileRangeIterFunc func(level, x, y int)
)

// InitLevelTileRange create a LevelTileRange instance
func InitLevelTileRange(level int, leftTop, rightBottom coordinate.Coordinate) (LevelTileRange, error) {
	minX, minY, err := coordinate.WGS84ToWebMercatorTile(leftTop, level)

	if err != nil {
		logger.Logger().Errorw("parse leftTop coordinate error", "err", err)
		return LevelTileRange{}, err
	}

	maxX, maxY, err := coordinate.WGS84ToWebMercatorTile(rightBottom, level)
	if err != nil {
		logger.Logger().Errorw("parse rightBottom coordinate error", "err", err)
		return LevelTileRange{}, err
	}

	if minY >= maxY {
		return LevelTileRange{}, fmt.Errorf("in WGS84 top [%f] should bigger than bottom [%f]", leftTop.Y, rightBottom.Y)
	}

	return LevelTileRange{
		level:        level,
		levelTileMax: (1 << level) - 1,
		minX:         minX,
		minY:         minY,
		maxX:         maxX,
		maxY:         maxY,
	}, nil
}

func (levelTileRange *LevelTileRange) EffectAll(runFunc LevelTileRangeIterFunc) {
	for y := levelTileRange.minY; y <= levelTileRange.maxY; y++ {
		// consider situation of crossing prime meridian
		for x := levelTileRange.minX; x != levelTileRange.maxX; x++ {
			if x > levelTileRange.levelTileMax {
				x = 0
			}

			xCopy, yCopy := x, y
			pool.Submit(func() {
				runFunc(levelTileRange.level, xCopy, yCopy)
			})
		}
	}
}
