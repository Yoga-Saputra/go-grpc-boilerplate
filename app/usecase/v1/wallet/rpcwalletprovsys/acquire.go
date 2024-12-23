package rpcwalletprovsys

import (
	// "github.com/go-redis/cache/v8"
	"github.com/Yoga-Saputra/go-grpc-boilerplate/pkg/kemu"
	"github.com/Yoga-Saputra/go-grpc-boilerplate/pkg/sw_pb_go/wallet/v1/provsys/syncp"
	"github.com/go-redis/cache/v9"
	"github.com/go-redsync/redsync/v4"

	"github.com/Yoga-Saputra/go-grpc-boilerplate/app/usecase/v1/wallet"
	"google.golang.org/grpc"
)

type (
	// provSysSyncpServer defines gRPC service server for "Provider System Sync" service.
	provSysSyncpServer struct {
		syncp.UnimplementedSyncpServer
		meta   wallet.Meta
		xloger func(lctx string, m ...string)
	}
)

// RegisterWalletRpcProvSys into GRpc Server.
func RegisterWalletRpcProvSys(
	s *grpc.Server,
	m wallet.Meta,
	xlg func(lctx string, m ...string),
	cch *cache.Cache,
	kemu *kemu.Mutex,
	rSync *redsync.Redsync,
) func(s *grpc.Server) {
	return func(s *grpc.Server) {

		// Provider system service
		syncp.RegisterSyncpServer(s, &provSysSyncpServer{meta: m, xloger: xlg})
	}
}
