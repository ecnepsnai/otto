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
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		stop()
	}()
	start()
}

func start() {
	fmt.Printf("Otto %s (Runtime %s). Copyright Ian Spence 2020.\n", server.ServerVersion, runtime.Version())
	server.Start()
	stop()
}
func stop() {
	fmt.Printf("Shutting down\n")
	server.Stop()
	os.Exit(1)
}
