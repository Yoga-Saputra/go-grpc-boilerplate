package wallet

import (
	"errors"
)

type (
	// Will be exposed struct of wallet.
	publicAPI struct {
		meta *Meta
	}
)

// Local variable of exposed wallet struct.
var public *publicAPI

// Create new exposed wallet instances.
func newPublicAPI(meta *Meta) {
	public = &publicAPI{meta}
}

// Validating local variable pointer.
func validatePointer() error {
	if public == nil {
		return errors.New("package cannot be accessed, please implement the usecase first")
	}

	return nil
}

func CreateWallet(
	BranchID uint16,
	MemberID uint64,
	pId, Currency, Username string,
) error {
	if err := validatePointer(); err != nil {
		return err
	}

	return public.meta.CreateNewWallet(BranchID, MemberID, pId, Currency, Username)
}
