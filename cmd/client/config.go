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

func loadConfig() error {
	if _, err := os.Stat("otto_client.conf"); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "The otto client must be configured before use. See https://github.com/ecnepsnai/otto/blob/%s/docs/client.md for more information.\n", MainVersion)
		os.Exit(1)
	}

	f, err := os.OpenFile("otto_client.conf", os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	c := clientConfig{
		ListenAddr: "0.0.0.0:12444",
		LogPath:    ".",
		DefaultUID: 0,
		DefaultGID: 0,
		AllowFrom:  "0.0.0.0/0",
	}
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
		return fmt.Errorf("invalid allow_from CIDR")
	}

	return nil
}

func mustLoadConfig() {
	if err := loadConfig(); err != nil {
		panic(err)
	}
}
