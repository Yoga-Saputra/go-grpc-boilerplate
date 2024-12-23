package dbgorm

import (
	"github.com/Yoga-Saputra/go-grpc-boilerplate/app/entity"
)

// AlterAllTableProviderLog for altering all table provider log table.
func AlterAllTableProviderLog(a *AdvanceDBMigrate) error {
	tx := a.Tx
	var schemas []string
	notInSchemas := []string{
		"pg_toast",
		"pg_catalog",
		"public",
		"information_schema",
		"partman",
		"cron",
	}

	if err := tx.
		Raw("SELECT schema_name FROM information_schema.schemata WHERE schema_name NOT IN ?", notInSchemas).
		Scan(&schemas).
		Error; err != nil {
		return err
	}

	for _, s := range schemas {
		if err := tx.
			Scopes(entity.ProvSchemaTable(a.Entity.(*entity.TransactionLogProvider), s)).
			AutoMigrate(&entity.TransactionLogProvider{}); err != nil {
			return err
		}

		if err := tx.
			Scopes(entity.ProvSchemaTablePartmanTemplate(a.Entity.(*entity.TransactionLogProvider), s)).
			AutoMigrate(&entity.TransactionLogProvider{}); err != nil {
			return err
		}
	}
	return nil
}
