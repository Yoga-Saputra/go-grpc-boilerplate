package repo

import (
	"errors"
	"strings"

	"github.com/Yoga-Saputra/go-grpc-boilerplate/app/entity"
	"gorm.io/gorm"
)

type (
	opTxnRepoDB struct {
		db *gorm.DB
	}
)

// NewOpTxnRepoDB create new DB repo for operator transaction entity.
func NewOpTxnRepoDB(db *gorm.DB) *opTxnRepoDB {
	if db != nil {
		return &opTxnRepoDB{db}
	}

	return nil
}

// Create new operator transaction record.
func (ord *opTxnRepoDB) Create(o *entity.RegisteredService) error {
	// o should != nil
	if o == nil {
		return errors.New("entity operator transaction argument pointer cannot be nil")
	}

	// Manually freed memory
	defer func() {
		o = nil
	}()

	// Insert
	if err := ord.db.Create(o).Error; err != nil {
		switch {
		case strings.Contains(err.Error(), "SQLSTATE 23505"):
			err = errors.New("operator transaction is already exists duplicate")
		}

		return err
	}

	return nil
}

// Finds operator transaction records by given conditions and return slice of operator transactions.
func (ord *opTxnRepoDB) Finds(conds map[string]interface{}) (res []entity.Wallet, err error) {
	tx := ord.db.Find(&res, conds)
	if tx.Error != nil {
		err = tx.Error
	} else if tx.RowsAffected <= 0 {
		err = errors.New("no records found")
	}
	return
}

// Find operator transaction records by given conditions and return operator transaction.
func (ord *opTxnRepoDB) Find(conds map[string]interface{}) (res entity.Wallet, err error) {
	tx := ord.db.Find(&res, conds)
	if tx.Error != nil {
		err = tx.Error
	} else if tx.RowsAffected <= 0 {
		err = errors.New("no records found")
	}

	return
}
