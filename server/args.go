package server

import (
	"fmt"
	"os"
)

func preBootstrapArgs() {
	args := os.Args[1:]
	i := 0
	count := len(args)
	for i < count {
		arg := args[i]

		if arg == "-d" || arg == "--data-dir" {
			if i == count-1 {
				fmt.Fprintf(os.Stderr, "%s requires exactly 1 parameter\n", arg)
				printHelpAndExit()
			}

			value := args[i+1]
			dataDirectory = value
			i++
		} else if arg == "-b" || arg == "--bind-addr" {
			if i == count-1 {
				fmt.Fprintf(os.Stderr, "%s requires exactly 1 parameter\n", arg)
				printHelpAndExit()
			}

			value := args[i+1]
			bindAddress = value
			i++
		} else if arg == "--no-schedule" {
			scheduleDisabled = true
		} else if arg == "-h" || arg == "--help" {
			printHelpAndExit()
		}

		i++
	}
}

func isVerbose() bool {
	args := os.Args[1:]
	for _, arg := range args {
		if arg == "-v" || arg == "--verbose" {
			return true
		}
	}
	return false
}

func postBootstrapArgs() {

}

func printHelpAndExit() {
	fmt.Printf("Usage %s [options]\n", os.Args[0])
	fmt.Printf("Options:\n")
	fmt.Printf("-d --data-dir <path>        Specify the absolute path to the data directory\n")
	fmt.Printf("-b --bind-addr <socket>     Specify the listen address for the web server\n")
	fmt.Printf("-v --verbose                Set the log level to debug\n")
	fmt.Printf("--no-schedule               Disable all automatic tasks\n")
	os.Exit(1)
}
