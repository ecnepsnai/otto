package server

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/ecnepsnai/otto"
	"github.com/ecnepsnai/secutil"
)

func randLocalhostIP() string {
	return fmt.Sprintf("127.%d.%d.%d",
		secutil.RandomNumber(1, 254),
		secutil.RandomNumber(1, 254),
		secutil.RandomNumber(1, 254))
}

func TestExecutePing(t *testing.T) {
	t.Parallel()

	ip := randLocalhostIP()
	version := randomString(4)
	clientId, _ := otto.NewIdentity()

	group, err := GroupStore.NewGroup(newGroupParameters{
		Name: randomString(6),
	})
	if err != nil {
		t.Fatalf("Error making group: %s", err.Message)
	}

	host, err := HostStore.NewHost(newHostParameters{
		Name:           randomString(6),
		Address:        ip,
		Port:           uint32(secutil.RandomNumber(0, 65535)),
		ClientIdentity: clientId.PublicKeyString(),
		GroupIDs:       []string{group.ID},
	})
	if err != nil {
		t.Fatalf("Error making host: %s", err.Message)
	}

	serverId := IdentityStore.Get(host.ID)
	_, allowFrom, _ := net.ParseCIDR("0.0.0.0/0")
	l, listenErr := otto.SetupListener(otto.ListenOptions{
		Address:          ip + ":0",
		AllowFrom:        []net.IPNet{*allowFrom},
		Identity:         clientId.Signer(),
		TrustedPublicKey: serverId.PublicKeyString(),
	}, func(conn *otto.Connection) {
		messageType, _, err := conn.ReadMessage()
		if err != nil {
			panic("Error reading message: " + err.Error())
		}
		if messageType != otto.MessageTypeHeartbeatRequest {
			panic("Incorrect message type")
		}
		if err := conn.WriteMessage(otto.MessageTypeHeartbeatResponse, otto.MessageHeartbeatResponse{ClientVersion: version}); err != nil {
			t.Fatalf("Error writing message: " + err.Error())
		}
		conn.Close()
	})
	if listenErr != nil {
		panic("error listening: " + listenErr.Error())
	}

	host.Port = uint32(l.Port())
	HostStore.Table.Update(host)

	go func() {
		l.Accept()
	}()

	time.Sleep(50 * time.Millisecond)

	if err := host.Ping(); err != nil {
		t.Fatalf("Error pinging host: %s", err.Message)
	}

	hb := heartbeatStore.LastHeartbeat(host)
	if hb == nil {
		t.Fatalf("No heartbeat found for host")
	}

	if hb.Version != version {
		t.Errorf("Unexpected version for host. Expected '%s' got '%s'", version, hb.Version)
	}
}

func TestExecuteAction(t *testing.T) {
	t.Parallel()

	ip := randLocalhostIP()
	version := randomString(4)
	clientId, _ := otto.NewIdentity()

	group, err := GroupStore.NewGroup(newGroupParameters{
		Name: randomString(6),
	})
	if err != nil {
		t.Fatalf("Error making group: %s", err.Message)
	}

	host, err := HostStore.NewHost(newHostParameters{
		Name:           randomString(6),
		Address:        ip,
		Port:           uint32(secutil.RandomNumber(0, 65535)),
		ClientIdentity: clientId.PublicKeyString(),
		GroupIDs:       []string{group.ID},
	})
	if err != nil {
		t.Fatalf("Error making host: %s", err.Message)
	}

	actionType := otto.ActionRunScript

	serverId := IdentityStore.Get(host.ID)
	_, allowFrom, _ := net.ParseCIDR("0.0.0.0/0")
	l, listenErr := otto.SetupListener(otto.ListenOptions{
		Address:          ip + ":0",
		AllowFrom:        []net.IPNet{*allowFrom},
		Identity:         clientId.Signer(),
		TrustedPublicKey: serverId.PublicKeyString(),
	}, func(conn *otto.Connection) {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			panic("Error reading message: " + err.Error())
		}
		if messageType != otto.MessageTypeTriggerAction {
			panic("Incorrect message type")
		}
		action, ok := message.(otto.MessageTriggerAction)
		if !ok {
			panic("Incorrect message type")
		}
		if action.Action != actionType {
			panic("Incorrect action type")
		}
		if err := conn.WriteMessage(otto.MessageTypeActionResult, otto.MessageActionResult{ClientVersion: version}); err != nil {
			panic("Error writing message: " + err.Error())
		}
		conn.Close()
	})
	if listenErr != nil {
		panic("error listening: " + listenErr.Error())
	}

	host.Port = uint32(l.Port())
	HostStore.Table.Update(host)

	go func() {
		l.Accept()
	}()

	time.Sleep(50 * time.Millisecond)

	result, err := host.TriggerAction(otto.MessageTriggerAction{Action: actionType}, nil, nil)
	if err != nil {
		t.Fatalf("Error triggering action: %s", err.Message)
	}

	if result.ClientVersion != version {
		t.Errorf("Unexpected version for host. Expected '%s' got '%s'", version, result.ClientVersion)
	}

	hb := heartbeatStore.LastHeartbeat(host)
	if hb == nil {
		t.Fatalf("No heartbeat found for host")
	}
}

