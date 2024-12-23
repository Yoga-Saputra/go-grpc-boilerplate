package wallet

import (
	"github.com/Yoga-Saputra/go-grpc-boilerplate/app/entity"
)

type (
	Repository interface {
		Transaction(txFunc func(interface{}) error) (err error)
		Create(w *entity.Wallet) error

		Find(conds map[string]interface{}) (res entity.Wallet, rows int, err error)
		FindWalletPromo(conds map[string]interface{}) (res entity.WalletPromo, rows int, err error)
	}
)
