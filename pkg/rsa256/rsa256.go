package rsa256

import (
	"crypto/rsa"
	"os"
)

// New create new config pointer of RSA256 package
func New(config ...Config) *Config {
	// Set default config
	cfg := configDefault(config...)

	return &cfg
}

// Generate generates new RSA key pair and write to given file path on config
func (c *Config) Generate() error {
	// RSA private key
	prvk, err := createPrivateKey(c.PrivateKeyFilePath, c.BitSize)
	if err != nil {
		return err
	}

	// RSA public key
	if err := createPublicKey(c.PublicKeyFilePath, prvk); err != nil {
		return err
	}

	return nil
}

// ReadPrivateKey reads RSA private key from given file path
func ReadPrivateKey(f string) (*rsa.PrivateKey, error) {
	return readPrivateKey(f)
}

// ReadPublicKey reads RSA public key from given file path
func ReadPublicKey(f string) (interface{}, error) {
	return readPublicKey(f)
}

// CheckRSAFileExists return true if RSA private and public key is exists or readable,
// otherwise wull return false
func CheckRSAFileExists(private, public string) bool {
	var (
		privateExists = true
		publicExists  = true
	)

	// RSA private key check
	if _, err := os.Stat(private); err != nil {
		privateExists = false
	}

	// RSA public key check
	if _, err := os.Stat(public); err != nil {
		publicExists = false
	}

	return privateExists && publicExists
}
