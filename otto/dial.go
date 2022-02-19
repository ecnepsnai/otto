package otto

import (
	"encoding/base64"
	"fmt"
	"net"
	"time"

	"golang.org/x/crypto/ssh"
)

// DialOptions describes options for dialing to a host
type DialOptions struct {
	Network          string
	Address          string
	Identity         ssh.Signer
	TrustedPublicKey string
	Timeout          time.Duration
}

// Dial will dial the host specified by the options and perform a SSH handshake with it.
func Dial(options DialOptions) (*Connection, error) {
	clientConfig := &ssh.ClientConfig{
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(options.Identity),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			log.PDebug("Handshake", map[string]interface{}{
				"public_key": base64.StdEncoding.EncodeToString(key.Marshal()),
			})
			if options.TrustedPublicKey == base64.StdEncoding.EncodeToString(key.Marshal()) {
				log.Debug("Recognized public key")
				return nil
			}
			log.PWarn("Rejecting connection from untrusted public key", map[string]interface{}{
				"public_key": base64.StdEncoding.EncodeToString(key.Marshal()),
			})
			return fmt.Errorf("unknown public key: %x", key.Marshal())
		},
		HostKeyAlgorithms: []string{ssh.KeyAlgoED25519},
		ClientVersion:     fmt.Sprintf("SSH-2.0-OTTO-%d", ProtocolVersion),
		Timeout:           options.Timeout,
	}

	log.PDebug("Dialing", map[string]interface{}{
		"network": options.Network,
		"address": options.Address,
		"timeout": options.Timeout.String(),
	})
	client, err := ssh.Dial(options.Network, options.Address, clientConfig)
	if err != nil {
		log.PError("Error connecting to host", map[string]interface{}{
			"address": options.Address,
			"error":   err.Error(),
		})
		return nil, err
	}

	log.PDebug("Opening channel", map[string]interface{}{
		"address":      options.Address,
		"channel_name": sshChannelName,
	})
	channel, _, err := client.OpenChannel(sshChannelName, nil)
	if err != nil {
		log.PError("Error connecting to host", map[string]interface{}{
			"address": options.Address,
			"error":   err.Error(),
		})
		return nil, err
	}
	log.PDebug("Connected to host", map[string]interface{}{
		"address": options.Address,
	})

	return &Connection{
		w:          channel,
		remoteAddr: client.RemoteAddr(),
		localAddr:  client.LocalAddr(),
	}, nil
}
