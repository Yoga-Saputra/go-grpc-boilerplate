package dbgorm

import (
	"fmt"

	"github.com/Yoga-Saputra/go-grpc-boilerplate/app/entity"
)

// CreateIndexByreferenceOnly for transaction log table for VRT provider.
func CreateIndexByreferenceOnly(a *AdvanceDBMigrate) error {
	type TransactionLogProvider struct {
		Reference string `gorm:"index"`
	}

	return a.Tx.
		Scopes(entity.ProvSchemaTable(a.Entity.(*entity.TransactionLogProvider), fmt.Sprintf("%v", a.Schema))).
		Migrator().CreateIndex(&TransactionLogProvider{}, "Reference")
}
