package main

import (
	"fmt"
	"net"
	"os"
	"path"
	"runtime"

	"github.com/ecnepsnai/logtic"
	"github.com/ecnepsnai/otto"
)

var log *logtic.Source

func main() {
	loadRegisterProperties()
	parseArgs()
	tryAutoRegister()
	mustLoadConfig()

	logtic.Log.FilePath = path.Join(config.LogPath, "otto_client.log")
	logtic.Log.Level = logtic.LevelWarn
	env := envMap()
	if _, verbose := env["OTTO_VERBOSE"]; verbose {
		logtic.Log.Level = logtic.LevelDebug
	}

	logtic.Open()
	log = logtic.Connect("otto")

	l, err := net.Listen("tcp", config.ListenAddr)
	if err != nil {
		panic(err)
	}
	log.Info("Otto client listening on %s", config.ListenAddr)
	for {
		c, err := l.Accept()
		if err != nil {
			continue
		}

		_, network, _ := net.ParseCIDR(config.AllowFrom)
		if !network.Contains(c.RemoteAddr().(*net.TCPAddr).IP) {
			log.Warn("Rejecting connection from server outside of allowed network: %s", c.RemoteAddr().String())
			c.Close()
			continue
		}

		go newServerConnection(c).Start()
	}
}

func parseArgs() {
	args := os.Args[1:]
	if len(args) == 0 {
		return
	}

	i := 0
	for i < len(args) {
		arg := args[i]
		if arg == "-v" || arg == "--version" {
			fmt.Printf("Otto Client:\n")
			fmt.Printf("    Version: %s\n    Protocol version: %d\n    Go Runtime: %s\n", MainVersion, otto.ProtocolVersion, runtime.Version())
			printOSInfo()
			os.Exit(0)
		} else if arg == "-s" {
			tryGuidedSetup()
		}
		i++
	}
}

func printOSInfo() {
	fmt.Printf("\nHost properties:\n")
	fmt.Printf("    Hostname: %s\n    Kernel Name: %s\n    Kernel Version: %s\n    Distribution Name: %s\n    Distribution Version: %s\n", registerProperties.Hostname, registerProperties.KernelName, registerProperties.KernelVersion, registerProperties.DistributionName, registerProperties.DistributionVersion)
}
