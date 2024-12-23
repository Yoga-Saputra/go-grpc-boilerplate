package entity

import (
	"math"
	"time"
)

type (
	// Wallet Promo entity.
	WalletPromo struct {
		ID             uint64 `gorm:"primaryKey" json:"id"`
		WalletID       uint64 `gorm:"not null;" json:"walletId"`
		InternalWallet Wallet `gorm:"foreignKey:WalletID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"internalWallet"`

		MemberID uint64 `gorm:"not null;uniqueIndex:wallet_promo_member_id_provider_code_p_idx" json:"memberId"`
		PID      string `gorm:"not null;size:25;uniqueIndex:wallet_promo_member_id_provider_code_p_idx;" json:"pId"`

		ServiceID        uint8             `gorm:"not null;" json:"serviceId"`
		InternalServices RegisteredService `gorm:"foreignKey:ServiceID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"internalServices"`

		// prepare promo by category or by provider
		ProviderCode string `gorm:"not null;size:5;uniqueIndex:wallet_promo_member_id_provider_code_p_idx;" json:"providerCode"`
		IsRunning    bool   `gorm:"not null; default:false;comment:true = running, false = expired" json:"is_running"` // status 0 = expired, 1 = running

		Currency    string    `gorm:"not null;size:5" json:"currency"`
		Amount      float64   `gorm:"type:numeric(30,6);default:0;not null" json:"amount"`
		ProcessedBy string    `gorm:"not null;size:25" json:"processed_by"`
		CreatedAt   time.Time `gorm:"not null" json:"createdAt"`
		UpdatedAt   time.Time `json:"updatedAt"`
	}
)

// TableName return table name of entity.
func (t *WalletPromo) TableName() string {
	return "wallet_promo"
}

func (w *WalletPromo) Amount2DecimalPlaces() float64 {
	return math.Floor(w.Amount*100) / 100
}
func (w *WalletPromo) Amount2DecimalPlacesAll(amount float64) float64 {
	return math.Floor((w.Amount+amount)*100) / 100
}
