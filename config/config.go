package config

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// The configuration that underlying all config
type config struct {
	// Base App configuration
	App app `json:"app" yaml:"app"`

	// Database configuration
	Database database `json:"database" yaml:"database"`

	// Cache configuration
	Cache cache `json:"cache"`

	// Queue configuration
	Queue queue `json:"queue" yaml:"queue"`

	// Kafka configuration
	Kafka kafka `json:"kafka" yaml:"kafka"`

	// External API/Microservices configuration
	External external `json:"external" yaml:"external"`
}

// Of is the config context that will be called by another package
var Of config

// Run viper setup and marshaling the config once at the runtime
func init() {
	viper.SetConfigName(".config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../")
	viper.AddConfigPath("/etc/seamless-wallet/")
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("Config file has been changed, re-load that")
		load()
	})

	load()
}

// Load and marshaling the config file
func load() {
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(&Of); err != nil {
		panic(err)
	}
}
