package db

import (
	"fmt"
	"github.com/cockroachdb/pebble"
	"map-server/config"
	"map-server/logger"
)

// Set set a k-v pair to db
func Set(key, value []byte) error {
	return SetAtomic(key, value)
}

// SetAtomic set k-v pair list to db, this function is atomically
func SetAtomic(kvList ...[]byte) error {
	db, err := pebble.Open(config.ReadConfig().Service.DbPath, &pebble.Options{DisableWAL: false})
	if err != nil {
		logger.Logger().Errorw("[db] can not open database", "err", err)
		panic(err)
	}
	defer db.Close()

	if len(kvList)%2 != 0 {
		err = fmt.Errorf("[db] kv list length not a even number")
		logger.Logger().Errorw("[db] kv list length not a even number", "err", err)
		return err
	}

	batch := db.NewBatch()
	for i := 0; i < len(kvList); i += 2 {
		if err = batch.Set(kvList[i], kvList[i+1], pebble.Sync); err != nil {
			logger.Logger().Errorw("[db] set error", "err", err)
			batchErr := batch.Close()
			if batchErr != nil {
				logger.Logger().Errorw("[db] close batch error", "err", err)
			}

			return err
		}
	}

	if err = batch.Commit(pebble.Sync); err != nil {
		logger.Logger().Errorw("[db] commit batch error", "err", err)
		return err
	}

	if err = db.Flush(); err != nil {
		logger.Logger().Errorw("[db] flush error", "err", err)
		return err
	}

	return nil
}

// Get get value from db by key
func Get(key []byte) ([]byte, error) {
	db, err := pebble.Open(config.ReadConfig().Service.DbPath, &pebble.Options{DisableWAL: false})
	if err != nil {
		logger.Logger().Errorw("[db] can not open database", "err", err)
		panic(err)
	}
	defer db.Close()

	value, closer, err := db.Get(key)
	if err != nil {
		logger.Logger().Errorw("[db] get value error", "err", err)
		return nil, err
	}

	if err = closer.Close(); err != nil {
		logger.Logger().Errorw("[db] closer close error", "err", err)
		return nil, err
	}

	return value, nil
}

// Del delete a k-v pair from db
func Del(key []byte) error {
	return DelAtomic(key)
}

// DelAtomic delete a k-v pair list from db, this function is atomically
func DelAtomic(keyList ...[]byte) error {
	db, err := pebble.Open(config.ReadConfig().Service.DbPath, &pebble.Options{DisableWAL: false})
	if err != nil {
		logger.Logger().Errorw("[db] can not open database", "err", err)
		panic(err)
	}
	defer db.Close()

	batch := db.NewBatch()
	for i := 0; i < len(keyList); i += 2 {
		if err = batch.Delete(keyList[i], pebble.Sync); err != nil {
			logger.Logger().Errorw("[db] del error", "err", err)
			batchErr := batch.Close()
			if batchErr != nil {
				logger.Logger().Errorw("[db] close batch error", "err", err)
			}

			return err
		}
	}

	if err = batch.Commit(pebble.Sync); err != nil {
		logger.Logger().Errorw("[db] commit batch error", "err", err)
		return err
	}

	if err = db.Flush(); err != nil {
		logger.Logger().Errorw("[db] flush error", "err", err)
		return err
	}

	return nil
}
