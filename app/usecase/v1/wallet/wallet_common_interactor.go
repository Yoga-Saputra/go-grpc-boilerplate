package wallet

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/Yoga-Saputra/go-grpc-boilerplate/app/entity"
)

// CreateNewWallet record
func (m *Meta) CreateNewWallet(
	branchID uint16,
	memberID uint64,
	pId, currency, username string,
) error {
	walletCategories := []entity.WalletCategory{entity.COMMON}

	// Do insert wallet loop through categories
	for _, v := range walletCategories {
		// Prepare wallet pointer
		w := &entity.Wallet{
			BranchID: branchID,
			MemberID: memberID,
			Category: v,
			PID:      pId,
			Username: sql.NullString{Valid: true, String: username},
			Currency: currency,
		}
		if err := m.repo.Create(w); err != nil {
			return err
		}
	}

	return nil
}

// GetWalletByMember record
func (m *Meta) GetWalletByMember(
	memberID,
	pId interface{},
	category entity.WalletCategory,
) (*entity.Wallet, int, error) {
	// Prepare query conditions
	conds := make(map[string]interface{})
	conds["category"] = category
	switch {
	case memberID != nil:
		conds["member_id"] = memberID

	case pId != nil:
		conds["p_id"] = pId

	default:
		return nil, 0, fmt.Errorf("argument memberID nor pId cannot be empty -> memberID: %d, pId: %s", memberID, pId)
	}

	// Select record
	wallet, rows, err := m.repo.Find(conds)
	switch {
	// Always check error appear first
	case err != nil:
		return nil, rows, err

	case rows == 0:
		return nil, rows, nil

	default:
		return &wallet, rows, nil
	}
}

// GetWalletByProviderCode record
func (m *Meta) GetWalletPromoByProviderCode(
	memberID,
	pId interface{},
	serviceId interface{},
	isRunning bool,
	providerCode string,
) (*entity.WalletPromo, int, error) {
	// Prepare query conditions
	conds := make(map[string]interface{})

	conds["is_running"] = isRunning
	conds["provider_code"] = strings.ToLower(providerCode)
	switch {
	case memberID != nil:
		conds["member_id"] = memberID

	case pId != nil:
		conds["p_id"] = pId

	case serviceId != nil:
		conds["service_id"] = serviceId

	default:
		return nil, 0, fmt.Errorf("argument memberID nor pId cannot be empty -> memberID: %d, pId: %s, isRunning : %v", memberID, pId, isRunning)
	}

	// Select record
	wPromo, rows, err := m.repo.FindWalletPromo(conds)
	switch {
	// Always check error appear first
	case err != nil:
		return nil, rows, err

	case rows == 0:
		return nil, rows, nil

	default:
		return &wPromo, rows, nil
	}
}
