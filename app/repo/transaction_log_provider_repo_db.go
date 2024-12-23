package repo

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Yoga-Saputra/go-grpc-boilerplate/app/entity"
	"gorm.io/gorm"
)

type (
	txnLogProviderRepoDB struct {
		db *gorm.DB
		tz *time.Location
	}
)

// Return gorm DB scopes.
func scope(
	db *gorm.DB,
	tz *time.Location,
	provider string,
) *gorm.DB {
	log.Println(tz.String())
	return db.Scopes(entity.ProvSchemaTable(&entity.TransactionLogProvider{}, provider))
}

// NewTxnLogRepoDB create new DB repo for transaction log entity.
func NewTxnProviderLogRepoDB(db *gorm.DB, tz *time.Location) *txnLogProviderRepoDB {
	if db != nil {
		return &txnLogProviderRepoDB{db, tz}
	}

	return nil
}

// Create new txn log provider record.
func (t *txnLogProviderRepoDB) Create(
	prov string,
	l *entity.TransactionLogProvider,
	itx interface{},
) error {
	// l should != nil
	if l == nil {
		return errors.New("entity transaction log provider argument pointer cannot be nil")
	}
	defer func() {
		l = nil
	}()

	// Get scope
	tx := scope(t.db, t.tz, prov)

	// Prepare transaction session
	if itx != nil {
		itxdb, ok := itx.(*gorm.DB)
		if !ok {
			return errors.New("cannot assert argument itx to the *gorm.DB type")
		}

		tx = scope(itxdb, t.tz, prov)
	}

	// Insert
	if err := tx.Create(l).Error; err != nil {
		switch {
		case strings.Contains(err.Error(), "SQLSTATE 23505"):
			err = fmt.Errorf("duplicate transaction provider: %s, reference: %s, pId: %s", prov, l.Reference, l.PID)
		}

		return err
	}

	return nil
}

// Count transaction log provider record by reference and pId.
func (t *txnLogProviderRepoDB) Count(
	prov, ref, pId, bDate string,
) (count int64) {
	tx := scope(t.db, t.tz, prov)

	tx.Select("COUNT(1) AS count").
		Where("reference = ?", ref).
		Where("p_id = ?", pId).
		Where("date = ?", bDate).
		Count(&count)

	return
}

// DeterminedFind to select transaction log provider record by reference, pId and date.
func (t *txnLogProviderRepoDB) DeterminedFind(
	prov, ref, pId, bDate string,
	bt time.Time,
) (e entity.TransactionLogProvider, ok bool) {
	// Query
	result := scope(t.db, t.tz, prov).
		Where("reference = ?", ref).
		Where("p_id = ?", pId)

	// Check if b date time is early day with max 5 minutes ahead
	// otherwise query where equal single date
	if bt.Hour() == 0 && bt.Minute() <= 5 {
		dayBeforeStr := bt.AddDate(0, 0, -1).Format("2006-01-02")
		result = result.Where("date BETWEEN ? AND ?", dayBeforeStr, bDate)
	} else {
		result = result.Where("date = ?", bDate)
	}

	result.First(&e)
	if result.RowsAffected <= 0 {
		return
	}

	// Finish
	ok = true
	return
}

// DeterminedFinds same like DeterminedFind(), but ref/transaction id slice of string
// and will return slice/array.
func (t *txnLogProviderRepoDB) DeterminedFinds(
	refSingle string,
	prov, pId, bDate string,
	bt time.Time,
	refDiff ...string,
) (e []entity.TransactionLogProvider, ok bool) {
	// Query
	result := scope(t.db, t.tz, prov)

	// Check if have argument many ref/trx_id
	if len(refDiff) > 0 {
		result = result.Where("reference IN ?", refDiff)
	} else {
		result = result.Where("reference = ?", refSingle)
	}
	result = result.Where("p_id = ?", pId)

	// Check if b date time is early day with max 5 minutes ahead
	// otherwise query where equal single date
	if bt.Hour() == 0 && bt.Minute() <= 5 {
		dayBeforeStr := bt.AddDate(0, 0, -1).Format("2006-01-02")
		result = result.Where("date BETWEEN ? AND ?", dayBeforeStr, bDate)
	} else {
		result = result.Where("date = ?", bDate)
	}

	result.First(&e)
	if result.RowsAffected <= 0 {
		return
	}

	// Finish
	ok = true
	return
}

// DeterminedFind to select transaction log provider record by reference, pId.
func (t *txnLogProviderRepoDB) DeterminedFindWithoutDate(
	prov, ref, pId string,
) (e entity.TransactionLogProvider, ok bool) {
	// Query
	result := scope(t.db, t.tz, prov).
		Where("reference = ?", ref).
		Where("p_id = ?", pId)

	result.First(&e)
	if result.RowsAffected <= 0 {
		return
	}

	// Finish
	ok = true
	return
}

// DeterminedFindWithoutDateOnlyByReference to select transaction log provider record only by reference.
func (t *txnLogProviderRepoDB) DeterminedFindWithoutDateOnlyByReference(
	prov, ref string,
) (e entity.TransactionLogProvider, ok bool) {
	// Query
	result := scope(t.db, t.tz, prov).
		Where("reference = ?", ref)

	result.First(&e)
	if result.RowsAffected <= 0 {
		return
	}

	// Finish
	ok = true
	return
}

// GetDataByTicketId select data by ticket id and date time.
func (t *txnLogProviderRepoDB) GetDataByTicketId(
	prov, ticketId string,
	bt time.Time,
) (
	datas []entity.TransactionLogProvider,
	err error,
) {
	// Query
	result := scope(t.db, t.tz, prov).Where("ticket_id = ?", ticketId)

	// Check if b date time is early day with max 5 minutes ahead
	// otherwise query where equal single date
	bDate := bt.Format("2006-01-02")
	if bt.Hour() == 0 && bt.Minute() <= 5 {
		dayBeforeStr := bt.AddDate(0, 0, -1).Format("2006-01-02")
		result = result.Where("date BETWEEN ? AND ?", dayBeforeStr, bDate)
	} else {
		result = result.Where("date = ?", bDate)
	}

	tx := result.Find(&datas)
	if tx.Error != nil {
		err = tx.Error
	}

	// check if count datas is 1, will use query beetwen sub 1 month from payout date
	if len(datas) == 1 {
		reuslt1 := scope(t.db, t.tz, prov).Where("ticket_id = ?", ticketId)
		monthBeforeStr := bt.AddDate(0, -1, 0).Format("2006-01-02")

		tx1 := reuslt1.Where("date BETWEEN ? AND ?", monthBeforeStr, bDate).Find(&datas)
		if tx1.Error != nil {
			err = tx1.Error
		}
	}

	return
}
