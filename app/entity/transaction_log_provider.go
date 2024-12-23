package entity

import (
	"database/sql"
	"strings"
	"time"

	"gorm.io/gorm"
)

type (
	// To hold provider transaction type.
	TxnProvType string

	// TransactionLog of provider ralted entity.
	TransactionLogProvider struct {
		ID        uint64         `gorm:"primaryKey" json:"id"`
		Reference string         `gorm:"not null;size:50;index:,unique,composite:reference_p_id_date_key" json:"reference"`
		TicketID  sql.NullString `gorm:"size:50;index:,where:ticket_id Is NOT NULL" json:"ticketId"`
		GameCode  sql.NullString `gorm:"size:120" json:"gameCode"`

		MemberID sql.NullInt64  `json:"memberId"`
		PID      string         `gorm:"not null;size:35;index:,unique,composite:reference_p_id_date_key" json:"pId"`
		Username sql.NullString `gorm:"size:50;" json:"username"`
		Currency string         `gorm:"not null;size:5" json:"currency"`
		BranchID sql.NullInt16  `json:"branchId"`

		GCategory   string  `gorm:"size:5" json:"gCategory"`
		Progressive float64 `gorm:"type:numeric(30,6)" json:"progressive"`
		JAmount     float64 `gorm:"type:numeric(30,6)" json:"jAmount"`
		BAmount     float64 `gorm:"type:numeric(30,6)" json:"bAmount"`

		WlAmount       float64        `gorm:"type:numeric(30,6)" json:"wlAmount"`
		PAmount        float64        `gorm:"type:numeric(30,6)" json:"pAmount"`
		TransferAmount float64        `gorm:"type:numeric(30,6);default:0;not null" json:"transferAmount"`
		BalanceBefore  float64        `gorm:"type:numeric(30,6);default:0;not null" json:"balanceBefore"`
		BalanceAfter   float64        `gorm:"type:numeric(30,6);default:0;not null" json:"balanceAfter"`
		Wl             float64        `gorm:"type:numeric(30,6);default:0;not null" json:"wl"`
		BDateTime      sql.NullString `gorm:"size:22" json:"bDateTime"`

		TxnType   TxnProvType `gorm:"size:10" json:"txnType"`
		Date      string      `gorm:"primaryKey;not null;index:,unique,composite:reference_p_id_date_key;type:date" json:"date"`
		CreatedAt time.Time   `gorm:"not null" json:"createdAt"`
	}
)

// Enum of provider transaction type.
const (
	OTHER   TxnProvType = "OTHER"
	TRFRIN  TxnProvType = "TRFRIN"
	TRFROUT TxnProvType = "TRFROUT"
)

// TableName return table name of entity.
func (t *TransactionLogProvider) TableName() string {
	return "transaction_logs"
}

// ProvSchemaTable return dynamic table schema name based on provider code as schema.
func ProvSchemaTable(t *TransactionLogProvider, prov string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		prov = strings.ToLower(prov)
		tableName := prov + "." + t.TableName()
		return db.Table(tableName)
	}
}

// ProvSchemaTablePartmanTemplate return dynamic table partman template schema name based on provider code as schema.
func ProvSchemaTablePartmanTemplate(t *TransactionLogProvider, prov string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		prov = strings.ToLower(prov)
		tableName := "partman.template_" + prov + "_" + t.TableName()
		return db.Table(tableName)
	}
}

// ProvDynamicTable return dynamic table name based on provider code as schema
// and month as a prefix.
func ProvDynamicTable(t *TransactionLogProvider, month, prov string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		prov = strings.ToLower(prov)
		month = strings.ToLower(month)
		tableName := prov + "." + month + "_" + t.TableName()
		return db.Table(tableName)
	}
}
