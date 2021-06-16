package main

import (
	"fmt"
	"os"
	"strings"
)

func tryGuidedSetup() {
	config := defaultConfig()
	uid, gid := getCurrentUIDandGID()
	config.DefaultUID = uid
	config.DefaultGID = gid
	config.Path = os.Getenv("PATH")

	var getConfigValue func(string, string) string
	getConfigValue = func(label, defaultVal string) string {
		if defaultVal == "" {
			fmt.Printf("%s: ", label)
		} else {
			fmt.Printf("%s [%s]: ", label, defaultVal)
		}

		var result string
		fmt.Scanln(&result)
		result = strings.Trim(result, "\r\n")

		if result == "" && defaultVal != "" {
			result = defaultVal
		}

		if result == "" {
			return getConfigValue(label, defaultVal)
		}

		return result
	}

	config.ListenAddr = getConfigValue("Listen Address", config.ListenAddr)
	config.PSK = getConfigValue("Pre-Shared-Key", "")
	config.AllowFrom = getConfigValue("Allow Connections From", config.AllowFrom)

	saveNewConfig(config)
	fmt.Printf("Otto is now configured! Run %s to start the client\n", os.Args[0])
	os.Exit(0)
}
