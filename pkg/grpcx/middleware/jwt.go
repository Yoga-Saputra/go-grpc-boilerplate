package middleware

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/Yoga-Saputra/go-grpc-boilerplate/pkg/grpcx"
	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type (
	// JWTConfig defines the config for JWT middleware.
	JWTConfig struct {
		// Signing key to validate token.
		// This is one of the three options to provide a token validation key.
		// The order of precedence is a user-defined KeyFunc, SigningKeys and SigningKey.
		// Required if neither user-defined KeyFunc nor SigningKeys is provided.
		SigningKey interface{}

		// Map of signing keys to validate token with kid field usage.
		// This is one of the three options to provide a token validation key.
		// The order of precedence is a user-defined KeyFunc, SigningKeys and SigningKey.
		// Required if neither user-defined KeyFunc nor SigningKey is provided.
		SigningKeys map[string]interface{}

		// Signing method used to check the token's signing algorithm.
		// Optional. Default value HS256.
		SigningMethod string

		// Claims are extendable claims data defining token content. Used by default ParseTokenFunc implementation.
		// Not used if custom ParseTokenFunc is set.
		// Optional. Default value jwt.MapClaims
		Claims jwt.Claims

		// TokenLookup is a string in the form of "<name>" or "<name>,<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "Authorization"
		// - "Auth"
		// - "X-Auth"
		// Multiply sources example:
		// - "header: Authorization,X-Auth"
		TokenLookup string

		// AuthScheme to be used in the Authorization header.
		// Optional. Default value "Bearer".
		AuthScheme string

		// KeyFunc defines a user-defined function that supplies the public key for a token validation.
		// The function shall take care of verifying the signing algorithm and selecting the proper key.
		// A user-defined KeyFunc can be useful if tokens are issued by an external party.
		// Used by default ParseTokenFunc implementation.
		//
		// When a user-defined KeyFunc is provided, SigningKey, SigningKeys, and SigningMethod are ignored.
		// This is one of the three options to provide a token validation key.
		// The order of precedence is a user-defined KeyFunc, SigningKeys and SigningKey.
		// Required if neither SigningKeys nor SigningKey is provided.
		// Not used if custom ParseTokenFunc is set.
		// Default to an internal implementation verifying the signing algorithm and selecting the proper key.
		KeyFunc jwt.Keyfunc

		// ParseTokenFunc defines a user-defined function that parses token from given auth. Returns an error when token
		// parsing fails or parsed token is invalid.
		// Defaults to implementation using `github.com/golang-jwt/jwt` as JWT implementation library
		ParseTokenFunc func(auth string) (interface{}, error)

		// IgnoreMethod to ignoring given methods to be authenticating.
		// This filed will be ignored if `ApplyOnlyOnMethod` setted.
		IgnoreMethod []string

		// ApplyOnlyOnMethod to apply the authenticating process only on given methods.
		ApplyOnlyOnMethod []string
	}

	jwtExtractor func(ctx context.Context, fullMethod string) (string, error)
)

const (
	// Algorithms
	AlgorithmHS256 = "HS256"
)

var (
	// DefaultJWTConfig is the default JWT auth middleware config.
	DefaultJWTConfig = JWTConfig{
		SigningMethod: AlgorithmHS256,
		TokenLookup:   "Authorization",
		AuthScheme:    "Bearer",
		Claims:        jwt.MapClaims{},
		KeyFunc:       nil,
	}
)

// JWT returns a JSON Web Token (JWT) auth middleware.
//
// For valid token, it sets the user in context and calls next handler.
// For invalid token, it returns "401 - Unauthorized" error.
// For missing token, it returns "400 - Bad Request" error.
//
// See: https://jwt.io/introduction
// See `JWTConfig.TokenLookup`
func JWT(key interface{}) *JWTConfig {
	c := DefaultJWTConfig
	c.SigningKey = key
	return JWTWithConfig(c)
}

// JWTWithConfig returns a JWT auth middleware with config.
// See: `JWT()`.
func JWTWithConfig(config JWTConfig) *JWTConfig {
	if config.SigningKey == nil && len(config.SigningKeys) == 0 && config.KeyFunc == nil && config.ParseTokenFunc == nil {
		panic("gRPCx: jwt middleware requires signing key")
	}

	// Set defaults
	if config.SigningMethod == "" || config.SigningMethod == " " {
		config.SigningMethod = DefaultJWTConfig.SigningMethod
	}
	if config.Claims == nil {
		config.Claims = DefaultJWTConfig.Claims
	}
	if config.TokenLookup == "" || config.TokenLookup == " " {
		config.TokenLookup = DefaultJWTConfig.TokenLookup
	}
	if config.AuthScheme == "" || config.AuthScheme == " " {
		config.AuthScheme = DefaultJWTConfig.AuthScheme
	}
	if config.KeyFunc == nil {
		config.KeyFunc = config.defaultKeyFunc
	}
	if config.ParseTokenFunc == nil {
		config.ParseTokenFunc = config.defaultParseToken
	}

	return &config
}

