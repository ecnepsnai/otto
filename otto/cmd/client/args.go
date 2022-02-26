package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"runtime"

	"github.com/ecnepsnai/otto"
)

func parseArgs() {
	args := os.Args[1:]
	if len(args) == 0 {
		return
	}

	i := 0
	for i < len(args) {
		arg := args[i]
		if arg == "-v" || arg == "--version" {
			message := `Otto client:
	Version: %s
	Protocol version: %d
	Go runtime: %s

Host Information:
	Hostname: %s
	Kernel Name: %s
	Kernel Version: %s
	Distribution Name: %s
	Distribution Version: %s
`
			fmt.Printf(message, MainVersion, otto.ProtocolVersion, runtime.Version(), registerProperties.Hostname, registerProperties.KernelName, registerProperties.KernelVersion, registerProperties.DistributionName, registerProperties.DistributionVersion)
			os.Exit(0)
		} else if arg == "-p" || arg == "--public-key" {
			signer, err := loadClientIdentity()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading identity: %s\n", err.Error())
				os.Exit(1)
			}
			fmt.Printf("%s\n", base64.RawURLEncoding.EncodeToString(signer.PublicKey().Marshal()))
			os.Exit(0)
		} else if arg == "-s" || arg == "--setup" {
			tryGuidedSetup()
		} else if arg == "-t" || arg == "--trust-identity" {
			if i == len(args)-1 {
				fmt.Fprintf(os.Stderr, "Arg %s requires a value\n", arg)
				os.Exit(1)
			}
			key := args[i+1]
			i++
			updateServerIdentity(key)
			os.Exit(0)
		} else {
			fmt.Printf(`Usage: %s [options]

Options:
-v --version              Print client version and host information
-p --public-key           Print the client public key
-s --setup                Start the interactive setup process
-t --trust-identity <id>  Trust the specified server identity

Environment variables:
OTTO_VERBOSE    If set with any value, increases the verbosity of the client log
`, os.Args[0])
			os.Exit(1)
		}
		i++
	}
}