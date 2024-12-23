package rpcwalletinsys

import (
	// "github.com/go-redis/cache/v8"
	"github.com/Yoga-Saputra/go-grpc-boilerplate/pkg/kemu"
	"github.com/Yoga-Saputra/go-grpc-boilerplate/pkg/sw_pb_go/wallet/v1/insys/common"
	"github.com/go-redis/cache/v9"
	"github.com/go-redsync/redsync/v4"

	"github.com/Yoga-Saputra/go-grpc-boilerplate/app/usecase/v1/wallet"
	"google.golang.org/grpc"
)

type (
	// inSysCommonServer defines gRPC service server for "Internal System Wallet Common Task".
	inSysCommonServer struct {
		common.UnimplementedCommonServer
		meta   wallet.Meta
		xloger func(lctx string, m ...string)
		cch    *cache.Cache
	}
)

// RegisterWalletRpc into GRpc Server.
func RegisterWalletRpcInSys(
	s *grpc.Server,
	m wallet.Meta,
	xlg func(lctx string, m ...string),
	cch *cache.Cache,
	kemu *kemu.Mutex,
	rSync *redsync.Redsync,
) func(s *grpc.Server) {
	return func(s *grpc.Server) {

		// Internal system service
		common.RegisterCommonServer(s, &inSysCommonServer{meta: m, xloger: xlg, cch: cch})
	}
}