// ***************** Interceptor Implement *****************

// Unary provides a hook to intercept the execution of a unary RPC on the server.
func (cfg *JWTConfig) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		// Do authenticating
		modCtx, err := cfg.authorize(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}

		return handler(modCtx, req)
	}
}

// Stream provides a hook to intercept the execution of a streaming RPC on the server.
func (cfg *JWTConfig) Stream() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		// Do authenticating
		modCtx, err := cfg.authorize(ss.Context(), info.FullMethod)
		if err != nil {
			return err
		}

		wrapped := grpcx.WrapServerStream(ss)
		wrapped.WrappedContext = modCtx
		return handler(srv, wrapped)
	}
}

// authorize will validating JWT token
func (cfg *JWTConfig) authorize(c context.Context, method string) (context.Context, error) {
	// Ignoring given method to be authenticating.
	// Or just apply authenticating process to the given methods only.
	if len(cfg.ApplyOnlyOnMethod) > 0 {
		for _, v := range cfg.ApplyOnlyOnMethod {
			if method != v {
				return c, nil
			}
		}
	} else if len(cfg.IgnoreMethod) > 0 {
		for _, v := range cfg.IgnoreMethod {
			if method == v {
				return c, nil
			}
		}
	}

	// Initialize
	// Split sources
	sources := strings.Split(cfg.TokenLookup, ",")
	var extractors []jwtExtractor
	for _, source := range sources {
		extractors = append(extractors, cfg.jwtFromMetaData(source, cfg.AuthScheme))
	}

	// Getting auth token
	var auth string
	var err error
	for _, extractor := range extractors {
		auth, err = extractor(c, method)
		if err != nil {
			break
		}
	}
	if err != nil {
		return nil, err
	}

	// Parse the token
	token, err := cfg.ParseTokenFunc(auth)
	if err != nil {
		return nil, err
	}

	// Assert JWT token
	t, tok := token.(*jwt.Token)
	if !tok {
		return nil, errors.New("failed assert token")
	}

	// Assert JWt map claims
	claims, cok := t.Claims.(jwt.MapClaims)
	if !cok {
		return nil, errors.New("failed assert token claims")
	}

	// Modify context on the go
	modCtx := context.WithValue(c, grpcx.ModCtxKey, claims)
	return modCtx, nil
}

// defaultParseToken return a JWT parsed token
func (config *JWTConfig) defaultParseToken(auth string) (interface{}, error) {
	var token *jwt.Token
	var err error

	// Issue #647, #656
	if _, ok := config.Claims.(jwt.MapClaims); ok {
		token, err = jwt.Parse(auth, config.KeyFunc)
	} else {
		t := reflect.ValueOf(config.Claims).Type().Elem()
		claims := reflect.New(t).Interface().(jwt.Claims)
		token, err = jwt.ParseWithClaims(auth, claims, config.KeyFunc)
	}

	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}
	if !token.Valid {
		return nil, status.Errorf(codes.Unauthenticated, "invalid JWT token")
	}
	return token, nil
}

// defaultKeyFunc returns a signing key of the given token.
func (config *JWTConfig) defaultKeyFunc(t *jwt.Token) (interface{}, error) {
	// Check the signing method
	if t.Method.Alg() != config.SigningMethod {
		return nil, fmt.Errorf("unexpected jwt signing method=%v", t.Header["alg"])
	}
	if len(config.SigningKeys) > 0 {
		if kid, ok := t.Header["kid"].(string); ok {
			if key, ok := config.SigningKeys[kid]; ok {
				return key, nil
			}
		}
		return nil, fmt.Errorf("unexpected jwt key id=%v", t.Header["kid"])
	}

	return config.SigningKey, nil
}

// jwtFromMetaData returns a `jwtExtractor` that extracts token from the request meta data.
func (cfg *JWTConfig) jwtFromMetaData(tokenLookup string, authScheme string) jwtExtractor {
	return func(c context.Context, method string) (string, error) {
		// Metadata extractor & checker
		meta, ok := metadata.FromIncomingContext(c)
		if !ok {
			return "", status.Errorf(codes.Unauthenticated, "metadata is not provided")
		}

		// Metadata found
		l := len(authScheme)
		vals := meta[strings.ToLower(tokenLookup)]
		for _, auth := range vals {
			if len(auth) > l+1 && auth[:l] == authScheme {
				return auth[l+1:], nil
			}
		}

		return "", status.Errorf(codes.Unauthenticated, "missing JWT token")
	}
}
