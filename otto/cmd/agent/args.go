package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"runtime"

	"github.com/ecnepsnai/logtic"
	"github.com/ecnepsnai/otto/shared/otto"
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
			message := `Otto agent:
	Version: %s
	Built On: %s
	Revision: %s
	Protocol version: %d
	Go runtime: %s

Host Information:
	Hostname: %s
	Kernel Name: %s
	Kernel Version: %s
	Distribution Name: %s
	Distribution Version: %s
`
			fmt.Printf(message, Version, BuildDate, BuildRevision, otto.ProtocolVersion, runtime.Version(), registerProperties.Hostname, registerProperties.KernelName, registerProperties.KernelVersion, registerProperties.DistributionName, registerProperties.DistributionVersion)
			os.Exit(0)
		} else if arg == "-d" || arg == "--debug" {
			logtic.Log.Level = logtic.LevelDebug
		} else if arg == "-p" || arg == "--public-key" {
			signer, err := loadAgentIdentity()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading identity: %s\n", err.Error())
				os.Exit(1)
			}
			fmt.Printf("%s\n", base64.StdEncoding.EncodeToString(signer.PublicKey().Marshal()))
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
			if err := updateServerIdentity(key); err != nil {
				fmt.Fprintf(os.Stderr, "Error updating server identity: %s", err.Error())
				os.Exit(1)
			}
			os.Exit(0)
		} else {
			fmt.Printf(`Usage: %s [options]

Options:
-v --version              Print agent version and host information
-p --public-key           Print the agent public key
-s --setup                Start the interactive setup process
-t --trust-identity <id>  Trust the specified server identity

Environment variables:
OTTO_VERBOSE    If set with any value prints the log to stdout/stderr
OTTO_DEBUG      If set with any value increases the log verbosity to debug
OTTO_DIR        Override the path used to store data & logs. By default the working directory is used.
`, os.Args[0])
			os.Exit(1)
		}
		i++
	}
}
