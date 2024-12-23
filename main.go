package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Yoga-Saputra/go-grpc-boilerplate/boot"
	"github.com/Yoga-Saputra/go-grpc-boilerplate/cli"
	"github.com/Yoga-Saputra/go-grpc-boilerplate/config"
	"github.com/jpillora/overseer"
	"github.com/jpillora/overseer/fetcher"
)

func main() {
	// Execute CLI app
	cli.Execute()

	// Run the app on the overseer
	overseer.Run(overseer.Config{
		Program: gracefulStart,
		Address: fmt.Sprintf(":%v", config.Of.App.Port),
		Fetcher: &fetcher.File{
			Path:     config.Of.App.ProgramFile,
			Interval: time.Duration(config.Of.App.AUWI) * time.Second,
		},
		Debug: config.Of.App.Debug(),
	})
}

// Starting system gracefully
func gracefulStart(s overseer.State) {
	boot.Up(&boot.AppArgs{
		NL: s.Listener,
	})

	// Listen OS Signal of interuption
	sigs := make(chan os.Signal, 1)
	signal.Notify(
		sigs,
		syscall.SIGTSTP,
		os.Interrupt,
		overseer.SIGTERM,
		overseer.SIGUSR1,
		overseer.SIGUSR2,
		syscall.SIGINT,
	)

	// Stop the app gracefully
	// As soon as posible after signal from OS that already registered listened
	<-sigs
	gracefulStop()
}

// Stoping system gracefully
func gracefulStop() {
	// Wait until current proccess finish their tasks
	boot.FinishTask()

	// Shutdown the process
	boot.Down()
}
