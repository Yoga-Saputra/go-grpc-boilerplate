package rsa256

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"os"
)

// Helper function to create RSA private key
func createPrivateKey(f string, bits int) (*rsa.PrivateKey, error) {
	// Generate RSA key
	prvk, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, err
	}

	// Marshal to bytes
	b := x509.MarshalPKCS1PrivateKey(prvk)

	// Create a pem block
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: b,
	}

	// Write to file
	prvpem, err := os.Create(f)
	if err != nil {
		return nil, err
	}
	if err := pem.Encode(prvpem, block); err != nil {
		return nil, err
	}

	return prvk, nil
}

// Helper function to read RSA private key
func readPrivateKey(f string) (*rsa.PrivateKey, error) {
	// Reading file
	r, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}

	// Decode pem block from readed file
	block, _ := pem.Decode(r)

	// Parse block bytes
	prvk, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return prvk, nil
}
