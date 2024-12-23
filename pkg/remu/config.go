// Package redis is singleton pattern for redis adapter
// This package using github.com/redis/go-redis/v9
package remu

// Config defines the config for redis client.
type Config struct {
	// Host name where the DB is hosted
	//
	// Optional. Default is "127.0.0.1"
	Host string

	// Port where the DB is listening on
	//
	// Optional. Default is 3306
	Port int

	// Server username
	//
	// Optional. Default is ""
	Username string

	// Server password
	//
	// Optional. Default is ""
	Password string

	// Database to be selected after connecting to the server.
	//
	// Optional. Default is 0
	Database int

	// URL the standard format redis url to parse all other options. If this is set all other config options, Host, Port, Username, Password, Database have no effect.
	//
	// Example: redis://<user>:<pass>@localhost:6379/<db>
	// Optional. Default is ""
	URL string

	// Maximum number of retries before giving up.
	//
	// Optional. Default is 3 retries; -1 (not 0) disables retries.
	MaxRetries int

	// Maximum number of socket connections.
	//
	// OPtional. Default is 10 connections per every available CPU as reported by runtime.GOMAXPROCS.
	PoolSize int

	// Redis key namespace
	//
	// Optional. default is "goredis-adapter"
	NameSpace string
}

// ConfigDefault is the default config
var ConfigDefault = Config{
	Host:       "127.0.0.1",
	Port:       6379,
	Username:   "",
	Password:   "",
	Database:   0,
	URL:        "",
	MaxRetries: 3,
	PoolSize:   10,
	NameSpace:  "goredis-adapter",
}

// Helper function to set default values
func configDefault(config ...Config) Config {
	// Return default config if nothing provided
	if len(config) < 1 {
		return ConfigDefault
	}

	// Override default config
	cfg := config[0]

	// Set default values
	if cfg.Host == "" {
		cfg.Host = ConfigDefault.Host
	}
	if cfg.Port <= 0 {
		cfg.Port = ConfigDefault.Port
	}
	if cfg.MaxRetries == 0 {
		cfg.MaxRetries = ConfigDefault.MaxRetries
	}
	if cfg.PoolSize <= 10 {
		cfg.PoolSize = ConfigDefault.PoolSize
	}
	if cfg.NameSpace == "" || cfg.NameSpace == " " {
		cfg.NameSpace = ConfigDefault.NameSpace
	}

	return cfg
}
