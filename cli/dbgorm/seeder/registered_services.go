package seeder

import (
	"strings"

	"github.com/Yoga-Saputra/go-grpc-boilerplate/app/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RegisteredServices struct{}

var regSvcDatas = []entity.RegisteredService{
	{ID: 1, Code: "topup"},
	{ID: 2, Code: "withdraw"},
	{ID: 3, Code: "adjustment"},
	{ID: 4, Code: "claim-promo"},
	{ID: 5, Code: "referral"},
	{ID: 6, Code: "voucher"},
	{ID: 7, Code: "cashback"},
}

func (RegisteredServices) Run(refresh bool, tx *gorm.DB) error {
	// Modify values
	for i, d := range regSvcDatas {
		regSvcDatas[i].Code = strings.ToUpper(d.Code)
	}

	if err := tx.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&regSvcDatas).Error; err != nil {
		return err
	}
	return nil
}
