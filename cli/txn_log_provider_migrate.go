package cli

// import (
// 	"fmt"
// 	"os"
// 	"strings"
// 	"time"

// 	"github.com/Yoga-Saputra/go-grpc-boilerplate/app/entity"
// 	"github.com/Yoga-Saputra/go-grpc-boilerplate/boot"
// 	"github.com/pterm/pterm"
// )

// // Main variable
// var txnLogProvider bool

// // Option variable
// var (
// 	txnLogProviderCode  string
// 	txnLogProviderMonth string
// )

// var txnLogProviderCommands = cli{
// 	argVar:   &txnLogProvider,
// 	argName:  "txn-log-provider",
// 	argUsage: "-txn-log-provider To migrating transaction log provider table",
// 	run:      txnLogProviderRun,
// 	stringOptions: []optionString{
// 		{
// 			optionVar:          &txnLogProviderCode,
// 			optionName:         "log-prov-code",
// 			optionUsage:        "-log-prov-code=<provider code> Desired provider code",
// 			optiondefaultValue: "",
// 		},
// 		{
// 			optionVar:          &txnLogProviderMonth,
// 			optionName:         "log-prov-month",
// 			optionUsage:        "-log-prov-month=<month name 3 char> Desired month name",
// 			optiondefaultValue: "",
// 		},
// 	},
// }

// func txnLogProviderRun() {
// 	spinnerLiveText, _ := pterm.DefaultSpinner.Start("Start Preparing Migration TxnLogProvider...")
// 	time.Sleep(time.Second)

// 	// Check given flag value
// 	// if len(strings.TrimSpace(txnLogProviderCode)) <= 0 {
// 	// 	spinnerLiveText.Fail("Provider Code is required, use -log-prov-code")
// 	// 	return
// 	// } else if len(strings.TrimSpace(txnLogProviderCode)) != 3 {
// 	// 	spinnerLiveText.Fail("Provider Code must be 3 character")
// 	// 	return
// 	// }

// 	// Prepare DB
// 	spinnerLiveText.UpdateText("Prepare DB...")

// 	// Open DB connection
// 	spinnerLiveText.UpdateText("Opening DB connection...")
// 	boot.Up(&boot.AppArgs{}, 1)
// 	defer func() {
// 		// Closing DB connection
// 		spinnerLiveText.UpdateText("Closing DB connection...")
// 		boot.Down(1)

// 		spinnerLiveText.Success("DB successfully migrated")
// 		fmt.Println()
// 	}()

// 	// Check DB connection
// 	if boot.DBA == nil {
// 		spinnerLiveText.Fail("Failed to open DB connection")
// 		return
// 	}

// 	// Prepare create schema and table
// 	codeLower := strings.ToLower(txnLogProviderCode)
// 	tx := boot.DBA

// 	// Get DB dialect
// 	dialect := tx.DB.Dialector.Name()
// 	ent := &entity.TransactionLogProvider{}

// 	// Do create DB stuff based on dialect
// 	spinnerLiveText.Success("Using dialect: ", dialect)
// 	switch dialect {
// 	case "postgres":
// 		type schemaName struct {
// 			SchemaName string `json:"schema_name"`
// 		}
// 		var schemas []schemaName

// 		if len(strings.TrimSpace(codeLower)) > 0 {
// 			// Create schema if not exists
// 			spinnerLiveText.Success("Create db schema: ", codeLower)
// 			tx.DB.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", codeLower))
// 		} else {
// 			tx.DB.Raw(`
// 			SELECT
// 				schema_name
// 			FROM
// 				information_schema.schemata
// 			WHERE
// 				"schema_name" NOT IN('information_schema', 'public', 'pg_catalog', 'pg_toast')
// 			`).Scan(&schemas)
// 		}

// 		// Create as many table as month in a year
// 		if len(strings.TrimSpace(txnLogProviderMonth)) <= 0 {
// 			if len(schemas) > 0 {
// 				for _, s := range schemas {
// 					for _, m := range monthList {
// 						spinnerLiveText.Success("Create table -> ", m)
// 						tx.DB.Scopes(entity.ProvDynamicTable(ent, m, s.SchemaName)).AutoMigrate(&entity.TransactionLogProvider{})
// 					}
// 				}
// 			} else {
// 				for _, m := range monthList {
// 					spinnerLiveText.Success("Create table -> ", m)
// 					tx.DB.Scopes(entity.ProvDynamicTable(ent, m, codeLower)).AutoMigrate(&entity.TransactionLogProvider{})
// 				}
// 			}
// 		} else {
// 			if len(schemas) > 0 {
// 				for _, s := range schemas {
// 					monthLower := strings.ToLower(txnLogProviderMonth)
// 					spinnerLiveText.Success("Create table -> ", monthLower)
// 					tx.DB.Scopes(entity.ProvDynamicTable(ent, monthLower, s.SchemaName)).AutoMigrate(&entity.TransactionLogProvider{})
// 				}
// 			} else {
// 				monthLower := strings.ToLower(txnLogProviderMonth)
// 				spinnerLiveText.Success("Create table -> ", monthLower)
// 				tx.DB.Scopes(entity.ProvDynamicTable(ent, monthLower, codeLower)).AutoMigrate(&entity.TransactionLogProvider{})
// 			}
// 		}

// 		spinnerLiveText.Success("DB stuff finish")
// 	default:
// 		spinnerLiveText.Fail("DB Dialect not supported yet")
// 		os.Exit(0)
// 	}
// }
