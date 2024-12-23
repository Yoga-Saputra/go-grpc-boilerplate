package rpcwalletinsys

import (
	"context"
	"fmt"
	"strings"

	"github.com/Yoga-Saputra/go-grpc-boilerplate/app/contract"
	"github.com/Yoga-Saputra/go-grpc-boilerplate/app/job/createwallet"
	"github.com/Yoga-Saputra/go-grpc-boilerplate/pkg/sw_pb_go/wallet/v1/insys/common"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// CreateWallet GRpc implemented method to create new wallet record.
func (ics *inSysCommonServer) CreateWallet(
	c context.Context,
	p *common.CreateWallet_Req,
) (*common.CreateWallet_Res, error) {
	// Prepare response
	resp := &common.CreateWallet_Res{
		Success: true,
		Code:    int32(contract.OK),
	}
	defer func() {
		resp = nil
	}()

	// Validate payload
	if err := ics.createWalletValidate(p); err != nil {
		ics.xloger("CreateWallet", "Validation Error:", err.String())

		resp.Success = false
		resp.Code = int32(err.Code)
		resp.Error = &wrapperspb.StringValue{Value: err.String()}
		return resp, nil
	}

	// Enqueue task to the queue
	info, e := createwallet.Enqueue(&createwallet.QPayload{
		BranchID: uint16(p.BranchId),
		MemberID: p.MemberId,
		PID:      p.PId,
		Currency: p.Currency,
		Username: p.Username,
	})

	if e != nil {
		ics.xloger("CreateWallet", "Error enqueue task:", e.Error())

		// What error task
		if strings.Contains(e.Error(), "exists") {
			resp.Code = int32(contract.QUEUETASKEXISTS)
			resp.Error = &wrapperspb.StringValue{Value: contract.QUEUETASKEXISTS.String()}
		} else {
			resp.Code = int32(contract.INTERNALERROR)
			resp.Error = &wrapperspb.StringValue{Value: contract.INTERNALERROR.String()}
		}

		resp.Success = false
		return resp, nil
	}

	resp.Data = fmt.Sprintf("Create Wallet task already enqueued with id: %s", info.ID)

	return resp, nil
}
