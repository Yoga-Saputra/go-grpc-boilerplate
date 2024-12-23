package service

import (
	"log"
	"strings"
)

// GrpcxLogger logging
func GrpcxLogger(lctx string, m ...string) {
	svc := "[gRPC]"
	log.Printf("%s[%s] - %s", svc, lctx, strings.Join(m, " "))
}
