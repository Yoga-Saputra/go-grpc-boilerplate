package rsa256

// Config defines the config for RSA Generator
type Config struct {
	// Full file path of private key that want to be generate
	//
	// Required. Default  "private-key.pem"
	PrivateKeyFilePath string

	// Full file path of public key that want to be generate
	//
	// Required. Default  "public-key.pem"
	PublicKeyFilePath string

	// Size of bits taht will be using to generate RSA key pair
	//
	// Optional. Default 2048
	BitSize int
}

// ConfigDefault is the defualt config
var ConfigDefault = Config{
	PrivateKeyFilePath: "private-key.pem",
	PublicKeyFilePath:  "public-key.pem",
	BitSize:            2048,
}

// Helper function to set default config
func configDefault(config ...Config) Config {
	// Return the default config
	if len(config) <= 0 {
		return ConfigDefault
	}

	// Overide the defualt config
	cfg := config[0]

	if cfg.PrivateKeyFilePath == "" || cfg.PrivateKeyFilePath == " " {
		cfg.PrivateKeyFilePath = ConfigDefault.PrivateKeyFilePath
	}
	if cfg.PublicKeyFilePath == "" || cfg.PublicKeyFilePath == " " {
		cfg.PublicKeyFilePath = ConfigDefault.PublicKeyFilePath
	}
	if cfg.BitSize <= 0 {
		cfg.BitSize = ConfigDefault.BitSize
	}

	return cfg
}
