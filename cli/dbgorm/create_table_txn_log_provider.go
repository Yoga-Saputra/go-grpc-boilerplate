package dbgorm

import (
	"errors"
	"fmt"

	"github.com/Yoga-Saputra/go-grpc-boilerplate/app/entity"
	"gorm.io/plugin/dbresolver"
)

func CreateTableTxnLogProvider(a *AdvanceDBMigrate) error {
	if a.Source != "" {
		a.Tx = a.Tx.Clauses(dbresolver.Use(a.Source))
	}
	tx := a.Tx.Debug()

	// Get DB dialect
	dialect := tx.Dialector.Name()

	// Do create DB stuff based on dialect
	switch dialect {
	case "postgres":
		// Create schema if not exists
		tx.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %v", a.Schema))

		// Declare migration options
		declared := migrateWithTableOpt{
			Entity: a.Entity,
			opts: map[string]interface{}{
				"gorm:table_options": " PARTITION BY RANGE (date)",
			},
		}

		// Append migrate options
		for optK, optV := range declared.opts {
			tx = tx.Set(optK, optV)
		}
		ent := declared.Entity.(*entity.TransactionLogProvider)
		tx = tx.Scopes(entity.ProvSchemaTable(ent, fmt.Sprintf("%v", a.Schema)))

		// Create parent table partition
		if err := tx.AutoMigrate(ent); err != nil {
			return err
		}

		// Setup for partition table
		if err := SetupPartitionedTable(a); err != nil {
			return err
		}

		return nil

	default:
		return errors.New("db Dialect not supported yet")
	}
}
