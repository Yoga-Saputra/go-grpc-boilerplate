package mcslog

import (
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/Yoga-Saputra/go-grpc-boilerplate/app/entity"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var public *Meta

// Validating local variable pointer.
func validatePointer() error {
	if public == nil {
		return fmt.Errorf("public API for mcslog package cannot be accessed, not implemented yet")
	}

	return nil
}

// Insert param to the mcslog data.
func Insert(
	param *Param,
	service *entity.RegisteredService,
	provider ...string,
) error {
	// Validation
	if err := validatePointer(); err != nil {
		return err
	}
	switch {
	case param == nil:
		return errors.New("pointer argument param cannot be nil")
	case service == nil:
		return errors.New("pointer argument service cannot be nil")
	case len(strings.TrimSpace(service.Code)) == 0:
		return errors.New("argument service.Code cannot be empty")
	case len(strings.TrimSpace(param.WalletCode)) == 0:
		return errors.New("argument param.WalletCode cannot be empty")
	}

	// Create origin and destination
	var origin, destination string
	if param.TrxType == ADD_CREDIT {
		origin = strings.ToUpper(service.Code)
		if len(provider) > 0 {
			origin = fmt.Sprintf("%s %s", origin, strings.ToUpper(provider[0]))
		}

		destination = strings.ToUpper(param.WalletCode)
		if param.Amount < 0 {
			param.Amount = math.Abs(param.Amount)
		}
	} else {
		origin = strings.ToUpper(param.WalletCode)
		destination = strings.ToUpper(service.Code)
		if len(provider) > 0 {
			destination = fmt.Sprintf("%s %s", destination, strings.ToUpper(provider[0]))
		}

		if param.Amount > 0 {
			param.Amount = param.Amount * -1
		}
	}

	// Create description
	desc := fmt.Sprintf("Before (%s), After (%s)",
		AmountFormat(param.AmountBefore),
		AmountFormat(param.AmountAfter),
	)
	if len(strings.TrimSpace(param.WalletAdditional.Code)) > 0 {
		desc += fmt.Sprintf(". %s: Before (%s), After(%s)",
			strings.ToUpper(param.WalletAdditional.Code),
			AmountFormat(param.WalletAdditional.BeforeCredit),
			AmountFormat(param.WalletAdditional.AfterCredit),
		)
	}
	if len(strings.TrimSpace(param.ProcessedBy)) > 0 {
		desc += fmt.Sprintf(". Processed by: %s", param.ProcessedBy)
	}

	// Set data
	data := &entity.Transfer{
		TrxID:             param.TrxID,
		BranchCode:        param.BranchCode,
		PID:               param.PId,
		Username:          param.Username,
		WalletOrigin:      origin,
		WalletDestination: destination,
		RetryAttempts:     0,
		Amount:            param.Amount,
		Currency:          param.Currency,
		TransactionDate:   param.Date.In(public.tz),
		TransactionStatus: param.Status.String(),
		TransactionType:   strings.ToUpper(service.Code),
		Description:       desc,
	}

	// Insert
	return public.repo.Insert(data)
}

// GetRegisteredServiceByID select registered service by id.
func GetRegisteredServiceByID(id uint8) (*entity.RegisteredService, error) {
	return public.repo.GetServiceByID(id)
}

func AmountFormat(amount float64) string {
	p := message.NewPrinter(language.English)
	return p.Sprintf("%.2f", amount)
}
