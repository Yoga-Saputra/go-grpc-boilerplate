package repo

import (
	"errors"
	"log"
	"strings"

	"github.com/Yoga-Saputra/go-grpc-boilerplate/app/entity"
	"gorm.io/gorm"
)

type (
	WalletCreditOperator string

	walletRepoDB struct {
		db *gorm.DB
	}
)

// Enum for wallet credit opperations
const (
	WalletAddCreditOp       WalletCreditOperator = "+"
	WalletSubstractCreditOp WalletCreditOperator = "-"
	WalletMultiplyCreditOp  WalletCreditOperator = "*"
	WalletDivideCreditOp    WalletCreditOperator = "/"
)

// NewWalletRepoDB create new DB repo for wallet entity.
func NewWalletRepoDB(db *gorm.DB) *walletRepoDB {
	if db != nil {
		return &walletRepoDB{db}
	}

	return nil
}

// Validating operator contant enum
func (op WalletCreditOperator) validate() bool {
	valid := false

	switch op {
	case WalletAddCreditOp:
		valid = true

	case WalletSubstractCreditOp:
		valid = true

	case WalletMultiplyCreditOp:
		valid = true

	case WalletDivideCreditOp:
		valid = true
	}

	return valid
}

// Transaction repo method of wallet that approach DB process with transaction.
func (wrd *walletRepoDB) Transaction(txFunc func(interface{}) error) (err error) {
	tx := wrd.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if err != nil {
			log.Printf("[DBTxSession] - Rollback with reason: %s", err.Error())
			tx.Rollback()
		} else {
			err = tx.Commit().Error
			if err != nil {
				log.Printf("[DBTxSession] - Commit error: %s", err.Error())
			}
		}
	}()

	err = txFunc(tx)
	return
}

// Create new wallet record.
func (wrd *walletRepoDB) Create(w *entity.Wallet) error {
	// w should != nil
	if w == nil {
		return errors.New("entity wallet argument pointer cannot be nil")
	}

	// Manually freed memory
	defer func() {
		w = nil
	}()

	// Insert
	if err := wrd.db.Omit(
		"IsNew",
		"UpdatedAt",
	).Create(w).Error; err != nil {
		switch {
		case strings.Contains(err.Error(), "SQLSTATE 23505"):
			err = errors.New("member wallet is already exists duplicate")
		}

		return err
	}

	return nil
}

// Find wallet records by given conditions and return wallet.
func (wrd *walletRepoDB) Find(conds map[string]interface{}) (res entity.Wallet, rows int, err error) {
	tx := wrd.db.Find(&res, conds)
	rows = int(tx.RowsAffected)
	err = tx.Error

	return
}

// Find wallet promo record by given conditions and return wallet.
func (wprd *walletRepoDB) FindWalletPromo(conds map[string]interface{}) (res entity.WalletPromo, rows int, err error) {
	tx := wprd.db.Where(conds).Find(&res)
	rows = int(tx.RowsAffected)
	err = tx.Error

	return
}
