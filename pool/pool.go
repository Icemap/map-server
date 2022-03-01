package pool

import (
	"github.com/panjf2000/ants/v2"
	"map-server/config"
	"sync"
)

var (
	workingPool *ants.Pool
	once        sync.Once
)

func Submit(submitFunc func()) error {
	once.Do(func() {
		var err error
		workingPool, err = ants.NewPool(config.ReadConfig().Service.Concurrency)
		if err != nil {
			panic(err)
		}
	})

	return workingPool.Submit(submitFunc)
}
