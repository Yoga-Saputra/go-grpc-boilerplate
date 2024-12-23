package entity

import "time"

type Transfer struct {
	ID                int       `gorm:"primaryKey" json:"id"`
	TrxID             string    `gorm:"column:trxId" json:"trxId"`
	BranchCode        int16     `gorm:"column:branchCode" json:"branchCode"`
	PID               string    `gorm:"column:pId" json:"pId"`
	Username          string    `gorm:"column:username" json:"username"`
	WalletOrigin      string    `gorm:"column:walletOrigin" json:"walletOrigin"`
	WalletDestination string    `gorm:"column:walletDestination" json:"walletDestination"`
	RetryAttempts     int16     `gorm:"column:retryAttempts" json:"retryAttempts"`
	Amount            float64   `gorm:"column:amount" json:"amount"`
	Currency          string    `gorm:"column:currency" json:"currency"`
	TransactionDate   time.Time `gorm:"column:transactionDate" json:"transactionDate"`
	TransactionStatus string    `gorm:"column:transactionStatus" json:"transactionStatus"`
	TransactionType   string    `gorm:"column:transactionType" json:"transactionType"`
	Description       string    `gorm:"column:description" json:"description"`
}

func (t Transfer) TableName() string {
	return "transfer"
}
