package cli

import (
	"flag"
	"os"
	"sort"
)

type cli struct {
	argVar         *bool
	argName        string
	argUsage       string
	run            func()
	cb             func(map[string]cli)
	intOptions     []optionInt
	uintOptions    []optionUInt
	stringOptions  []optionString
	boolOptions    []optionBool
	float64Options []optionFloat64
}

type optionInt struct {
	optionVar          *int
	optionName         string
	optionUsage        string
	optiondefaultValue int
}

type optionUInt struct {
	optionVar          *uint
	optionName         string
	optionUsage        string
	optiondefaultValue uint
}

type optionString struct {
	optionVar          *string
	optionName         string
	optionUsage        string
	optiondefaultValue string
}

type optionBool struct {
	optionVar          *bool
	optionName         string
	optionUsage        string
	optiondefaultValue bool
}

type optionFloat64 struct {
	optionVar          *float64
	optionName         string
	optionUsage        string
	optiondefaultValue float64
}

// Commands defines map of cli function.
// If have new CLIservices, just create new file and their method and register here
var Commands = map[string]cli{
	"app":                appCommands,
	"rsa":                rsaCommands,
	"jwt_rsa":            jwtRSACommands,
	"db_migrate":         dbMigrateCommands,
	"jwt_aud":            jwtAudCommands,
	"new_provider":       newProviderCommands,
	"export_status_code": exportStatusCodeContractCommands,
}

// Execute run the cli function
func Execute() {
	var appc cli

	// Sort map slice
	keys := make([]string, 0, len(Commands))
	for k := range Commands {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		if k == "app" {
			appc = Commands[k]
			continue
		} else if *Commands[k].argVar {
			Commands[k].run()
			os.Exit(0)
		}
	}

	// App info command
	if !*appc.argVar {
		appc.run()
		if appc.cb != nil {
			appc.cb(Commands)
		}
		os.Exit(0)
	}
}

func init() {
	for _, v := range Commands {
		flag.BoolVar(v.argVar, v.argName, false, v.argUsage)

		// Int option
		for _, o := range v.intOptions {
			flag.IntVar(o.optionVar, o.optionName, o.optiondefaultValue, o.optionUsage)
		}

		// UInt option
		for _, o := range v.uintOptions {
			flag.UintVar(o.optionVar, o.optionName, o.optiondefaultValue, o.optionUsage)
		}

		// String option
		for _, o := range v.stringOptions {
			flag.StringVar(o.optionVar, o.optionName, o.optiondefaultValue, o.optionUsage)
		}

		// Bool option
		for _, o := range v.boolOptions {
			flag.BoolVar(o.optionVar, o.optionName, o.optiondefaultValue, o.optionUsage)
		}

		// Float64 option
		for _, o := range v.float64Options {
			flag.Float64Var(o.optionVar, o.optionName, o.optiondefaultValue, o.optionUsage)
		}
	}

	flag.Parse()
}
