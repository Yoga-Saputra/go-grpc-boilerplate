package middleware

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type (
	softMaintenance struct{}
)

func SoftMaintenance() *softMaintenance {
	return &softMaintenance{}
}

// ***************** Interceptor Implement *****************

// Unary provides a hook to intercept the execution of a unary RPC on the server.
func (cfg *softMaintenance) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		return nil, status.Errorf(codes.Aborted, "gRPC service is under maintenance")
	}
}

// Stream provides a hook to intercept the execution of a streaming RPC on the server.
func (cfg *softMaintenance) Stream() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		return status.Errorf(codes.Aborted, "gRPC service is under maintenance")
	}
}
