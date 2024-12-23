package cli

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/Yoga-Saputra/go-grpc-boilerplate/pkg/rsa256"
	"github.com/golang-jwt/jwt"
	"github.com/pterm/pterm"
)

const jwtRSAChar = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// Main variable argument
var newJWTRSA bool

// Option variable argument
var (
	jwtRSAPrivKey   string
	jwtRSASecretKey string
)

var jwtRSASecretKeyLen uint

var jwtRSACommands = cli{
	argVar:   &newJWTRSA,
	argName:  "new-jwt-rsa",
	argUsage: "-new-jwt-rsa To generate new JWT token with RSA method",
	run:      jwtRSARun,
	stringOptions: []optionString{
		{
			optionVar:          &jwtRSAPrivKey,
			optionName:         "private-key",
			optionUsage:        "-private-key=<file path> RSA private key file path",
			optiondefaultValue: "",
		},
		{
			optionVar:          &jwtRSASecretKey,
			optionName:         "secret-key",
			optionUsage:        "-secret-key=<file path> JWT generated secret key file output",
			optiondefaultValue: "secret.key",
		},
	},
	uintOptions: []optionUInt{
		{
			optionVar:          &jwtRSASecretKeyLen,
			optionName:         "secret-key-len",
			optionUsage:        "-secret-key-len=<int len> JWT generated secret key length",
			optiondefaultValue: 20,
		},
	},
}

func jwtRSARun() {
	spinnerLiveText, _ := pterm.DefaultSpinner.Start("Start Generating new JWT token...")
	time.Sleep(time.Second)

	// Check given flag value
	if jwtRSAPrivKey == "" || jwtRSAPrivKey == " " {
		spinnerLiveText.Fail("RSA private key file path must be given, use --private-key")
		return
	}

	// Read RSA private key
	spinnerLiveText.UpdateText("Reading RSA private key")
	privk, err := rsa256.ReadPrivateKey(jwtRSAPrivKey)
	if err != nil {
		spinnerLiveText.Fail(fmt.Sprintf("Failed to read RSA private key: %v", err.Error()))
		return
	}

	// Generate and write the secret key
	spinnerLiveText.UpdateText("Generate and writing the secret key into file")
	secretKey, err := writeSecretKey()
	if err != nil {
		spinnerLiveText.Fail(fmt.Sprintf("Failed to generate and write secret key: %v", err.Error()))
		return
	}

	// Signing RSA private key to JWT
	spinnerLiveText.UpdateText("Signing RSA private key to JWT")
	token := jwt.New(jwt.SigningMethodRS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sec"] = secretKey
	k, err := token.SignedString(privk)
	if err != nil {
		spinnerLiveText.Fail("Failed to signed JWT with the RSA private key")
		return
	}

	spinnerLiveText.Success("New JWT token has been generated")
	fmt.Println() // Print spacer
	pterm.Info.Println(
		fmt.Sprintf("* Secret Key: %v \n* JWT Token: %v", secretKey, k),
	)
}

// Helper function to generate and write the secret key
func writeSecretKey() (string, error) {
	// Create file
	f, err := os.Create(jwtRSASecretKey)
	if err != nil {
		return "", err
	}
	defer f.Close()

	// Make random char
	seedrand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, jwtRSASecretKeyLen)
	for i := range b {
		b[i] = jwtRSAChar[seedrand.Intn(len(jwtRSAChar))]
	}

	// Convert random char to md5
	hasher := md5.New()
	hasher.Write([]byte(string(b)))
	secretKey := hex.EncodeToString(hasher.Sum(nil))

	// Write to file
	if _, err := f.WriteString(secretKey); err != nil {
		return "", err
	}

	return secretKey, nil
}
