package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/Yoga-Saputra/go-grpc-boilerplate/boot"
	"github.com/Yoga-Saputra/go-grpc-boilerplate/cli/dbgorm"
	"github.com/pterm/pterm"
)

// Main variable argument
var dbMigrate bool

// Option variable argument
var tableName string

var seedRefresh bool

var dbMigrateCommands = cli{
	argVar:   &dbMigrate,
	argName:  "db-migrate",
	argUsage: "-db-migrate To start migrations. If without sub-argument will migrate all table",
	run:      dbMigrateRun,
	stringOptions: []optionString{
		{
			optionVar:          &tableName,
			optionName:         "table",
			optionUsage:        "-table=<table name> Just migrate specific table instead migrate all",
			optiondefaultValue: "",
		},
	},
	boolOptions: []optionBool{
		{
			optionVar:          &seedRefresh,
			optionName:         "seed-refresh",
			optionUsage:        "-seed-refresh seed with refresh",
			optiondefaultValue: false,
		},
	},
}

func dbMigrateRun() {
	spinnerLiveText, _ := pterm.DefaultSpinner.Start("Doing DB migrations...")
	time.Sleep(time.Second)

	// Open DB connection
	spinnerLiveText.UpdateText("Opening DB connection...")
	boot.Up(&boot.AppArgs{}, 1)
	defer func() {
		// Closing DB connection
		spinnerLiveText.UpdateText("Closing DB connection...")
		boot.Down(1)

		spinnerLiveText.Success("DB successfully migrated")
		fmt.Println()
	}()

	// Check DB connection
	if boot.DBA == nil {
		spinnerLiveText.Fail("Failed to open DB connection")
		return
	}

	if tableName != "" && tableName != " " {
		// Start migration based on given table name if any
		spinnerLiveText.UpdateText(fmt.Sprintf("Just migrate %s table...", tableName))
		t := strings.Split(tableName, ",")
		if err := startMigrator(spinnerLiveText, t...); err != nil {
			spinnerLiveText.Fail(err.Error())
		}

		return
	}

	// Start migration all
	if err := startMigrator(spinnerLiveText); err != nil {
		spinnerLiveText.Fail(err.Error())
	}
}

// Helper function to execute migrator up
func startMigrator(st *pterm.SpinnerPrinter, t ...string) error {
	st.UpdateText("Start migration...")

	return dbgorm.RunDBMigrate(boot.DBA.DB, seedRefresh, t...)
}
