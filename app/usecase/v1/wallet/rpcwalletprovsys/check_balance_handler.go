package rpcwalletprovsys

import (
	"context"

	"github.com/Yoga-Saputra/go-grpc-boilerplate/app/contract"
	"github.com/Yoga-Saputra/go-grpc-boilerplate/app/entity"
	"github.com/Yoga-Saputra/go-grpc-boilerplate/pkg/sw_pb_go/wallet/v1/provsys/syncp"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// CheckBalance is GRpc implemented method to check wallet balance
func (syc *provSysSyncpServer) Credit(
	c context.Context,
	p *syncp.SyncCredit_Req,
) (*syncp.SyncCredit_Res, error) {
	// Get context value
	ctxVal := ctxValue(c)

	// Prepare response
	resp := &syncp.SyncCredit_Res{
		Success: true,
		Code:    int32(contract.OK),
		Data:    nil,
	}
	defer func() {
		resp = nil
	}()

	// Validate payload
	if err := syc.checkBalanceValidate(p); err != nil {
		syc.xloger("CheckBalance", "Validation Error:", err.String())

		resp.Success = false
		resp.Code = int32(err.Code)
		resp.Error = &wrapperspb.StringValue{Value: err.String()}
		return resp, nil
	} else if ctxVal == nil {
		syc.xloger("CheckBalance", "Error: Context value not found")

		resp.Success = false
		resp.Code = int32(contract.CTXVALNOTFOUNDMETADATANULL)
		resp.Error = &wrapperspb.StringValue{Value: contract.CTXVALNOTFOUNDMETADATANULL.String()}
		return resp, nil
	}

	// Get the record
	wallet, rows, e := syc.meta.GetWalletByMember(
		nil,
		p.GetPId(),
		entity.WalletCategory(entity.COMMON), // <- for now set default to be wallet common
		// entity.WalletCategory(ctxVal.cat),
	)
	switch {
	case e != nil:
		syc.xloger("CheckBalance", "Select from DB Error:", e.Error())

		resp.Success = false
		resp.Code = int32(contract.INTERNALERROR)
		resp.Error = &wrapperspb.StringValue{Value: contract.INTERNALERROR.String()}
		return resp, nil

	case rows == 0:
		syc.xloger("CheckBalance", "Record not found")

		resp.Success = false
		resp.Code = int32(contract.WALLETNOTFOUND)
		resp.Error = &wrapperspb.StringValue{Value: contract.WALLETNOTFOUND.String()}
		return resp, nil

	case wallet.IsDisabled:
		syc.xloger("CheckBalance", "Wallet status is disabled")

		resp.Success = false
		resp.Code = int32(contract.WALLETDISABLED)
		resp.Error = &wrapperspb.StringValue{Value: contract.WALLETDISABLED.String()}
		return resp, nil
	}

	wPromo, rows, e := syc.meta.GetWalletPromoByProviderCode(nil, p.GetPId(), nil, true, p.GetProviderCode())
	switch {
	case e != nil:
		syc.xloger("CheckBalancePromo", "Select from DB Error:", e.Error())

		resp.Success = false
		resp.Code = int32(contract.INTERNALERROR)
		resp.Error = &wrapperspb.StringValue{Value: contract.INTERNALERROR.String()}
		return resp, nil

	case rows == 0:
		syc.xloger("CheckBalancePromo", "Record not found")

		wPromo = nil
	}

	credit := wallet.Amount2DecimalPlaces()
	if wPromo != nil {
		credit = wPromo.Amount2DecimalPlacesAll(wallet.Amount)
	}

	// Return the success response
	resp.Data = &syncp.SyncCredit_Data{}
	resp.Data.PId = wallet.PID
	resp.Data.Currency = wallet.Currency
	resp.Data.Credit = credit
	resp.Data.LastUpdate = 0
	if !wallet.UpdatedAt.IsZero() {
		resp.Data.LastUpdate = wallet.UpdatedAt.UnixMilli()
	}
	return resp, nil
}
