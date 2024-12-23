package httpclient

import (
	"net/http"
	"net/http/cookiejar"
	"time"
)

// Config defines http client initial config
type Config struct {
	// Base URL
	//
	// Optional. Default ""
	BaseUrl string

	// HTTP Headers
	//
	// Optional. Default "User-Agent" : "GoNoob/1.0"
	Headers map[string]interface{}

	// Cookie jar
	//
	// Optional. Default http.CookieJar(nil)
	CookieJar http.CookieJar

	// TLS verify enable
	//
	// Optional. Default true/disabled
	SkipInsecure bool

	// Given auth token. Not exportable.
	token string

	// Given request timeout. Not exportable.
	timeout time.Duration
}

// ConfigDefault is the default config
var ConfigDefault = Config{
	BaseUrl: "",
	Headers: map[string]interface{}{
		"User-Agent": "GoNoob/1.0",
	},
	SkipInsecure: true,
}

// Helper function to set default config
func configDefault(config ...Config) *Config {
	jar, _ := cookiejar.New(nil)

	// Return default config
	if len(config) <= 0 {
		cfg := &ConfigDefault
		cfg.CookieJar = jar
		return cfg
	}

	// Overide default config
	cfg := config[0]

	// Headers
	if len(cfg.Headers) <= 0 {
		cfg.Headers = ConfigDefault.Headers
	} else {
		headers := make(map[string]interface{})
		if cfg.Headers["User-Agent"] == nil {
			cfg.Headers["User-Agent"] = ConfigDefault.Headers["User-Agent"]
		}

		for k, v := range cfg.Headers {
			headers[k] = v
		}

		cfg.Headers = headers
	}

	// Cookie jar
	if cfg.CookieJar == nil {
		cfg.CookieJar = jar
	}

	return &cfg
}
