package gormadp

import (
	"database/sql"
	"log"

	"gorm.io/gorm"
)

// DBAdapter defines pointer value of opened database
type DBAdapter struct {
	DB  *gorm.DB
	SQL *sql.DB
}

// Open creates and open a new database connection
func Open(config Config, driverOptions ...interface{}) *DBAdapter {
	// Set default config
	cfg := configDefault(config, driverOptions...)

	// Open new database connection
	db, err := gorm.Open(cfg.dialector, config.Options)
	if err != nil {
		panic(err)
	}

	// Expose go native sql
	sql, err := db.DB()
	if err != nil {
		panic(err)
	}

	return &DBAdapter{
		DB:  db,
		SQL: sql,
	}
}

// NewConnectionConfig just creates new final config of database connection
// Useful when using gorm plugin such as dbresolver
func NewConnectionConfig(config Config, driverOptions ...interface{}) *connectionConfig {
	return configDefault(config, driverOptions...)
}

// GetNilSession return gorm DB session
func (d *DBAdapter) GetNilSession() *gorm.Session {
	return &gorm.Session{}
}

// Close to closing database connection
// It is a good approach for graceful shutdown?
func (d *DBAdapter) Close() {
	if d != nil {
		if d.SQL != nil {
			if err := d.SQL.Close(); err != nil {
				log.Printf("[DB Adapter] - Failed to close db: %v \n", err.Error())
			} else {
				log.Println("[DB Adapter] - DB already closed")
			}
		}
	}
}
