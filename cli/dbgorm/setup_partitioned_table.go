package dbgorm

import (
	"fmt"
	"time"

	"github.com/Yoga-Saputra/go-grpc-boilerplate/app/entity"
)

func SetupPartitionedTable(a *AdvanceDBMigrate) error {
	loc, err := time.LoadLocation("Asia/Manila")
	if err != nil {
		return err
	}

	ent := a.Entity.(*entity.TransactionLogProvider)
	now := time.Now().In(loc)
	timestampStr := now.AddDate(0, 0, -1).Format("2006-01-02")
	tableName := fmt.Sprintf("%v.%v", a.Schema, ent.TableName())

	// Call setup partition procedure
	if err := a.Tx.
		Exec(`CALL setup_partition_proc(?, ?)`, tableName, timestampStr).
		Error; err != nil {
		return err
	}

	return nil
}
