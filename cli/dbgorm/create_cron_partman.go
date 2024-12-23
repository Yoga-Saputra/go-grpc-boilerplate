package dbgorm

import (
	"errors"

	"github.com/Yoga-Saputra/go-grpc-boilerplate/config"
	"github.com/Yoga-Saputra/go-grpc-boilerplate/pkg/gormadp"
)

// Will create cron schedule to call partition manager process.
func CreateCronPartman(a *AdvanceDBMigrate) error {
	// Prepare and open new DB connection
	var err error
	var db *gormadp.DBAdapter
	for _, t := range config.Of.Database.Tools {
		if t.Identifier == "pgcron" {
			cfg := gormadp.Config{
				Host:     t.Host,
				Port:     t.Port,
				User:     t.User,
				Password: t.Password,
				DBName:   t.Name,
				Dialect:  gormadp.Dialect(t.Dialect),
			}
			opts := cfg.Dialect.PgOptions(gormadp.PgConfig{
				SSLMode:  false,
				TimeZone: "Asia/Manila",
			})
			db = gormadp.Open(cfg, opts)
			break
		}
	}
	if db == nil {
		return errors.New("cannot open or create new DB connection")
	}

	// Prepare and execute SQL statements
	tx := db.DB
	jobname := "maintenance_partition"
	cronSchedule := "0 8 1 * *"
	database := config.Of.Database.Name

	err = tx.Exec(
		`SELECT cron.schedule(?, ?, $$CALL partman.run_maintenance_proc()$$)`,
		jobname,
		cronSchedule,
	).Error
	if err != nil {
		return err
	}

	err = tx.Exec(
		`UPDATE cron.job SET database = ? WHERE jobname = ?`,
		database,
		jobname,
	).Error
	if err != nil {
		return err
	}

	// Close DB connection and finish
	db.Close()
	return err
}
