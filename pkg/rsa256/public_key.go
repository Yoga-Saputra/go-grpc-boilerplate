package rsa256

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"os"
)

// Helper function to create RSA public key
func createPublicKey(f string, prvk *rsa.PrivateKey) error {
	// Get public key from private key
	pubk := &prvk.PublicKey

	// Marshal to bytes
	b, err := x509.MarshalPKIXPublicKey(pubk)
	if err != nil {
		return err
	}

	// Create a pem block
	block := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: b,
	}

	// Write to file
	pubpem, err := os.Create(f)
	if err != nil {
		return err
	}
	if err := pem.Encode(pubpem, block); err != nil {
		return err
	}

	return nil
}

// Helper cuntion to read RSA public key
func readPublicKey(f string) (interface{}, error) {
	// Reading file
	r, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}

	// Decode pem block from readed file
	block, _ := pem.Decode(r)

	// Parse block bytes
	pubk, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return pubk, nil
}
