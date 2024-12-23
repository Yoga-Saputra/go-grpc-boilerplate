package repo

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Yoga-Saputra/go-grpc-boilerplate/app/entity"
	"gorm.io/gorm"
)

type (
	txnLogRepoDB struct {
		db *gorm.DB
	}
)

// NewTxnLogRepoDB create new DB repo for transaction log entity.
func NewTxnLogRepoDB(db *gorm.DB) *txnLogRepoDB {
	if db != nil {
		return &txnLogRepoDB{db}
	}

	return nil
}

// Create new wallet record.
func (trd *txnLogRepoDB) Create(t *entity.TransactionLogInternal, itx interface{}) error {
	// t should != nil
	if t == nil {
		return errors.New("entity transaction log argument pointer cannot be nil")
	}

	// Manually freed memory
	defer func() {
		t = nil
	}()

	// Prepare transaction session
	tx := trd.db
	if itx != nil {
		itxdb, ok := itx.(*gorm.DB)
		if !ok {
			return errors.New("cannot assert argument itx to the *gorm.DB type")
		}

		tx = itxdb
	}

	// Insert
	if err := tx.Create(t).Error; err != nil {
		switch {
		case strings.Contains(err.Error(), "SQLSTATE 23505"):
			err = fmt.Errorf("duplicate transaction with serviceId: %d, reference: %s", t.ServiceID, t.Reference)
		}

		return err
	}

	return nil
}

// delete log internal transaction
func (trd *txnLogRepoDB) Delete(conds map[string]interface{}, itx interface{}) error {

	// Prepare transaction session
	tx := trd.db
	if itx != nil {
		itxdb, ok := itx.(*gorm.DB)
		if !ok {
			return errors.New("cannot assert argument itx to the *gorm.DB type")
		}

		tx = itxdb
	}

	// delete transaction
	tx = tx.Where(conds).Delete(&entity.TransactionLogInternal{})

	// Check error
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

// Count transaction log record by operator transaction id and reference.
func (trd *txnLogRepoDB) Count(serviceID uint8, ref string) (count int64) {
	trd.db.Model(&entity.TransactionLogInternal{}).
		Select("COUNT(1) AS count").
		Where("service_id = ?", serviceID).
		Where("reference = ?", ref).
		Count(&count)

	return
}
