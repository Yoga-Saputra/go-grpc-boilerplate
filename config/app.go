package config

import (
	"io/ioutil"
	"log"
	"path"

	"github.com/Yoga-Saputra/go-grpc-boilerplate/pkg/rsa256"
)

// Base app configuration key value
type app struct {
	// App name
	Name string `json:"name" yaml:"name"`

	// App description
	Desc string `json:"desc" yaml:"desc"`

	// Port number that app will running
	Port int `json:"port" yaml:"port"`

	// Environtment of the app
	//
	// "development" or "production"
	Env string `json:"env" yaml:"env"`

	// Prefork mode on REST API server
	//
	// If true, app will running on prefork mode (multi child proccess)
	Prefork bool `json:"prefork" yaml:"prefork"`

	// Is abbreviation of Auto Update Watcher Interval
	//
	// Value must be int, and interval should be in second(s)
	AUWI int `json:"auwi" yaml:"auwi"`

	// Time Zone of the app
	//
	// i.e "Asia/Manila"
	TimeZone string `json:"timeZone" yaml:"timeZone"`

	// Full path where program executable places
	ProgramFile string `json:"programFile" yaml:"programFile"`

	// Full path of working directory location
	WorkingDir string `json:"workingDir" yaml:"workingDir"`

	// Runtime modified
	publicKey interface{}
	secretKey []byte
}

// Debug global debug flag based on app env.
// Return debug true/false
func (a *app) Debug() bool {
	return a.Env != "production"
}

// ResolveFilePathInWorkDir return full path of given file name in working directory
func (a *app) ResolveFilePathInWorkDir(f string) string {
	return path.Join(a.WorkingDir, "/", f)
}

// GetPublicKey return public key
func (a *app) GetPublicKey() interface{} {
	if a.publicKey == nil {
		f := a.ResolveFilePathInWorkDir("public-key.pem")
		pubk, err := rsa256.ReadPublicKey(f)
		if err != nil {
			log.Printf("[App][GetPublicKey] - Error: %s", err.Error())
			return nil
		}

		a.publicKey = pubk
	}

	return a.publicKey
}

// GetSecretKey return secret key
func (a *app) GetSecretKey() string {
	if a.secretKey == nil {
		// Read secret.key file
		f := a.ResolveFilePathInWorkDir("secret.key")
		r, err := ioutil.ReadFile(f)
		if err != nil {
			log.Printf("[App][GetSecretKey] - Error: %s", err.Error())
			return ""
		}

		a.secretKey = r
	}

	return string(a.secretKey)
}
