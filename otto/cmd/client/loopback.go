package main

import (
	"encoding/base64"
	"time"

	"github.com/ecnepsnai/otto"
	"github.com/ecnepsnai/secutil"
)

var loopbackIdentity otto.Identity

func setupLoopback() {
	id, err := otto.NewIdentity()
	if err != nil {
		log.Fatal("Error making loopback identity: %s", err.Error())
	}
	loopbackIdentity = id
	log.Debug("Created loopback identity: %s", loopbackIdentity.PublicKeyString())
}

func startLoopbackRepeater() {
	for {
		time.Sleep(1 * time.Minute)
		sendLoopbackHeartbeat()
	}
}

func sendLoopbackHeartbeat() {
	c, err := otto.Dial(otto.DialOptions{
		Network:          "tcp",
		Address:          config.ListenAddr,
		Identity:         loopbackIdentity.Signer(),
		TrustedPublicKey: base64.StdEncoding.EncodeToString(clientIdentity.PublicKey().Marshal()),
		Timeout:          2 * time.Second,
	})
	if err != nil {
		log.Error("Error making loopback connection to otto client: %s", err.Error())
		log.Fatal("Otto client is not listening correctly, exiting...")
	}
	nonce := secutil.RandomString(8)
	resp, err := c.SendHeartbeat(otto.MessageHeartbeatRequest{
		Version: Version,
		Nonce:   nonce,
	})
	if err != nil {
		log.Error("Error sending loopback heartbeat to client: %s", err.Error())
		log.Fatal("Otto client is not listening correctly, exiting...")
	}
	if resp.Nonce != nonce {
		log.Error("Unexpected nonce in loopback heartbeat, this should never happen")
		log.Fatal("Otto client is not listening correctly, exiting...")
	}
	if resp.ClientVersion != Version {
		log.Error("Unexpected version in loopback heartbeat, this should never happen")
		log.Fatal("Otto client is not listening correctly, exiting...")
	}
	c.Close()
	log.Debug("Loopback heartbeat successful")
}
