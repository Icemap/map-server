package config

import (
	"fmt"
	"github.com/spf13/viper"
	"sync"
)

var (
	once       sync.Once
	configData Config
)

func ReadConfig() Config {
	once.Do(func() {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("config")
		if err := viper.ReadInConfig(); err != nil {
			panic(fmt.Errorf("Fatal error config file: %w \n", err))
		}

		err := viper.Unmarshal(&configData)
		if err != nil {
			panic(fmt.Errorf("Fatal error config file: %w \n", err))
		}
	})
	return configData
}
