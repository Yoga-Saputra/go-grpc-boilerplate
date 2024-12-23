package txnlogprovider

import (
	"errors"
	"fmt"
	"time"

	"github.com/Yoga-Saputra/go-grpc-boilerplate/app/entity"
)

// Will be exposed struct of txnlogprovider.
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

// Insert the transaction log provider record.
func Insert(
	prov string,
	l *entity.TransactionLogProvider,
	itx interface{},
) error {
	if err := validatePointer(); err != nil {
		return err
	}

	return public.meta.repo.Create(prov, l, itx)
}

// CheckCount the transaction log provider by reference, pId and bDate.
func CheckCount(
	prov, ref, pId, bDate string,
) error {
	count := public.meta.repo.Count(
		prov,
		ref,
		pId,
		bDate,
	)
	if count > 0 {
		return fmt.Errorf(
			"duplicate transaction provider: %s, reference: %s, pId: %s, bDate: %s",
			prov,
			ref,
			pId,
			bDate,
		)
	}

	return nil
}

// Check the transaction log provider by reference, pId and bDate.
func CheckOneTxnID(
	prov, ref, pId, bDate string,
	bt time.Time,
) (entity.TransactionLogProvider, bool) {
	return public.meta.repo.DeterminedFind(
		prov,
		ref,
		pId,
		bDate,
		bt,
	)
}

// Check the transaction log provider by one or more reference, pId and bDate.
func CheckManyTxnID(
	refSingle string,
	prov, pId, bDate string,
	bt time.Time,
	refDiff ...string,
) ([]entity.TransactionLogProvider, bool) {
	return public.meta.repo.DeterminedFinds(
		refSingle,
		prov,
		pId,
		bDate,
		bt,
		refDiff...,
	)
}

// Check the transaction log provider by reference, pId.
func CheckOneTxnIDWithoutDate(
	prov, ref, pId string,
) (entity.TransactionLogProvider, bool) {
	return public.meta.repo.DeterminedFindWithoutDate(
		prov,
		ref,
		pId,
	)
}

// Check the transaction log provider only by reference.
func CheckOneTxnIDWithoutDateOnlyByReference(
	prov, ref string,
) (entity.TransactionLogProvider, bool) {
	return public.meta.repo.DeterminedFindWithoutDateOnlyByReference(
		prov,
		ref,
	)
}

// GetDataByTicketId return log provider by given ticket id and date time.
func GetDataByTicketId(
	prov, ticketId string,
	bt time.Time,
) (
	[]entity.TransactionLogProvider,
	error,
) {
	return public.meta.repo.GetDataByTicketId(prov, ticketId, bt)
}
