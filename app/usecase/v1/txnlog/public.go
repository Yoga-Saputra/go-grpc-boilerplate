package txnlog

import (
	"errors"
	"fmt"

	"github.com/Yoga-Saputra/go-grpc-boilerplate/app/entity"
)

// Will be exposed struct of txnlog.
type publicAPI struct {
	meta *Meta
}

// Local variable of exposed txnlog struct.
var public *publicAPI

// Create new exposed txnlog instances.
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

// Check the transaction log by service id and reference.
func Check(opTxnID uint8, ref string) error {
	count := public.meta.repo.Count(opTxnID, ref)
	if count > 0 {
		return fmt.Errorf("duplicate transaction with serviceId: %d, reference: %s", opTxnID, ref)
	}

	return nil
}

// Insert the transaction log record.
func Insert(l *entity.TransactionLogInternal, itx interface{}) error {
	if err := validatePointer(); err != nil {
		return err
	}

	return public.meta.repo.Create(l, itx)
}

// Delete the transaction log record.
func Delete(conds map[string]interface{}, itx interface{}) error {
	if err := validatePointer(); err != nil {
		return err
	}

	return public.meta.repo.Delete(conds, itx)
}
