package grpcx

import (
	"context"

	"google.golang.org/grpc"
)

type (
	modCtxKey string

	// WrappedServerStream is a wrapper for grpc.ServerStream that makes context can be modifying.
	WrappedServerStream struct {
		grpc.ServerStream

		// WrappedContext is the modified context.
		WrappedContext context.Context
	}
)

const ModCtxKey modCtxKey = "user"

// Context return modified context of ServerStream.
func (w *WrappedServerStream) Context() context.Context {
	return w.WrappedContext
}

// WrapServerStream return ServerStream that has ability to be overwrited context.
func WrapServerStream(s grpc.ServerStream) *WrappedServerStream {
	if existing, ok := s.(*WrappedServerStream); ok {
		return existing
	}

	return &WrappedServerStream{
		ServerStream:   s,
		WrappedContext: s.Context(),
	}
}
