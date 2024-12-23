package boot

import (
	"fmt"
	"os"
	"time"

	// "github.com/go-redis/cache/v8"
	"github.com/Yoga-Saputra/go-grpc-boilerplate/app/repo"
	"github.com/Yoga-Saputra/go-grpc-boilerplate/app/service"
	"github.com/Yoga-Saputra/go-grpc-boilerplate/app/usecase/v1/mcslog"
	"github.com/Yoga-Saputra/go-grpc-boilerplate/app/usecase/v1/txnlog"
	"github.com/Yoga-Saputra/go-grpc-boilerplate/app/usecase/v1/txnlogprovider"
	"github.com/Yoga-Saputra/go-grpc-boilerplate/app/usecase/v1/wallet"
	"github.com/Yoga-Saputra/go-grpc-boilerplate/app/usecase/v1/wallet/rpcwalletinsys"
	"github.com/Yoga-Saputra/go-grpc-boilerplate/app/usecase/v1/wallet/rpcwalletprovsys"
	"github.com/Yoga-Saputra/go-grpc-boilerplate/config"
	"github.com/Yoga-Saputra/go-grpc-boilerplate/pkg/grpcx"
	"github.com/Yoga-Saputra/go-grpc-boilerplate/pkg/grpcx/middleware"
	"github.com/Yoga-Saputra/go-grpc-boilerplate/pkg/kemu"
	"github.com/go-redis/cache/v9"
	"github.com/go-redsync/redsync/v4"
	"google.golang.org/grpc"
)

// RPC singleton of gRPC server instance.
var RPC *grpcx.Instance

// Start gRPC server.
func rpcUp(args *AppArgs) {
	if HardMaintenance == "false" {
		// Determine config based on service mode
		// is it set to force maintenance or not
		var gc *grpcx.Config
		if SoftMaintenance == "false" {
			gc = &grpcx.Config{
				// gRPC server middleware
				Middleware: []grpcx.MiddlewareFunc{
					middleware.JWTWithConfig(middleware.JWTConfig{
						SigningMethod: "RS256",
						SigningKey:    config.Of.App.GetPublicKey(),
					}),
					middleware.JWTPostValidationWithConfig(middleware.JWTPostValidationConfig{
						Required:      true,
						OnlyForMethod: []string{"provsys"},
					}),
				},
			}
		} else {
			gc = &grpcx.Config{
				Middleware: []grpcx.MiddlewareFunc{middleware.SoftMaintenance()},
			}
		}

		// Create new instance
		rpc := grpcx.NewServer(gc)

		// Registry gRPC server services
		rpc = rpc.RegisterService(acquireUsecaseV1(rpc, Cache.RCch, Cache.Rs)...)

		// Hold grpcx instance pointer to local variable
		RPC = rpc

		// Start gRPC server using goroutine
		go func() {
			if err := rpc.Start(args.NL, func(i *grpcx.Instance) {
				// Print info
				printOutUp(fmt.Sprintf("⇨ gRPC services running on: %v  pId: %v", args.NL.Addr().String(), os.Getpid()))
				for k, v := range i.Server.GetServiceInfo() {
					printOutUp(fmt.Sprintf("⇨ gRPC %s -> %v", k, v.Methods))
				}
			}); err != nil {
				panic(err)
			}
		}()
	}
}

// Stop gRPC server.
func rpcDown() {
	if RPC != nil {
		printOutDown("Shutting down gRPC services...")

		RPC.Stop()
	}
}

// acquiring usecase v1 into the GRpc service handler.
func acquireUsecaseV1(i *grpcx.Instance, cch *cache.Cache, rSync *redsync.Redsync) (rs []func(s *grpc.Server)) {
	printOutUp("Acquire usecase V1...")

	// Load time zone
	tz, err := time.LoadLocation(config.Of.App.TimeZone)
	if err != nil {
		panic(err)
	}

	// Acquire implemented gRPC interface
	rs = acquireUsecaseV1Rpc(i, tz, cch, rSync)

	// Registering transaction log usecase v1.
	txnlog.RegisterUsecase(
		repo.NewTxnLogRepoDB(DBA.DB),
		printOutUp,
	)

	// Registering transaction log provider usecase v1.
	txnlogprovider.RegisterUsecase(
		repo.NewTxnProviderLogRepoDB(DBA.DB, tz),
		printOutUp,
	)

	// Register mcs log user case v1
	mcslogRepo, err := repo.NewTransferLogRepoDB(DBA.DB)
	if err != nil {
		printOutUp(fmt.Sprintf("Failed acquire mcslog repo: %s", err.Error()))
	} else {
		mcslog.RegisterUseCase(mcslogRepo, tz, printOutUp)
	}

	return
}

func acquireUsecaseV1Rpc(
	i *grpcx.Instance,
	tz *time.Location,
	cch *cache.Cache,
	rSync *redsync.Redsync,
) (rs []func(s *grpc.Server)) {
	// define global kemu
	kemu := kemu.New()
	// Register Usecase & gRPC wallet V1
	walletMeta := wallet.RegisterUsecase(
		repo.NewWalletRepoDB(DBA.DB),
		tz,
		printOutUp,
	)
	walIntSysRpcServer := rpcwalletinsys.RegisterWalletRpcInSys(
		i.Server,
		*walletMeta,
		service.GrpcxLogger,
		cch,
		kemu,
		rSync,
	) // <- gRPC service for "Internal System" wallet
	walProvSysRpcServer := rpcwalletprovsys.RegisterWalletRpcProvSys(
		i.Server,
		*walletMeta,
		service.GrpcxLogger,
		cch,
		kemu,
		rSync,
	) // <- gRPC service for "Provider System" wallet
	rs = append(rs, walIntSysRpcServer, walProvSysRpcServer)

	return
}
