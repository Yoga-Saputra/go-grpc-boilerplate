package repo

import (
	"errors"
	"strings"

	"github.com/Yoga-Saputra/go-grpc-boilerplate/app/entity"
	"gorm.io/gorm"
)

type transferLogRepoDB struct {
	db *gorm.DB
}

func NewTransferLogRepoDB(db *gorm.DB) (*transferLogRepoDB, error) {
	if db == nil {
		return nil, errors.New("pointer argument is nil")
	}

	return &transferLogRepoDB{db: db}, nil
}

func (tlrd *transferLogRepoDB) Insert(a *entity.Transfer) error {
	if a == nil {
		return errors.New("pointer argument cannot be nil")
	}

	err := tlrd.db.Create(a).Error
	if err != nil {
		if strings.Contains(err.Error(), "SQLSTATE 23505") {
			return nil
		}
	}
	return err
}

func (tlrd *transferLogRepoDB) GetServiceByID(id uint8) (
	data *entity.RegisteredService,
	err error,
) {
	err = tlrd.db.Take(&data, id).Error
	return
}
