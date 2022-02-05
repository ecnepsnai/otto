package server

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/ecnepsnai/otto"
)

func TestExecute(t *testing.T) {
	t.Parallel()

	version := randomString(4)
	var port uint32 = 12401
	clientId, _ := otto.NewIdentity()

	group, err := GroupStore.NewGroup(newGroupParameters{
		Name: randomString(6),
	})
	if err != nil {
		t.Fatalf("Error making group: %s", err.Message)
	}

	host, err := HostStore.NewHost(newHostParameters{
		Name:           randomString(6),
		Address:        "127.0.0.1",
		Port:           port,
		ClientIdentity: clientId.PublicKeyString(),
		GroupIDs:       []string{group.ID},
	})
	if err != nil {
		t.Fatalf("Error making host: %s", err.Message)
	}

	serverId := IdentityStore.Get(host.ID)

	go func() {
		_, allowFrom, _ := net.ParseCIDR("0.0.0.0/0")
		l, err := otto.SetupListener(otto.ListenOptions{
			Address:          fmt.Sprintf("127.0.0.1:%d", port),
			AllowFrom:        []net.IPNet{*allowFrom},
			Identity:         clientId.Signer(),
			TrustedPublicKey: serverId.PublicKeyString(),
		}, func(conn *otto.Connection) {
			conn.WriteMessage(otto.MessageTypeHeartbeatResponse, otto.MessageHeartbeatResponse{ClientVersion: version})
			conn.Close()
		})
		if err != nil {
			panic("error listening: " + err.Error())
		}
		l.Accept()
	}()

	time.Sleep(50 * time.Millisecond)

	if err := host.Ping(); err != nil {
		t.Errorf("Error pinging host: %s", err.Message)
	}

	hb := heartbeatStore.LastHeartbeat(host)
	if hb == nil {
		t.Fatalf("No heartbeat found for host")
	}

	if hb.Version != version {
		t.Errorf("Unexpected version for host. Expected '%s' got '%s'", version, hb.Version)
	}
}

func TestExecuteUntrustedKey(t *testing.T) {
	t.Parallel()

	version := randomString(4)
	var port uint32 = 12402
	clientId, _ := otto.NewIdentity()

	group, err := GroupStore.NewGroup(newGroupParameters{
		Name: randomString(6),
	})
	if err != nil {
		t.Fatalf("Error making group: %s", err.Message)
	}

	host, err := HostStore.NewHost(newHostParameters{
		Name:     randomString(6),
		Address:  "127.0.0.2",
		Port:     port,
		GroupIDs: []string{group.ID},
	})
	if err != nil {
		t.Fatalf("Error making host: %s", err.Message)
	}

	serverId := IdentityStore.Get(host.ID)

	go func() {
		_, allowFrom, _ := net.ParseCIDR("0.0.0.0/0")
		l, err := otto.SetupListener(otto.ListenOptions{
			Address:          fmt.Sprintf("127.0.0.2:%d", port),
			AllowFrom:        []net.IPNet{*allowFrom},
			Identity:         clientId.Signer(),
			TrustedPublicKey: serverId.PublicKeyString(),
		}, func(conn *otto.Connection) {
			conn.WriteMessage(otto.MessageTypeHeartbeatResponse, otto.MessageHeartbeatResponse{ClientVersion: version})
			conn.Close()
		})
		if err != nil {
			panic("error listening: " + err.Error())
		}
		l.Accept()
	}()

	time.Sleep(50 * time.Millisecond)

	host.Ping()

	host = HostCache.ByID(host.ID)

	if host.Trust.UntrustedIdentity != clientId.PublicKeyString() {
		t.Errorf("Unrecognized public key for host. Expected '%s' got '%s'", clientId.PublicKeyString(), host.Trust.UntrustedIdentity)
	}
}

