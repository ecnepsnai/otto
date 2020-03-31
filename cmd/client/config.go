package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type clientConfig struct {
	PSK        string `json:"psk"`
	LogPath    string `json:"log_path"`
	DefaultUID uint32 `json:"default_uid"`
	DefaultGID uint32 `json:"default_gid"`
	Path       string `json:"path"`
}

var config *clientConfig

func loadConfig() error {
	f, err := os.OpenFile("otto_client.conf", os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	c := clientConfig{}
	if err := json.NewDecoder(f).Decode(&c); err != nil {
		return err
	}
	config = &c

	if config.PSK == "" {
		return fmt.Errorf("empty PSK prohibited")
	}

	return nil
}

func mustLoadConfig() {
	if err := loadConfig(); err != nil {
		panic(err)
	}
}
