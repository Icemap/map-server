package db

import (
	"fmt"
	"github.com/cockroachdb/pebble"
	"github.com/spf13/viper"
	"testing"
)

func TestAllBaseFunc(t *testing.T) {
	viper.AddConfigPath("../config")

	Del([]byte("not_exist_key"))

	_, err := Get([]byte("not_exist_key"))
	if err != pebble.ErrNotFound {
		fmt.Println("should be pebble.ErrNotFound")
		t.Error("should be pebble.ErrNotFound")
		return
	}

	err = Set([]byte("key"), []byte("value"))
	if err != nil {
		fmt.Println(err)
		t.Error(err)
		return
	}

	value, err := Get([]byte("key"))
	if err != nil {
		fmt.Println(err)
		t.Error(err)
		return
	}

	if string(value) != "value" {
		t.Error("value error")
		return
	}

	err = Del([]byte("key"))
	if err != nil {
		fmt.Println(err)
		t.Error(err)
		return
	}

	_, err = Get([]byte("key"))
	if err != pebble.ErrNotFound {
		fmt.Println("should be pebble.ErrNotFound")
		t.Error("should be pebble.ErrNotFound")
		return
	}
}
