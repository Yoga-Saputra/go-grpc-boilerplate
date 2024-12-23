package txnlogprovider

import (
	"time"

	"github.com/Yoga-Saputra/go-grpc-boilerplate/app/entity"
)

type (
	Repository interface {
		Create(
			prov string,
			l *entity.TransactionLogProvider,
			itx interface{},
		) error

		Count(
			prov, ref, pId, bDate string,
		) (count int64)
		DeterminedFind(
			prov, ref, pId, bDate string,
			bt time.Time,
		) (e entity.TransactionLogProvider, ok bool)
		DeterminedFinds(
			refSingle string,
			prov, pId, bDate string,
			bt time.Time,
			refDiff ...string,
		) (e []entity.TransactionLogProvider, ok bool)
		DeterminedFindWithoutDate(
			prov, ref, pId string,
		) (e entity.TransactionLogProvider, ok bool)
		DeterminedFindWithoutDateOnlyByReference(
			prov, ref string,
		) (e entity.TransactionLogProvider, ok bool)

		GetDataByTicketId(
			prov, ticketId string,
			bt time.Time,
		) (
			datas []entity.TransactionLogProvider,
			err error,
		)
	}
)
