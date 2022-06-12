package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/ecnepsnai/otto"
	"golang.org/x/crypto/ssh"
)

type clientConfig struct {
	ListenAddr     string   `json:"listen_addr"`
	IdentityPath   string   `json:"identity_path"`
	ServerIdentity string   `json:"server_identity"`
	LogPath        string   `json:"log_path"`
	DefaultUID     uint32   `json:"default_uid"`
	DefaultGID     uint32   `json:"default_gid"`
	Path           string   `json:"path"`
	AllowFrom      []string `json:"allow_from"`
}

var config *clientConfig
var clientIdentity ssh.Signer

const (
	otto_CONFIG_FILE_NAME          = "otto_client.conf"
	otto_CONFIG_ATOMIC_FILE_NAME   = ".otto_client.conf.tmp"
	otto_IDENTITY_FILE_NAME        = ".otto_id.der"
	otto_IDENTITY_ATOMIC_FILE_NAME = ".otto_id.der.tmp"
)

func loadConfig() error {
	if _, err := os.Stat(otto_CONFIG_FILE_NAME); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "The otto client must be configured before use. See https://github.com/ecnepsnai/otto/blob/%s/docs/client.md for more information.\n\nUse -s to run interactive setup.\n", MainVersion)
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

	if config.IdentityPath == "" {
		return fmt.Errorf("empty identity path prohibited")
	}

	if config.ServerIdentity == "" {
		return fmt.Errorf("empty server identity prohibited")
	}

	if config.ListenAddr == "" {
		return fmt.Errorf("empty listen address prohibited")
	}

	for i, a := range c.AllowFrom {
		if _, _, err := net.ParseCIDR(a); err != nil {
			return fmt.Errorf("invalid allow_from CIDR '%s' at index %d: %s", a, i, err.Error())
		}
	}

	return nil
}

func mustLoadConfig() {
	if err := loadConfig(); err != nil {
		panic(err)
	}
}

func generateIdentity() error {
	id, err := otto.NewIdentity()
	if err != nil {
		log.PError("Error generating identity", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}

	f, err := os.OpenFile(otto_IDENTITY_ATOMIC_FILE_NAME, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		log.PError("Error opening identity file for writing", map[string]interface{}{
			"file_path": otto_IDENTITY_ATOMIC_FILE_NAME,
			"error":     err.Error(),
		})
		return err
	}

	if _, err := f.Write(id); err != nil {
		log.PError("Error writing identity file", map[string]interface{}{
			"file_path": otto_IDENTITY_ATOMIC_FILE_NAME,
			"error":     err.Error(),
		})
		f.Close()
		return err
	}
	f.Close()

	if err := os.Rename(otto_IDENTITY_ATOMIC_FILE_NAME, otto_IDENTITY_FILE_NAME); err != nil {
		log.PError("Error renaming identity file", map[string]interface{}{
			"old_name": otto_IDENTITY_ATOMIC_FILE_NAME,
			"new_name": otto_IDENTITY_FILE_NAME,
			"error":    err.Error(),
		})
		return err
	}

	return nil
}

func loadOrGenerateClientIdentity() (ssh.Signer, error) {
	_, err := os.Stat(otto_IDENTITY_FILE_NAME)
	if err != nil && os.IsNotExist(err) {
		if err := generateIdentity(); err != nil {
			return nil, err
		}
		return loadClientIdentity()
	}
	return loadClientIdentity()
}

func loadClientIdentity() (ssh.Signer, error) {
	info, err := os.Stat(otto_IDENTITY_FILE_NAME)
	if err != nil {
		if !os.IsNotExist(err) {
			log.PError("Error reading identity file", map[string]interface{}{
				"file_path": otto_IDENTITY_FILE_NAME,
				"error":     err.Error(),
			})
		}
		return nil, err
	}

	mode := info.Mode()
	if mode&32 != 0 || mode&16 != 0 || mode&8 != 0 || mode&4 != 0 || mode&2 != 0 || mode&1 != 0 {
		fmt.Fprintf(os.Stderr, "The client identity file can be accessed by other users, this is very dangerous! You should delete the '%s' file, restart the client, then re-trust the host on the Otto server.\n", otto_IDENTITY_FILE_NAME)
		os.Exit(1)
	}

	data, err := os.ReadFile(otto_IDENTITY_FILE_NAME)
	if err != nil {
		log.PError("Error reading identity file", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}
	id, err := otto.ParseIdentity(data)
	if err != nil {
		log.PError("Error reading identity file", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}
	return id.Signer(), nil
}

func mustLoadIdentity() {
	signer, err := loadOrGenerateClientIdentity()
	if err != nil {
		panic(err)
	}
	clientIdentity = signer
	log.Debug("Client identity loaded: %s", base64.StdEncoding.EncodeToString(clientIdentity.PublicKey().Marshal()))
}

func saveNewConfig(c clientConfig) error {
	f, err := os.OpenFile(otto_CONFIG_ATOMIC_FILE_NAME, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.PError("Error opening atomic file for writing", map[string]interface{}{
			"file_path": otto_CONFIG_ATOMIC_FILE_NAME,
			"error":     err.Error(),
		})
		return err
	}

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "    ")
	err = encoder.Encode(c)
	f.Close()
	if err != nil {
		log.PError("Error writing config JSON", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}

	if err := os.Rename(otto_CONFIG_ATOMIC_FILE_NAME, otto_CONFIG_FILE_NAME); err != nil {
		log.PError("Error renaming atomic config file", map[string]interface{}{
			"file_path":        otto_CONFIG_FILE_NAME,
			"atomic_file_path": otto_CONFIG_ATOMIC_FILE_NAME,
			"error":            err.Error(),
		})
		return err
	}

	return nil
}

func updateServerIdentity(newPublicKey string) error {
	f, err := os.OpenFile(otto_CONFIG_FILE_NAME, os.O_RDONLY, 0644)
	if err != nil {
		log.PError("Error updating server identity", map[string]interface{}{
			"where": "reading existing config file",
			"error": err.Error(),
		})
		return err
	}

	c := defaultConfig()
	err = json.NewDecoder(f).Decode(&c)
	f.Close()
	if err != nil {
		log.PError("Error updating server identity", map[string]interface{}{
			"where": "decoding existing config file",
			"error": err.Error(),
		})
		return err
	}

	oldPublicKey := c.ServerIdentity
	c.ServerIdentity = newPublicKey

	if err := saveNewConfig(c); err != nil {
		log.PError("Error updating server identity", map[string]interface{}{
			"where": "saving config file",
			"error": err.Error(),
		})
		return err
	}

	log.PWarn("Server identity updated", map[string]interface{}{
		"old_identity": oldPublicKey,
		"new_identity": newPublicKey,
	})
	return nil
}

func defaultConfig() clientConfig {
	return clientConfig{
		ListenAddr:   "0.0.0.0:12444",
		IdentityPath: otto_IDENTITY_FILE_NAME,
		LogPath:      ".",
		DefaultUID:   0,
		DefaultGID:   0,
		AllowFrom:    []string{"0.0.0.0/0", "::/0"},
	}
}

func getAllowFroms() []net.IPNet {
	nets := make([]net.IPNet, len(config.AllowFrom))
	for i, a := range config.AllowFrom {
		_, network, err := net.ParseCIDR(a)
		if err != nil {
			panic(fmt.Sprintf("invalid CIDR address: %s", a))
		}
		nets[i] = *network
	}

	return nets
}
