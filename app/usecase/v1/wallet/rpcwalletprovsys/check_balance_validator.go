package rpcwalletprovsys

import (
	"strings"

	"github.com/Yoga-Saputra/go-grpc-boilerplate/app/contract"
	"github.com/Yoga-Saputra/go-grpc-boilerplate/pkg/sw_pb_go/wallet/v1/provsys/syncp"
)

// Validator "CheckBalance".
func (syc *provSysSyncpServer) checkBalanceValidate(p *syncp.SyncCredit_Req) (err *contract.Error) {
	switch {
	case len(strings.TrimSpace(p.GetPId())) <= 0:
		err = &contract.Error{Code: contract.VALIDATIONERROR, AppendFormat: []string{"payload p_id is required"}}

	case len(strings.TrimSpace(p.GetProviderCode())) <= 0:
		err = &contract.Error{Code: contract.VALIDATIONERROR, AppendFormat: []string{"payload provider_code is required"}}

	case len(strings.TrimSpace(p.GetProviderCode())) < 2 || len(strings.TrimSpace(p.GetProviderCode())) > 5:
		err = &contract.Error{Code: contract.VALIDATIONERROR, AppendFormat: []string{"payload provider_code is must be min 2 char and max 5 char"}}
	}

	return
}
