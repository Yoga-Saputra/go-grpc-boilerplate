package rpcwalletprovsys

import (
	"context"
	"fmt"

	"github.com/Yoga-Saputra/go-grpc-boilerplate/pkg/grpcx"
	"github.com/golang-jwt/jwt"
)

// Struct to hold value of parsed auth meta from context.
type authMetaContext struct {
	aud string
	jti string
	cat string
}

// Return structured value of parsed auth meta from context.
func ctxValue(ctx context.Context) *authMetaContext {
	// Get value fropm context by key from grpcx
	val := ctx.Value(grpcx.ModCtxKey)
	if val == nil {
		return nil
	}

	// Assert interface value to jwt.MapClaims
	v, ok := val.(jwt.MapClaims)
	if !ok {
		return nil
	}

	// Get set wallet category
	walletCat := "COMMON"
	if cat, ok := v["cat"]; ok {
		walletCat = fmt.Sprintf("%v", cat)
	}

	return &authMetaContext{
		aud: fmt.Sprintf("%v", v["aud"]),
		jti: fmt.Sprintf("%v", v["jti"]),
		cat: walletCat,
	}
}
