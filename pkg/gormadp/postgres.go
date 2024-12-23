package gormadp

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// PgConfig defines config for postgres driver
type PgConfig struct {
	// SSLMode mode of ssl that using to connect to Database
	//
	// Optional. Default is true
	SSLMode bool

	// TimeZone that want to use as a default timezone in Database
	//
	// Optional. Default is "Asia/Jakarta"
	TimeZone string

	cfg Config
}

// Driver name identity
const Postgres Dialect = "postgres"

// PgOptions create new options of postgres driver
func (d Dialect) PgOptions(c PgConfig) *PgConfig {
	return &PgConfig{
		SSLMode:  c.SSLMode,
		TimeZone: c.TimeZone,
	}
}

// Return postgres dialector
func (c *PgConfig) dialector() gorm.Dialector {
	return postgres.Open(resolveDSN(c))
}

// Return postgres dsn
func resolveDSN(c *PgConfig) string {
	// Variable
	var (
		sslMode  string
		timezone = "Asia/Jakarta"
	)

	// Determine option
	if c.SSLMode {
		sslMode = "enable"
	} else {
		sslMode = "disable"
	}
	if len(c.TimeZone) > 0 && c.TimeZone != " " {
		timezone = c.TimeZone
	}

	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		c.cfg.Host,
		c.cfg.User,
		c.cfg.Password,
		c.cfg.DBName,
		c.cfg.Port,
		sslMode,
		timezone,
	)
}