func TestExecuteIncorrectIdentity(t *testing.T) {
	t.Parallel()

	version := randomString(4)
	var port uint32 = 12403
	clientId, _ := otto.NewIdentity()
	otherId, _ := otto.NewIdentity()

	group, err := GroupStore.NewGroup(newGroupParameters{
		Name: randomString(6),
	})
	if err != nil {
		t.Fatalf("Error making group: %s", err.Message)
	}

	host, err := HostStore.NewHost(newHostParameters{
		Name:           randomString(6),
		Address:        "127.0.0.3",
		Port:           port,
		ClientIdentity: otherId.PublicKeyString(),
		GroupIDs:       []string{group.ID},
	})
	if err != nil {
		t.Fatalf("Error making host: %s", err.Message)
	}

	serverId := IdentityStore.Get(host.ID)

	go func() {
		_, allowFrom, _ := net.ParseCIDR("0.0.0.0/0")
		l, err := otto.SetupListener(otto.ListenOptions{
			Address:          fmt.Sprintf("127.0.0.3:%d", port),
			AllowFrom:        []net.IPNet{*allowFrom},
			Identity:         clientId.Signer(),
			TrustedPublicKey: serverId.PublicKeyString(),
		}, func(conn *otto.Connection) {
			conn.WriteMessage(otto.MessageTypeHeartbeatResponse, otto.MessageHeartbeatResponse{ClientVersion: version})
			conn.Close()
		})
		if err != nil {
			panic("error listening: " + err.Error())
		}
		l.Accept()
	}()

	time.Sleep(50 * time.Millisecond)

	if host.Ping() == nil {
		t.Fatalf("No error found when one expected connecting to host with wrong key")
	}
}

func TestExecuteReconnect(t *testing.T) {
	t.Parallel()

	version := randomString(4)
	var port uint32 = 12404
	clientId, _ := otto.NewIdentity()

	group, err := GroupStore.NewGroup(newGroupParameters{
		Name: randomString(6),
	})
	if err != nil {
		t.Fatalf("Error making group: %s", err.Message)
	}

	host, err := HostStore.NewHost(newHostParameters{
		Name:           randomString(6),
		Address:        "127.0.0.4",
		Port:           port,
		ClientIdentity: clientId.PublicKeyString(),
		GroupIDs:       []string{group.ID},
	})
	if err != nil {
		t.Fatalf("Error making host: %s", err.Message)
	}

	serverId := IdentityStore.Get(host.ID)

	go func() {
		_, allowFrom, _ := net.ParseCIDR("0.0.0.0/0")
		l, err := otto.SetupListener(otto.ListenOptions{
			Address:          fmt.Sprintf("127.0.0.4:%d", port),
			AllowFrom:        []net.IPNet{*allowFrom},
			Identity:         clientId.Signer(),
			TrustedPublicKey: serverId.PublicKeyString(),
		}, func(conn *otto.Connection) {
			conn.ReadMessage()
			conn.WriteMessage(otto.MessageTypeHeartbeatResponse, otto.MessageHeartbeatResponse{ClientVersion: version})
			conn.Close()
		})
		if err != nil {
			panic("error listening: " + err.Error())
		}
		l.Accept()
	}()

	time.Sleep(50 * time.Millisecond)

	if err := host.Ping(); err != nil {
		t.Errorf("Error pinging host: %s", err.Message)
	}

	hb := heartbeatStore.LastHeartbeat(host)
	if hb == nil {
		t.Fatalf("No heartbeat found for host")
	}

	if hb.Version != version {
		t.Errorf("Unexpected version for host. Expected '%s' got '%s'", version, hb.Version)
	}

	time.Sleep(50 * time.Millisecond)

	if err := host.Ping(); err != nil {
		t.Errorf("Error pinging host: %s", err.Message)
	}

	hb = heartbeatStore.LastHeartbeat(host)
	if hb == nil {
		t.Fatalf("No heartbeat found for host")
	}

	if hb.Version != version {
		t.Errorf("Unexpected version for host. Expected '%s' got '%s'", version, hb.Version)
	}
}
