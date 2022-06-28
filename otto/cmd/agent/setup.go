package main

import (
	"encoding/base64"
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

	signer, err := loadOrGenerateAgentIdentity()
	if err != nil {
		panic(err)
	}

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
	config.ServerIdentity = getConfigValue("Server Identity (Copy from Otto Server)", "")
	config.AllowFrom = strings.Split(getConfigValue("Allow Connections From (comma-separated list of CIDR addresses)", strings.Join(config.AllowFrom, ",")), ",")
	fmt.Printf("Agent identity: %s\n", base64.StdEncoding.EncodeToString(signer.PublicKey().Marshal()))

	saveNewConfig(config)
	fmt.Printf("Otto is now configured! Run %s to start the agent\n", os.Args[0])
	os.Exit(0)
}
