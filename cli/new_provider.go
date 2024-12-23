package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/Yoga-Saputra/go-grpc-boilerplate/app/entity"
	"github.com/Yoga-Saputra/go-grpc-boilerplate/boot"
	"github.com/Yoga-Saputra/go-grpc-boilerplate/cli/dbgorm"
	"github.com/Yoga-Saputra/go-grpc-boilerplate/config"
	"github.com/pterm/pterm"
)

// Main variable
var newProvider bool

// Option variable
var (
	newProviderCode      string
	newProviderWalletCat string
)

var newProviderCommands = cli{
	argVar:   &newProvider,
	argName:  "new-provider",
	argUsage: "-new-provider To initialize all needed stuf for new provider",
	run:      newProviderRun,
	stringOptions: []optionString{
		{
			optionVar:          &newProviderCode,
			optionName:         "code",
			optionUsage:        "-code=<provider code> Desired provider code",
			optiondefaultValue: "",
		},
		{
			optionVar:          &newProviderWalletCat,
			optionName:         "cat",
			optionUsage:        "-cat=<wallet category> Desired wallet category",
			optiondefaultValue: "common",
		},
	},
}

func newProviderRun() {
	spinnerLiveText, _ := pterm.DefaultSpinner.Start("Start Preparing Initializer for new provider...")
	time.Sleep(time.Second)

	// Check given flag value
	switch {
	case len(strings.TrimSpace(newProviderCode)) <= 0:
		spinnerLiveText.Fail("Provider Code is required, use -code")
		return

	case len(strings.TrimSpace(newProviderCode)) != 3:
		spinnerLiveText.Fail("Provider Code must be 3 character")
		return
	}

	// Prepare DB
	spinnerLiveText.UpdateText("Prepare DB...")

	// Open DB connection
	spinnerLiveText.UpdateText("Opening DB connection...")
	boot.Up(&boot.AppArgs{}, 1)
	defer func() {
		// Closing DB connection
		spinnerLiveText.UpdateText("Closing DB connection...")
		boot.Down(2)

		spinnerLiveText.Success("DB successfully migrated")
		fmt.Println()
	}()

	// Check DB connection
	if boot.DBA == nil {
		spinnerLiveText.Fail("Failed to open DB connection")
		return
	}

	// Prepare create schema and table
	codeLower := strings.ToLower(newProviderCode)
	codeUpper := strings.ToUpper(newProviderCode)

	// Exec DB migrate
	a := &dbgorm.AdvanceDBMigrate{
		Entity: &entity.TransactionLogProvider{},
		Schema: codeLower,
		Tx:     boot.DBA.DB,
		Source: "stake_log",
	}
	if err := dbgorm.CreateTableTxnLogProvider(a); err != nil {
		spinnerLiveText.Fail("Failed create table ->", err.Error())
	}

	// Create token
	jwtAudPrivKey = config.Of.App.ResolveFilePathInWorkDir("private-key.pem")
	audiderCode = codeUpper
	audiderWalletCat = strings.ToUpper(newProviderWalletCat)
	jwtAudRun()
}
