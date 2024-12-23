package rpcwalletinsys

import (
	"encoding/json"
	"strings"

	"github.com/Yoga-Saputra/go-grpc-boilerplate/app/contract"
	"github.com/Yoga-Saputra/go-grpc-boilerplate/pkg/sw_pb_go/wallet/v1/insys/common"
)

// Validator CreateWallet
func (ics *inSysCommonServer) createWalletValidate(p *common.CreateWallet_Req) (err *contract.Error) {
	// Log payload
	b, e := json.Marshal(p)
	if err != nil {
		ics.xloger("CreateWallet", "Error marshal payload:", e.Error())
	} else {
		ics.xloger("CreateWallet", "Payload:", string(b))
	}

	switch {
	case p.GetBranchId() <= 0:
		err = &contract.Error{Code: contract.VALIDATIONERROR, AppendFormat: []string{"payload branch_id is required"}}

	case p.GetMemberId() <= 0:
		err = &contract.Error{Code: contract.VALIDATIONERROR, AppendFormat: []string{"payload member_id is required"}}

	case len(strings.TrimSpace(p.GetPId())) <= 0:
		err = &contract.Error{Code: contract.VALIDATIONERROR, AppendFormat: []string{"payload p_id is required"}}

	case len(strings.TrimSpace(p.GetCurrency())) <= 0:
		err = &contract.Error{Code: contract.VALIDATIONERROR, AppendFormat: []string{"payload currency is required"}}

	case len(strings.TrimSpace(p.GetCurrency())) < 2 || len(strings.TrimSpace(p.GetCurrency())) > 5:
		err = &contract.Error{Code: contract.VALIDATIONERROR, AppendFormat: []string{"payload currency is must be min 2 char and max 5 char"}}

	case len(strings.TrimSpace(p.GetUsername())) <= 0:
		err = &contract.Error{Code: contract.VALIDATIONERROR, AppendFormat: []string{"payload username is required"}}
	}

	return
}
