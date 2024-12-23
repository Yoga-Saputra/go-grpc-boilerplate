package entity

import "time"

type (
	// TransactionLogInternal entity.
	TransactionLogInternal struct {
		ID               uint64            `gorm:"primaryKey" json:"id"`
		ServiceID        uint8             `gorm:"not null;uniqueIndex:transaction_logs_service_id_reference_idx" json:"serviceId"`
		InternalServices RegisteredService `gorm:"foreignKey:ServiceID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"internalServices"`
		Reference        string            `gorm:"not null;size:36;index;uniqueIndex:transaction_logs_service_id_reference_idx" json:"reference"`

		PID      string `gorm:"not null;size:25;index" json:"pId"`
		Currency string `gorm:"not null;size:5" json:"currency"`

		TransferAmount float64 `gorm:"type:numeric(30,6);default:0;not null" json:"transferAmount"`
		BalanceBefore  float64 `gorm:"type:numeric(30,6);default:0;not null" json:"balanceBefore"`
		BalanceAfter   float64 `gorm:"type:numeric(30,6);default:0;not null" json:"balanceAfter"`

		Date      string    `gorm:"->;not null;index;type:date;default:now()" json:"date"`
		CreatedAt time.Time `gorm:"not null" json:"createdAt"`
	}
)
