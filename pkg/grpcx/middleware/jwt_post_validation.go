package middleware

import (
	"context"
	"strings"

	"github.com/Yoga-Saputra/go-grpc-boilerplate/pkg/grpcx"
	"github.com/Yoga-Saputra/go-grpc-boilerplate/pkg/grpcx/util"
	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type (
	// JWTPostValidationConfig defines the config of JWT Post Validation middleware
	JWTPostValidationConfig struct {
		// TokenLookup is a string in the form of "<name>" or "<name>,<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Signature".
		// Possible values:
		// - "Signature"
		// - "Sig"
		// - "X-Sig"
		// Multiply sources example:
		// - "header: Signature,X-Sig"
		TokenLookup string

		// If this field set to false, validation will be ignored
		// if not have any context claim.
		// Default value is false.
		Required bool

		// OnlyForMethod to apply to certain gRPC service method.
		// Validation will using "contains" instead of matching exact.
		OnlyForMethod []string
	}
)

var (
	// DefaultJWTConfig is the default JWT auth middleware config.
	DefaultJWTPostValidationConfig = JWTPostValidationConfig{
		Required:    false,
		TokenLookup: "Signature",
	}
)

// JWTPostValidation return JWT Post Validation.
func JWTPostValidation() *JWTPostValidationConfig {
	c := DefaultJWTPostValidationConfig
	return JWTPostValidationWithConfig(c)
}

// JWTPostValidationWithConfig return JWT Post Validation middleware with config.
// See: `JWTPostValidation()`.
func JWTPostValidationWithConfig(config JWTPostValidationConfig) *JWTPostValidationConfig {
	// Set defaults
	if len(strings.TrimSpace(config.TokenLookup)) <= 0 {
		config.TokenLookup = DefaultJWTPostValidationConfig.TokenLookup
	}

	return &config
}

// ***************** Interceptor Implement *****************

// Unary provides a hook to intercept the execution of a unary RPC on the server.
func (cfg *JWTPostValidationConfig) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		// Do authenticating
		if err := cfg.validation(ctx, info.FullMethod); err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

// Stream provides a hook to intercept the execution of a streaming RPC on the server.
func (cfg *JWTPostValidationConfig) Stream() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		return cfg.validation(ss.Context(), info.FullMethod)
	}
}

// Do post validate on interceptor.
func (cfg *JWTPostValidationConfig) validation(c context.Context, method string) error {
	// Metadata extractor
	sources := strings.Split(cfg.TokenLookup, ",")
	var extractors []jwtExtractor
	for _, source := range sources {
		extractors = append(extractors, cfg.signatureFromMetaData(source))
	}

	// Getting signature
	var sig string
	var err error
	for _, extractor := range extractors {
		sig, err = extractor(c, method)
		if err != nil {
			break
		}
	}
	if err != nil {
		return err
	}

	// Get custom value from context
	user := c.Value(grpcx.ModCtxKey)
	if user == nil {
		if cfg.Required {
			return status.Error(codes.Unauthenticated, "failed to validate, user context if not found")
		}
	}

	// Parse token claims
	claims, ok := user.(jwt.MapClaims)
	if !ok {
		if cfg.Required {
			return status.Error(codes.Unauthenticated, "failed to validate, failed to asert token claims")
		}
	}

	// Do validation
	if len(cfg.OnlyForMethod) > 0 {
		for _, v := range cfg.OnlyForMethod {
			if strings.Contains(strings.ToLower(method), strings.ToLower(v)) {
				if valid := util.VerifyKey(sig, claims); !valid {
					return status.Error(codes.Unauthenticated, "your token seems to be invalid")
				}
			}
		}

		return nil
	} else {
		if valid := util.VerifyKey(sig, claims); !valid {
			return status.Error(codes.Unauthenticated, "your token seems to be invalid")
		}

		return nil
	}
}

// jwtFromMetaData returns a `jwtExtractor` that extracts token from the request meta data.
func (cfg *JWTPostValidationConfig) signatureFromMetaData(tokenLookup string) jwtExtractor {
	return func(c context.Context, method string) (string, error) {
		// Metadata extractor & checker
		meta, ok := metadata.FromIncomingContext(c)
		if !ok {
			return "", status.Errorf(codes.Unauthenticated, "metadata is not provided")
		}

		// Metadata found
		vals := meta[strings.ToLower(tokenLookup)]
		for _, sigs := range vals {
			return sigs, nil
		}

		if cfg.Required && len(cfg.OnlyForMethod) > 0 {
			for _, v := range cfg.OnlyForMethod {
				if strings.Contains(strings.ToLower(method), strings.ToLower(v)) {
					return "", status.Errorf(codes.Unauthenticated, "missing Signature")
				}
			}

			return "", nil
		} else {
			return "", nil
		}
	}
}
