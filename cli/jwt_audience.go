package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/Yoga-Saputra/go-grpc-boilerplate/pkg/grpcx/util"
	"github.com/Yoga-Saputra/go-grpc-boilerplate/pkg/rsa256"
	"github.com/golang-jwt/jwt"
	"github.com/pterm/pterm"
)

// Main variable argument
var newJWTAud bool

// Option variable argument
var (
	jwtAudPrivKey    string
	audiderCode      string
	audiderWalletCat string
)

var jwtAudSecretKeyLen uint

var jwtAudCommands = cli{
	argVar:   &newJWTAud,
	argName:  "new-jwt-aud",
	argUsage: "-new-jwt-aud To generate new JWT token for each audider",
	run:      jwtAudRun,
	stringOptions: []optionString{
		{
			optionVar:          &jwtAudPrivKey,
			optionName:         "aud-priv-key",
			optionUsage:        "-aud-priv-key=<file path> RSA private key file path",
			optiondefaultValue: "",
		},
		{
			optionVar:          &audiderCode,
			optionName:         "aud-code",
			optionUsage:        "-aud-code=<audider code> Audider code that used as aud JWT",
			optiondefaultValue: "",
		},
		{
			optionVar:          &audiderWalletCat,
			optionName:         "aud-cat",
			optionUsage:        "-aud-cat=<audider wallet category> Audider wallet category",
			optiondefaultValue: "common",
		},
	},
	uintOptions: []optionUInt{
		{
			optionVar:          &jwtAudSecretKeyLen,
			optionName:         "aud-secret-key-len",
			optionUsage:        "-aud-secret-key-len=<int len> JWT generated secret key length",
			optiondefaultValue: 20,
		},
	},
}

func jwtAudRun() {
	spinnerLiveText, _ := pterm.DefaultSpinner.Start("Start Generating new JWT token...")
	time.Sleep(time.Second)

	// Check given flag value
	switch {
	case len(strings.TrimSpace(jwtAudPrivKey)) <= 0:
		spinnerLiveText.Fail("RSA private key file path must be given, use -private-key")
		return

	case len(strings.TrimSpace(audiderCode)) <= 0:
		spinnerLiveText.Fail("Audider code is required, use -aud-code")
		return
	}

	// Read RSA private key
	spinnerLiveText.UpdateText("Reading RSA private key")
	privk, err := rsa256.ReadPrivateKey(jwtAudPrivKey)
	if err != nil {
		spinnerLiveText.Fail(fmt.Sprintf("Failed to read RSA private key: %v", err.Error()))
		return
	}

	// Generate and write the secret key
	spinnerLiveText.UpdateText("Generate and writing the secret key into file")
	aud := strings.ToUpper(audiderCode)
	sig, secretKey, iat := util.GenAuthKey(int(jwtAudSecretKeyLen), aud)

	// Signing RSA private key to JWT
	spinnerLiveText.UpdateText("Signing RSA private key to JWT")
	token := jwt.New(jwt.SigningMethodRS256)

	// Inject the token claims
	walletCat := strings.ToUpper(audiderWalletCat)
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = "SeamlessWallet"
	claims["sub"] = "AudSys"
	claims["cat"] = walletCat
	claims["aud"] = aud
	claims["iat"] = iat
	claims["jti"] = secretKey

	// Signed token
	k, err := token.SignedString(privk)
	if err != nil {
		spinnerLiveText.Fail("Failed to signed JWT with the RSA private key")
		return
	}

	spinnerLiveText.Success("New JWT token has been generated")
	fmt.Println() // Print spacer
	pterm.Info.Println(
		fmt.Sprintf(
			"* Audider Code: %s \n* Wallet Category: %s \n* Signature: %s \n* JWT Token: %s \n* Secret Key: %s",
			strings.ToUpper(audiderCode),
			strings.ToUpper(audiderWalletCat),
			sig,
			k,
			secretKey,
		),
	)
}
