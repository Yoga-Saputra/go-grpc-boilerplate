package gormadp

import (
	"gorm.io/gorm"
)

// Config defines the config for App Database
type Config struct {
	// Host name where the Database is hosted
	//
	// Required. Default is "127.0.0.1"
	Host string

	// Port number of Database connection
	//
	// Required. Default 5432
	Port int

	// User name of Database connection
	//
	// Required. Default is "test"
	User string

	// Password of Database conenction
	//
	// Required. Default is "test"
	Password string

	// DBName name of Database that want to connect
	//
	// Required. Default is "test"
	DBName string

	// Dialect is varian or type of database query language
	//
	// Required. Default is "postgres"
	Dialect Dialect

	// Options is config of gorm
	//
	// Optional. Default is nil
	Options gorm.Option
}

type DialectOptions func(d Dialect, c interface{}) interface{}

type Dialect string

type connectionConfig struct {
	dialector gorm.Dialector
}

// ConfigDefault is the default config
var ConfigDefault = Config{
	Host:     "127.0.0.1",
	Port:     5432,
	User:     "testing",
	Password: "testing123",
	DBName:   "testing",
	Dialect:  "postgres",
	Options:  &gorm.Config{},
}

// Helper function to set default config
func configDefault(config Config, driverOptions ...interface{}) *connectionConfig {
	// Overide default config
	cfg := config

	// Validating
	if cfg.Host == "" {
		cfg.Host = ConfigDefault.Host
	}
	if cfg.Port == 0 {
		cfg.Port = ConfigDefault.Port
	}
	if cfg.User == "" {
		cfg.User = ConfigDefault.User
	}
	if cfg.Password == "" {
		cfg.Password = ConfigDefault.Password
	}
	if cfg.DBName == "" {
		cfg.DBName = ConfigDefault.DBName
	}
	if cfg.Dialect == "" {
		cfg.Dialect = ConfigDefault.Dialect
	}
	if cfg.Options == nil {
		cfg.Options = config.Options
	}

	return finalConfigBasedOnDriver(cfg, driverOptions...)
}

// Helper function to set final "connection-use" config based on choosen driver
func finalConfigBasedOnDriver(config Config, driverOptions ...interface{}) *connectionConfig {
	switch config.Dialect {
	// PostgreSQL Driver
	case Postgres:
		pgc := &PgConfig{cfg: config}
		if len(driverOptions) > 0 {
			// dop, ok := driverOptions[0].(*PgConfig)
			dop, ok := driverOptions[0].(*PgConfig)
			if !ok {
				panic("Given option is not belongs to postgres")
			}
			pgc.SSLMode = dop.SSLMode
			pgc.TimeZone = dop.TimeZone
		}
		return &connectionConfig{dialector: pgc.dialector()}

	// Unknown Driver
	default:
		panic("Driver name is not available. Current supported driver: postgres")
	}
}
