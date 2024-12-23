package cli

import (
	"fmt"
	"time"

	"github.com/Yoga-Saputra/go-grpc-boilerplate/app/contract"
	"github.com/pterm/pterm"
)

// Main variable
var exportStatusCodeContract bool

var exportStatusCodeContractCommands = cli{
	argVar:   &exportStatusCodeContract,
	argName:  "export-status-code",
	argUsage: "-export-status-code To export status code list",
	run:      exportStatusCodeContractRun,
}

func exportStatusCodeContractRun() {
	spinnerLiveText, _ := pterm.DefaultSpinner.Start("Exporting status code contract...")
	time.Sleep(time.Second)
	spinnerLiveText.Success("Finish")

	fmt.Println(contract.StatusCodeToMDTable())
}
