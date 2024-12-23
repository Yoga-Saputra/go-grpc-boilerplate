package entity

import (
	"database/sql"
	"math"
	"time"
)

type (
	// To hold wallet category value.
	WalletCategory string

	// Wallet entity.
	Wallet struct {
		ID       uint64 `gorm:"primaryKey" json:"id"`
		BranchID uint16 `gorm:"not null;index" json:"branchId"`

		MemberID uint64         `gorm:"not null;index;uniqueIndex:wallets_member_id_category_idx" json:"memberId"`
		Category WalletCategory `gorm:"not null;size:25;uniqueIndex:wallets_member_id_category_idx;uniqueIndex:wallets_category_p_id_idx" json:"category"`

		PID      string         `gorm:"not null;index;size:25;uniqueIndex:wallets_category_p_id_idx" json:"pId"`
		Currency string         `gorm:"not null;size:5" json:"currency"`
		Username sql.NullString `gorm:"size:50;" json:"username"`

		Amount        float64 `gorm:"type:numeric(30,6);default:0;not null" json:"amount"`
		NetProfitLoss float64 `gorm:"type:numeric(30,6);default:0;not null" json:"netProfitLoss"`

		IsNew      bool `gorm:"default:true;not null" json:"isNew"`
		LockedIn   bool `gorm:"default:true;not null" json:"lockedIn"`
		LockedOut  bool `gorm:"default:true;not null" json:"lockedOut"`
		IsLocked   bool `gorm:"default:false;not null" json:"isLocked"`
		IsDisabled bool `gorm:"default:false;not null" json:"isDisabled"`

		CreatedAt time.Time `gorm:"not null" json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	}

	WalletMemberDataTable struct {
		WalletID       uint64  `json:"wallet_id"`
		MemberID       uint64  `json:"member_id"`
		BranchID       uint16  `json:"branch_id"`
		Currency       string  `json:"currency"`
		PID            string  `json:"p_id"`
		Username       string  `json:"username"`
		CommonLocked   bool    `json:"common_locked"`
		CommonDisabled bool    `json:"common_disabled"`
		PromoRunning   bool    `json:"promo_running"`
		CreatedAt      string  `json:"created_at"`
		AmountCommon   float64 `json:"amount_common"`
		AmountPromo    float64 `json:"amount_promo"`
		ProviderCode   string  `json:"provider_code"`
	}
	MemberWalletSummary struct {
		Count        int64   `json:"count"`
		AmountPromo  float64 `json:"amount_promo"`
		AmountCommon float64 `json:"amount_common"`
	}
)

// Enum of wallet category.
const (
	COMMON WalletCategory = "COMMON"
)

// Amount2DecimalPlaces return credit amount that already round down to 2 decimal places.
func (w *Wallet) Amount2DecimalPlaces() float64 {
	return math.Floor(w.Amount*100) / 100
}
