package cli

import (
	"fmt"
	"time"

	"github.com/Yoga-Saputra/go-grpc-boilerplate/boot"
	"github.com/Yoga-Saputra/go-grpc-boilerplate/config"
	"github.com/pterm/pterm"
)

var runApp bool

var appCommands = cli{
	argVar:   &runApp,
	argName:  "run",
	argUsage: "-run To run the App as a services",
	run:      printInfo,
	cb:       printUsage,
}

const (
	// Year and copyright
	// http://patorjk.com/software/taag/#p=display&f=Stop&t=Seamless%20Wallet
	yc     = "(c) 2024-%v Yoga-Saputra"
	banner = `
      ______                  ______                             ______        _ _                   _                 
 / _____)                / _____)                           (____  \      (_) |                 | |      _         
| /  ___  ___     ___   | /  ___  ____ ____   ____    ___    ____)  ) ___  _| | ____  ____ ____ | | ____| |_  ____ 
| | (___)/ _ \   (___)  | | (___)/ ___)  _ \ / ___)  (___)  |  __  ( / _ \| | |/ _  )/ ___)  _ \| |/ _  |  _)/ _  )
| \____/| |_| |         | \____/| |   | | | ( (___          | |__)  ) |_| | | ( (/ /| |   | | | | ( ( | | |_( (/ / 
 \_____/ \___/           \_____/|_|   | ||_/ \____)         |______/ \___/|_|_|\____)_|   | ||_/|_|\_||_|\___)____)
                                      |_|                                                 |_|                      
%s %s`
)

func printInfo() {
	pyc := fmt.Sprintf(yc, time.Now().Year())
	header := fmt.Sprintf(pterm.LightGreen(banner), pterm.Red(boot.AppVersion), pterm.LightGreen(pyc))
	pterm.DefaultCenter.Println(header)

	additional := config.Of.App.Desc

	// App version and last build info
	lastBuild := "N/A"
	if boot.LastBuildAt != "" && boot.LastBuildAt != " " {
		lastBuild = boot.LastBuildAt
	}
	additional += fmt.Sprintf("\nLast Build Binary at: %v", lastBuild)

	// Print additional info
	pterm.DefaultCenter.WithCenterEachLineSeparately().Println(pterm.Cyan(additional))

	// Command list and usage headers
	fmt.Println(" Usage: -<argument>...")
	fmt.Println("")
	fmt.Println(" Arguments:")
}

func printUsage(commands map[string]cli) {
	var lists []pterm.BulletListItem
	for _, c := range commands {
		text := fmt.Sprintf("%v  [%v]", c.argName, c.argUsage)
		lists = append(lists, pterm.BulletListItem{
			Level: 2,
			Text:  text,
		})

		for _, v := range c.boolOptions {
			lists = append(lists, pterm.BulletListItem{
				Level: 4,
				Text:  fmt.Sprintf("%v  [%v]", v.optionName, v.optionUsage),
			})
		}
		for _, v := range c.float64Options {
			lists = append(lists, pterm.BulletListItem{
				Level: 4,
				Text:  fmt.Sprintf("%v  [%v]", v.optionName, v.optionUsage),
			})
		}
		for _, v := range c.intOptions {
			lists = append(lists, pterm.BulletListItem{
				Level: 4,
				Text:  fmt.Sprintf("%v  [%v]", v.optionName, v.optionUsage),
			})
		}
		for _, v := range c.stringOptions {
			lists = append(lists, pterm.BulletListItem{
				Level: 4,
				Text:  fmt.Sprintf("%v  [%v]", v.optionName, v.optionUsage),
			})
		}
		for _, v := range c.uintOptions {
			lists = append(lists, pterm.BulletListItem{
				Level: 4,
				Text:  fmt.Sprintf("%v  [%v]", v.optionName, v.optionUsage),
			})
		}
	}

	pterm.DefaultBulletList.WithItems(lists).Render()
}
