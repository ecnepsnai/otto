package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path"

	"github.com/ecnepsnai/otto/shared/otto"
	"golang.org/x/crypto/ssh"
)

type agentConfig struct {
	ListenAddr      string   `json:"listen_addr"`
	IdentityPath    string   `json:"identity_path"`
	ServerIdentity  string   `json:"server_identity"`
	LogPath         string   `json:"log_path"`
	DefaultUID      uint32   `json:"default_uid"`
	DefaultGID      uint32   `json:"default_gid"`
	Path            string   `json:"path"`
	AllowFrom       []string `json:"allow_from"`
	ScriptTimeout   *int64   `json:"script_timeout,omitempty"`
	RebootCommand   *string  `json:"reboot_command,omitempty"`
	ShutdownCommand *string  `json:"shutdown_command,omitempty"`
}

var config *agentConfig
var agentIdentity ssh.Signer

const (
	otto_CONFIG_FILE_NAME          = "otto_agent.conf"
	otto_CONFIG_ATOMIC_FILE_NAME   = ".otto_agent.conf.tmp"
	otto_IDENTITY_FILE_NAME        = ".otto_id.der"
	otto_IDENTITY_ATOMIC_FILE_NAME = ".otto_id.der.tmp"
)

var otto_DIR = "."

func loadConfig() error {
	identityLock.Lock()
	defer identityLock.Unlock()

	if _, err := os.Stat(path.Join(otto_DIR, otto_CONFIG_FILE_NAME)); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "The otto agent must be configured before use. See https://github.com/ecnepsnai/otto/blob/%s/docs/agent.md for more information.\n\nUse -s to run interactive setup.\n", Version)
		os.Exit(1)
	}

	f, err := os.OpenFile(path.Join(otto_DIR, otto_CONFIG_FILE_NAME), os.O_RDONLY, 0644)
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

	f, err := os.OpenFile(path.Join(otto_DIR, otto_IDENTITY_ATOMIC_FILE_NAME), os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		log.PError("Error opening identity file for writing", map[string]interface{}{
			"file_path": path.Join(otto_DIR, otto_IDENTITY_ATOMIC_FILE_NAME),
			"error":     err.Error(),
		})
		return err
	}

	if _, err := id.Write(f); err != nil {
		log.PError("Error writing identity file", map[string]interface{}{
			"file_path": path.Join(otto_DIR, otto_IDENTITY_ATOMIC_FILE_NAME),
			"error":     err.Error(),
		})
		f.Close()
		return err
	}
	f.Close()

	if err := os.Rename(path.Join(otto_DIR, otto_IDENTITY_ATOMIC_FILE_NAME), path.Join(otto_DIR, otto_IDENTITY_FILE_NAME)); err != nil {
		log.PError("Error renaming identity file", map[string]interface{}{
			"old_name": path.Join(otto_DIR, otto_IDENTITY_ATOMIC_FILE_NAME),
			"new_name": path.Join(otto_DIR, otto_IDENTITY_FILE_NAME),
			"error":    err.Error(),
		})
		return err
	}

	return nil
}

func loadOrGenerateAgentIdentity() (ssh.Signer, error) {
	_, err := os.Stat(path.Join(otto_DIR, otto_IDENTITY_FILE_NAME))
	if err != nil && os.IsNotExist(err) {
		if err := generateIdentity(); err != nil {
			return nil, err
		}
		return loadAgentIdentity()
	}
	return loadAgentIdentity()
}

func loadAgentIdentity() (ssh.Signer, error) {
	info, err := os.Stat(path.Join(otto_DIR, otto_IDENTITY_FILE_NAME))
	if err != nil {
		if !os.IsNotExist(err) {
			log.PError("Error reading identity file", map[string]interface{}{
				"file_path": path.Join(otto_DIR, otto_IDENTITY_FILE_NAME),
				"error":     err.Error(),
			})
		}
		return nil, err
	}

	mode := info.Mode()
	if mode&32 != 0 || mode&16 != 0 || mode&8 != 0 || mode&4 != 0 || mode&2 != 0 || mode&1 != 0 {
		fmt.Fprintf(os.Stderr, "The agent identity file can be accessed by other users, this is very dangerous! You should delete the '%s' file, restart the agent, then re-trust the host on the Otto server.\n", path.Join(otto_DIR, otto_IDENTITY_FILE_NAME))
		if os.Getenv("OTTO_VERY_DANGEROUS_IGNORE_IDENTITY_PERMISSIONS") == "" {
			os.Exit(1)
		}
	}

	data, err := os.ReadFile(path.Join(otto_DIR, otto_IDENTITY_FILE_NAME))
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
	signer, err := loadOrGenerateAgentIdentity()
	if err != nil {
		panic(err)
	}
	agentIdentity = signer
	log.Debug("Agent identity loaded: %s", base64.StdEncoding.EncodeToString(agentIdentity.PublicKey().Marshal()))
}

func formatJSON(c interface{}) ([]byte, error) {
	buf := &bytes.Buffer{}
	encoder := json.NewEncoder(buf)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(c); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func saveNewConfig(c agentConfig) error {
	data, err := formatJSON(c)
	if err != nil {
		log.PError("Error encoding config JSON", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}
	if err := os.WriteFile(path.Join(otto_DIR, otto_CONFIG_ATOMIC_FILE_NAME), data, 0644); err != nil {
		log.PError("Error writing config JSON", map[string]interface{}{
			"file_path": path.Join(otto_DIR, otto_CONFIG_ATOMIC_FILE_NAME),
			"error":     err.Error(),
		})
		return err
	}

	if err := os.Rename(path.Join(otto_DIR, otto_CONFIG_ATOMIC_FILE_NAME), path.Join(otto_DIR, otto_CONFIG_FILE_NAME)); err != nil {
		log.PError("Error renaming atomic config file", map[string]interface{}{
			"file_path":        path.Join(otto_DIR, otto_CONFIG_FILE_NAME),
			"atomic_file_path": path.Join(otto_DIR, otto_CONFIG_ATOMIC_FILE_NAME),
			"error":            err.Error(),
		})
		return err
	}

	return nil
}

func updateServerIdentity(newPublicKey string) error {
	identityLock.Lock()
	defer identityLock.Unlock()

	f, err := os.OpenFile(path.Join(otto_DIR, otto_CONFIG_FILE_NAME), os.O_RDONLY, 0644)
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

func defaultConfig() agentConfig {
	return agentConfig{
		ListenAddr:   "0.0.0.0:12444",
		IdentityPath: path.Join(otto_DIR, otto_IDENTITY_FILE_NAME),
		LogPath:      otto_DIR,
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
