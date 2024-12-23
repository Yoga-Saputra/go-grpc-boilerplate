package boot

import (
	"fmt"
	"log"
	"net"
	"os"
	"sort"

	"github.com/Yoga-Saputra/go-grpc-boilerplate/config"
)

var (
	// Application name.
	// Default from config
	AppName = config.Of.App.Name

	// Description of the application
	// Default from config
	AppDesc = config.Of.App.Desc

	// Application version.
	// Will set from makefile, otherwise the default is v0.0.0
	AppVersion = "v.0.0.0"

	// Datetime of last application has been builded
	LastBuildAt string

	// To force set service into maintenance mode.
	// REST API nor gRPC service still available during soft maintenance.
	SoftMaintenance = "false"

	// To force set service into maintenance mode hard
	// REST API nor gRPC service will not available during hard maintenance.
	HardMaintenance = "false"
)

// AppArgs defines given argumen of any app services.
// If have new services, please add needed argument here,
// so the "Up()" function can be called on this code
type AppArgs struct {
	// Network listener for net server
	NL net.Listener
}

// Map of function that will be called on Up() method based on their order.
// If have new services, just create new file and their method and register here
var orderUp = map[int]func(arg *AppArgs){
	1: dbUp,
	2: kafkaProducerUp,
	3: cacheUp,
	4: rpcUp,
	5: queueUp,
}

// Map of function that will be called on Down() method based on their order.
// If have new services, just create new file and their method and register here
var orderDown = map[int]func(){
	1: rpcDown,
	2: dbDown,
	3: cacheDown,
	4: kafkaProducerDown,
}

// Map of function that will be called on FinishTask() method based on their order.
// If have new services, just create new file and their method and register here
var orderFinishTask = map[int]func(){
	1: queueDown,
}

// Up will turn up all services.
func Up(args *AppArgs, manual ...int) {
	log.Println("++--------------------[UP...]--------------------++")

	// Do manual start up if have manual argument given
	// Otherwise do auto
	if len(manual) > 0 {
		for _, v := range manual {
			orderUp[v](args)
		}
	} else {
		keys := make([]int, 0, len(orderUp))
		for k := range orderUp {
			keys = append(keys, k)
		}
		sort.Ints(keys)

		for _, k := range keys {
			orderUp[k](args)
		}
	}
}

// Down will turn down all services.
func Down(manual ...int) {
	log.Println("++--------------------[DOWN.]--------------------++")

	// Do manual down if have manual argument given
	// Otherwise do auto
	if len(manual) > 0 {
		for _, v := range manual {
			orderDown[v]()
		}
	} else {
		printOutDown(fmt.Sprintf("Running cleanup task on pId: %v ... \n", os.Getpid()))

		keys := make([]int, 0, len(orderDown))
		for k := range orderDown {
			keys = append(keys, k)
		}
		sort.Ints(keys)

		for _, k := range keys {
			orderDown[k]()
		}
	}
}

// FinishTask will tell the system to finish all before shutdown.
func FinishTask(manual ...int) {
	log.Println("++--------------------[FTask]--------------------++")

	// Do manual down if have manual argument given
	// Otherwise do auto
	if len(manual) > 0 {
		for _, v := range manual {
			orderFinishTask[v]()
		}
	} else {
		printOutFinishTask(fmt.Sprintf("Wait until all tasks has been finished on pId: %v ... \n", os.Getpid()))

		keys := make([]int, 0, len(orderFinishTask))
		for k := range orderFinishTask {
			keys = append(keys, k)
		}
		sort.Ints(keys)

		for _, k := range keys {
			orderFinishTask[k]()
		}
	}
}

// Helper function to print out message when servces up
func printOutUp(s string) {
	log.Printf("[UP...] - %v\n", s)
}

// Helper function to print out message when servces down
func printOutDown(s string) {
	log.Printf("[DOWN.] - %v\n", s)
}

// Helper function to print out message when servces wait finishing tasks
func printOutFinishTask(s string) {
	log.Printf("[FTask] - %v\n", s)
}
