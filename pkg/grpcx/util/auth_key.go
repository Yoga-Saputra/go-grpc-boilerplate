package util

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	char         = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	keySeparator = "<-->"
)

// Generate auth secret key
// the "Raw" of generated key should be
// "{audience}{keySeparator}{signature}{keySeparator}{iat}"
// without "{}" symbol.
func GenAuthKey(keyLength int, aud string) (
	signature, key string,
	iat int64,
) {
	now := time.Now()
	iat = now.Unix()

	// Make random char
	seedrand := rand.New(rand.NewSource(now.UnixNano()))
	b := make([]byte, keyLength)
	for i := range b {
		b[i] = char[seedrand.Intn(len(char))]
	}

	// Make signature
	hasherSig := md5.New()
	hasherSig.Write(b)
	signature = hex.EncodeToString(hasherSig.Sum(nil))

	// Make key
	key = makeKey(aud, signature, fmt.Sprintf("%v", iat))
	return
}

// Verifying secret key based on given JWT claims.
func VerifyKey(sig string, claims jwt.MapClaims) bool {
	jti := claims["jti"]
	aud := claims["aud"]
	iat := claims["iat"]

	var iatI64 int64
	switch iatType := iat.(type) {
	case float64:
		iatI64 = int64(iatType)
	case json.Number:
		v, _ := iatType.Int64()
		iatI64 = v
	}

	key := makeKey(fmt.Sprintf("%v", aud), sig, fmt.Sprintf("%v", iatI64))
	return jti == key
}

func makeKey(aud, sig, iat string) string {
	// Make key
	keyRaw := fmt.Sprintf(
		"%s%s%s%s%v",
		aud,
		keySeparator,
		sig,
		keySeparator,
		iat,
	)
	hasherKey := md5.New()
	hasherKey.Write([]byte(keyRaw))

	return hex.EncodeToString(hasherKey.Sum(nil))
}