func TestExecuteUntrustedKey(t *testing.T) {
	t.Parallel()

	ip := randLocalhostIP()
	version := randomString(4)
	clientId, _ := otto.NewIdentity()

	group, err := GroupStore.NewGroup(newGroupParameters{
		Name: randomString(6),
	})
	if err != nil {
		t.Fatalf("Error making group: %s", err.Message)
	}

	host, err := HostStore.NewHost(newHostParameters{
		Name:     randomString(6),
		Address:  ip,
		Port:     uint32(secutil.RandomNumber(0, 65535)),
		GroupIDs: []string{group.ID},
	})
	if err != nil {
		t.Fatalf("Error making host: %s", err.Message)
	}

	serverId := IdentityStore.Get(host.ID)

	_, allowFrom, _ := net.ParseCIDR("0.0.0.0/0")
	l, listenErr := otto.SetupListener(otto.ListenOptions{
		Address:          ip + ":0",
		AllowFrom:        []net.IPNet{*allowFrom},
		Identity:         clientId.Signer(),
		TrustedPublicKey: serverId.PublicKeyString(),
	}, func(conn *otto.Connection) {
		if err := conn.WriteMessage(otto.MessageTypeHeartbeatResponse, otto.MessageHeartbeatResponse{ClientVersion: version}); err != nil {
			t.Fatalf("Error writing message: " + err.Error())
		}
		conn.Close()
	})
	if listenErr != nil {
		panic("error listening: " + listenErr.Error())
	}

	host.Port = uint32(l.Port())
	HostStore.Table.Update(host)

	go func() {
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

	ip := randLocalhostIP()
	version := randomString(4)
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
		Address:        ip,
		Port:           uint32(secutil.RandomNumber(0, 65535)),
		ClientIdentity: otherId.PublicKeyString(),
		GroupIDs:       []string{group.ID},
	})
	if err != nil {
		t.Fatalf("Error making host: %s", err.Message)
	}

	serverId := IdentityStore.Get(host.ID)
	_, allowFrom, _ := net.ParseCIDR("0.0.0.0/0")
	l, listenErr := otto.SetupListener(otto.ListenOptions{
		Address:          ip + ":0",
		AllowFrom:        []net.IPNet{*allowFrom},
		Identity:         clientId.Signer(),
		TrustedPublicKey: serverId.PublicKeyString(),
	}, func(conn *otto.Connection) {
		if err := conn.WriteMessage(otto.MessageTypeHeartbeatResponse, otto.MessageHeartbeatResponse{ClientVersion: version}); err != nil {
			t.Fatalf("Error writing message: " + err.Error())
		}
		conn.Close()
	})
	if listenErr != nil {
		panic("error listening: " + listenErr.Error())
	}

	host.Port = uint32(l.Port())
	HostStore.Table.Update(host)

	go func() {
		l.Accept()
	}()

	time.Sleep(50 * time.Millisecond)

	if host.Ping() == nil {
		t.Fatalf("No error found when one expected connecting to host with wrong key")
	}
}

func TestExecuteReconnect(t *testing.T) {
	t.Parallel()

	ip := randLocalhostIP()
	version := randomString(4)
	clientId, _ := otto.NewIdentity()

	group, err := GroupStore.NewGroup(newGroupParameters{
		Name: randomString(6),
	})
	if err != nil {
		t.Fatalf("Error making group: %s", err.Message)
	}

	host, err := HostStore.NewHost(newHostParameters{
		Name:           randomString(6),
		Address:        ip,
		Port:           uint32(secutil.RandomNumber(0, 65535)),
		ClientIdentity: clientId.PublicKeyString(),
		GroupIDs:       []string{group.ID},
	})
	if err != nil {
		t.Fatalf("Error making host: %s", err.Message)
	}

	serverId := IdentityStore.Get(host.ID)

	_, allowFrom, _ := net.ParseCIDR("0.0.0.0/0")
	l, listenErr := otto.SetupListener(otto.ListenOptions{
		Address:          ip + ":0",
		AllowFrom:        []net.IPNet{*allowFrom},
		Identity:         clientId.Signer(),
		TrustedPublicKey: serverId.PublicKeyString(),
	}, func(conn *otto.Connection) {
		conn.ReadMessage()
		if err := conn.WriteMessage(otto.MessageTypeHeartbeatResponse, otto.MessageHeartbeatResponse{ClientVersion: version}); err != nil {
			t.Fatalf("Error writing message: " + err.Error())
		}
		conn.Close()
	})
	if listenErr != nil {
		panic("error listening: " + listenErr.Error())
	}

	host.Port = uint32(l.Port())
	HostStore.Table.Update(host)

	go func() {
		l.Accept()
	}()

	time.Sleep(50 * time.Millisecond)

	if err := host.Ping(); err != nil {
		t.Fatalf("Error pinging host: %s", err.Message)
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
		t.Fatalf("Error pinging host: %s", err.Message)
	}

	hb = heartbeatStore.LastHeartbeat(host)
	if hb == nil {
		t.Fatalf("No heartbeat found for host")
	}

	if hb.Version != version {
		t.Errorf("Unexpected version for host. Expected '%s' got '%s'", version, hb.Version)
	}
}
