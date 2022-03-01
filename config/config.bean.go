package config

// Service
type Service struct {
	DbPath        string `yaml:"dbPath"`
	DownloadRetry int    `yaml:"downloadRetry"`
	Port          int    `yaml:"port"`
	Concurrency   int    `yaml:"concurrency"`
	MapPath       string `yaml:"mapPath"`
}

// Config
type Config struct {
	Log     Log     `yaml:"log"`
	Service Service `yaml:"service"`
}

// Log
type Log struct {
	Level     string `yaml:"level"`
	LogName   string `yaml:"logName"`
	ErrorName string `yaml:"errorName"`
}

