package entity

type (
	// RegisteredService entity.
	RegisteredService struct {
		ID   uint8  `gorm:"primaryKey" json:"id"`
		Code string `gorm:"not null;size:100;unique" json:"code"`
	}
)
