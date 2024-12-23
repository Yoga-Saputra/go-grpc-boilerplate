package mcslog

import (
	"time"

	"github.com/Yoga-Saputra/go-grpc-boilerplate/app/entity"
)

type (
	Reppository interface {
		Insert(a *entity.Transfer) error
		GetServiceByID(id uint8) (*entity.RegisteredService, error)
	}

	Status uint8

	TrxType uint8

	Param struct {
		TrxID            string
		BranchCode       int16
		PId              string
		Username         string
		Currency         string
		Amount           float64
		AmountBefore     float64
		AmountAfter      float64
		Date             time.Time
		Status           Status
		WalletCode       string
		TrxType          TrxType
		WalletAdditional AdditionalWalletInfo
		ProcessedBy      string
	}

	AdditionalWalletInfo struct {
		Code         string
		BeforeCredit float64
		AfterCredit  float64
	}
)

const (
	SUCCESS Status = iota
	FAILED
	CANCELED
)

const (
	ADD_CREDIT TrxType = iota
	DEDUCT_CREDIT
)

func (s Status) String() string {
	switch s {
	default:
		return "UNKNOWN"
	case SUCCESS:
		return "SUCCESS"
	case FAILED:
		return "FAILED"
	case CANCELED:
		return "CANCELED"
	}
}
