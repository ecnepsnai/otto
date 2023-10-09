package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/ecnepsnai/otto/server"
)

func main() {
	killSignal := make(chan os.Signal, 1)
	signal.Notify(killSignal, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-killSignal
		stop()
	}()

	printSignal := make(chan os.Signal, 1)
	signal.Notify(printSignal, syscall.SIGTRAP)
	go func() {
		<-printSignal
		runtime.Breakpoint()
	}()
	start()
}

func start() {
	fmt.Printf("Otto %s (Built on '%s', Runtime %s). Copyright Ian Spence 2020-2022.\n", server.Version, server.BuildDate, runtime.Version())
	server.Start()
	stop()
}
func stop() {
	fmt.Printf("Shutting down\n")
	server.Stop()
	os.Exit(1)
}
