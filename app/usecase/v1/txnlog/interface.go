package txnlog

import "github.com/Yoga-Saputra/go-grpc-boilerplate/app/entity"

type (
	Repository interface {
		Create(t *entity.TransactionLogInternal, itx interface{}) error

		Delete(conds map[string]interface{}, itx interface{}) error

		Count(opTxnID uint8, ref string) (count int64)
	}
)
