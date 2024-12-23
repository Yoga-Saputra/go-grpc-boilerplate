package grpcx

import (
	"errors"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type (
	// MiddlewareFunc defines a function to process middleware.
	MiddlewareFunc interface {
		// Unary provides a hook to intercept the execution of a unary RPC on the server.
		Unary() grpc.UnaryServerInterceptor

		// Stream provides a hook to intercept the execution of a streaming RPC on the server.
		Stream() grpc.StreamServerInterceptor
	}

	// ServiceRegisterFunc gRPC register service function (callback).
	ServiceRegisterFunc func(server *grpc.Server)
)

type (
	// Config defines the config for gRPCx.
	Config struct {
		// Middleware defines a slice middleware function are will be used before create new gRPC server instance.
		Middleware []MiddlewareFunc

		// AddGrpcService defines a gRPC service function are will be used.
		AddGrpcService ServiceRegisterFunc
	}

	// Instance defines gRPC server instance.
	Instance struct {
		Server *grpc.Server

		registerServiceFunc []func(s *grpc.Server)
	}
)

// New create new gRPC server instance.
func NewServer(config ...*Config) *Instance {
	// Init config
	var cfg *Config
	if len(config) > 0 {
		cfg = config[0]
	}

	// Applying the middleware to the server interceptor
	unaries, streams := applyMiddleware(cfg)

	// Create gRPC server
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(unaries...),
		grpc.ChainStreamInterceptor(streams...),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: 5 * time.Minute,
		}),
	)
	i := &Instance{
		Server: server,
	}

	// Create gRPC services instance
	// and register the requierd gRPC method services
	if cfg != nil && cfg.AddGrpcService != nil {
		i.registerServiceFunc = append(i.registerServiceFunc, cfg.AddGrpcService)
	}

	return i
}

// RegisterService into GRpc server.
func (i *Instance) RegisterService(rs ...func(s *grpc.Server)) *Instance {
	i.registerServiceFunc = append(i.registerServiceFunc, rs...)
	return i
}

// Start the gRPC server.
func (i *Instance) Start(nl net.Listener, callback ...func(i *Instance)) error {
	if i != nil {
		for _, v := range i.registerServiceFunc {
			v(i.Server)
		}

		// Run given callback if any
		for _, v := range callback {
			v(i)
		}

		return i.Server.Serve(nl)
	}

	return errors.New("cannot starting gRPC server")
}

// Stop the gRPC server.
func (i *Instance) Stop() {
	if i != nil {
		i.Server.GracefulStop()
	} else {
		log.Printf("[gRPC Server] - gRPC server has been shut down, the instance is null")
	}
}

// applyMiddleware will applying defined middleware into gRPC server interceptor.
func applyMiddleware(cfg *Config) (usi []grpc.UnaryServerInterceptor, ssi []grpc.StreamServerInterceptor) {
	// Append user defined middleware to the interceptor
	if cfg != nil {
		for _, intc := range cfg.Middleware {
			usi = append(usi, intc.Unary())
			ssi = append(ssi, intc.Stream())
		}
	}

	return
}
