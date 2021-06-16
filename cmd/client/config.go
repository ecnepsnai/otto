package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type clientConfig struct {
	ListenAddr string `json:"listen_addr"`
	PSK        string `json:"psk"`
	LogPath    string `json:"log_path"`
	DefaultUID uint32 `json:"default_uid"`
	DefaultGID uint32 `json:"default_gid"`
	Path       string `json:"path"`
	AllowFrom  string `json:"allow_from"`
}

var config *clientConfig

const (
	otto_CONFIG_FILE_NAME        = "otto_client.conf"
	otto_CONFIG_ATOMIC_FILE_NAME = ".otto_client.conf.tmp"
)

func loadConfig() error {
	if _, err := os.Stat(otto_CONFIG_FILE_NAME); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "The otto client must be configured before use. See https://github.com/ecnepsnai/otto/blob/%s/docs/client.md for more information.\n", MainVersion)
		os.Exit(1)
	}

	f, err := os.OpenFile(otto_CONFIG_FILE_NAME, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	c := defaultConfig()
	if err := json.NewDecoder(f).Decode(&c); err != nil {
		return err
	}
	config = &c

	if config.PSK == "" {
		return fmt.Errorf("empty PSK prohibited")
	}

	if config.ListenAddr == "" {
		return fmt.Errorf("empty listen address prohibited")
	}

	if _, _, err := net.ParseCIDR(c.AllowFrom); err != nil {
		return fmt.Errorf("invalid allow_from CIDR '%s': %s", c.AllowFrom, err.Error())
	}

	return nil
}

func mustLoadConfig() {
	if err := loadConfig(); err != nil {
		panic(err)
	}
}

func updatePSK(newPSK string) error {
	f, err := os.OpenFile(otto_CONFIG_FILE_NAME, os.O_RDONLY, 0644)
	if err != nil {
		log.PError("Error updating PSK", map[string]interface{}{
			"where": "reading existing config file",
			"error": err.Error(),
		})
		return err
	}

	c := defaultConfig()
	err = json.NewDecoder(f).Decode(&c)
	f.Close()
	if err != nil {
		log.PError("Error updating PSK", map[string]interface{}{
			"where": "decoding existing config file",
			"error": err.Error(),
		})
		return err
	}

	c.PSK = newPSK

	f, err = os.OpenFile(otto_CONFIG_ATOMIC_FILE_NAME, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.PError("Error updating PSK", map[string]interface{}{
			"where": "opening atomic file for writing",
			"error": err.Error(),
		})
		return err
	}

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "    ")
	err = encoder.Encode(c)
	f.Close()
	if err != nil {
		log.PError("Error updating PSK", map[string]interface{}{
			"where": "writing atomic file",
			"error": err.Error(),
		})
		return err
	}

	if err := os.Rename(otto_CONFIG_ATOMIC_FILE_NAME, otto_CONFIG_FILE_NAME); err != nil {
		log.PError("Error updating PSK", map[string]interface{}{
			"where": "renaming atomic config file",
			"error": err.Error(),
		})
		return err
	}

	log.Warn("PSK updated")
	return nil
}

func defaultConfig() clientConfig {
	return clientConfig{
		ListenAddr: "0.0.0.0:12444",
		LogPath:    ".",
		DefaultUID: 0,
		DefaultGID: 0,
		AllowFrom:  "0.0.0.0/0",
	}
}
