package cli

import (
	"fmt"
	"time"

	"github.com/Yoga-Saputra/go-grpc-boilerplate/pkg/rsa256"
	"github.com/pterm/pterm"
)

// Main variable argument
var newRSAKey bool

var rsaCommands = cli{
	argVar:   &newRSAKey,
	argName:  "new-rsa-pair",
	argUsage: "-new-rsa-pair To generate new RSA private & public key",
	run:      rsaRun,
}

func rsaRun() {
	r := rsa256.New(rsa256.Config{
		PrivateKeyFilePath: "private-key.pem",
		PublicKeyFilePath:  "public-key.pem",
	})

	spinnerLiveText, _ := pterm.DefaultSpinner.Start("Generating RSA pair key...")
	time.Sleep(time.Second)
	if err := r.Generate(); err != nil {
		spinnerLiveText.Fail(fmt.Sprintf("Failed generate RSA pair key: %v", err.Error()))
	} else {
		spinnerLiveText.Success("RSA pair key has been generated")

		fmt.Println() // Print spacer
		pterm.Info.Println(
			fmt.Sprintf("* Private Key: %v \n* Public Key: %v", r.PrivateKeyFilePath, r.PublicKeyFilePath),
		)
	}
}
