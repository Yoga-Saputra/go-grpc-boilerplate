package xtalsvcs

import (
	"errors"
	"strings"

	"github.com/Yoga-Saputra/go-grpc-boilerplate/config"
)

// Helper function to check existences of service config
func checkConfigExists(serviceKey string) error {
	serviceKey = strings.ToLower(serviceKey)
	host := getConfigValue(serviceKey, "host")

	if len(host) <= 0 || host == " " {
		return errors.New("configuration not found")
	}

	return nil
}

// Helper function to get config value based on given service key name
func getConfigValue(serviceKey string, requisite string) string {
	serviceKey = strings.ToLower(serviceKey)
	requisite = strings.ToLower(requisite)

	if requisite == "host" {
		return config.Of.External.Host[serviceKey]
	}

	return config.Of.External.Token[serviceKey]
}
