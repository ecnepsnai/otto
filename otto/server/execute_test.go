package server

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/ecnepsnai/ds"
	"github.com/ecnepsnai/otto/server/environ"
	"github.com/ecnepsnai/otto/shared/otto"
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
	agentId, _ := otto.NewIdentity()

	group, err := GroupStore.NewGroup(newGroupParameters{
		Name: randomString(6),
	})
	if err != nil {
		t.Fatalf("Error making group: %s", err.Message)
	}

	host, err := HostStore.NewHost(newHostParameters{
		Name:          randomString(6),
		Address:       ip,
		Port:          uint32(secutil.RandomNumber(0, 65535)),
		AgentIdentity: agentId.PublicKeyString(),
		GroupIDs:      []string{group.ID},
	})
	if err != nil {
		t.Fatalf("Error making host: %s", err.Message)
	}

	_, allowFrom, _ := net.ParseCIDR("0.0.0.0/0")
	l, listenErr := otto.SetupListener(&otto.ListenOptions{
		Address:   ip + ":0",
		AllowFrom: []net.IPNet{*allowFrom},
		Identity:  agentId.Signer(),
		GetTrustedPublicKeys: func() []string {
			serverId, err := IdentityStore.Get(host.ID)
			if err != nil {
				panic(err)
			}
			if serverId == nil {
				return []string{}
			}
			return []string{serverId.PublicKeyString()}
		},
	}, func(conn *otto.Connection) {
		messageType, m, err := conn.ReadMessage()
		if err != nil {
			panic("Error reading message: " + err.Error())
		}
		if messageType != otto.MessageTypeHeartbeatRequest {
			panic("Incorrect message type")
		}
		message := m.(otto.MessageHeartbeatRequest)
		if err := conn.WriteMessage(otto.MessageTypeHeartbeatResponse, otto.MessageHeartbeatResponse{AgentVersion: version, Nonce: message.Nonce}); err != nil {
			t.Fatalf("Error writing message: " + err.Error())
		}
		conn.Close()
	})
	if listenErr != nil {
		panic("error listening: " + listenErr.Error())
	}

	host.Port = uint32(l.Port())
	HostStore.Table.StartWrite(func(tx ds.IReadWriteTransaction) error {
		return tx.Update(*host)
	})

	go func() {
		l.Accept()
	}()

	time.Sleep(100 * time.Millisecond)

	if err := host.Ping(); err != nil {
		t.Fatalf("Error pinging host: %s", err.Error())
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
	agentId, _ := otto.NewIdentity()

	environKey := randomString(4)
	environValue := randomString(4)

	scriptName := randomString(4)
	script, err := ScriptStore.NewScript(newScriptParameters{
		Name:       scriptName,
		Executable: "/bin/sh",
		Script:     "echo 'hello world'",
		RunLevel:   ScriptRunLevelReadOnly,
	})
	if err != nil {
		t.Fatalf("Error making script: %s", err.Message)
	}

	group, err := GroupStore.NewGroup(newGroupParameters{
		Name:      randomString(6),
		ScriptIDs: []string{script.ID},
		Environment: []environ.Variable{
			{
				Key:   environKey,
				Value: environValue,
			},
		},
	})
	if err != nil {
		t.Fatalf("Error making group: %s", err.Message)
	}

	host, err := HostStore.NewHost(newHostParameters{
		Name:          randomString(6),
		Address:       ip,
		Port:          uint32(secutil.RandomNumber(0, 65535)),
		AgentIdentity: agentId.PublicKeyString(),
		GroupIDs:      []string{group.ID},
	})
	if err != nil {
		t.Fatalf("Error making host: %s", err.Message)
	}

	_, allowFrom, _ := net.ParseCIDR("0.0.0.0/0")
	l, listenErr := otto.SetupListener(&otto.ListenOptions{
		Address:   ip + ":0",
		AllowFrom: []net.IPNet{*allowFrom},
		Identity:  agentId.Signer(),
		GetTrustedPublicKeys: func() []string {
			serverId, err := IdentityStore.Get(host.ID)
			if err != nil {
				panic(err)
			}
			if serverId == nil {
				return []string{}
			}
			return []string{serverId.PublicKeyString()}
		},
	}, func(conn *otto.Connection) {
		defer conn.Close()

		messageType, message, err := conn.ReadMessage()
		if err != nil {
			t.Fatalf("Error reading message: " + err.Error())
			return
		}
		if messageType != otto.MessageTypeTriggerActionRunScript {
			t.Fatalf("Incorrect message type")
			return
		}
		scriptInfo, ok := message.(otto.MessageTriggerActionRunScript)
		if !ok {
			t.Fatalf("Incorrect message type")
			return
		}
		if scriptInfo.Name != scriptName {
			t.Fatalf("Bad script name")
			return
		}
		if scriptInfo.Environment == nil {
			t.Fatalf("No environment variables passed")
			return
		}
		foundEnviron := false
		for key, value := range scriptInfo.Environment {
			if key == environKey && value == environValue {
				foundEnviron = true
			}
		}
		if !foundEnviron {
			t.Fatalf("Did not find environment variable")
			return
		}
		if err := conn.WriteMessage(otto.MessageTypeReadyForData, nil); err != nil {
			t.Fatalf("Error sending message: " + err.Error())
			return
		}
		var buf = make([]byte, scriptInfo.ScriptInfo.Length)
		conn.ReadData(buf)
		if err := conn.WriteMessage(otto.MessageTypeActionResult, otto.MessageActionResult{ScriptResult: otto.ScriptResult{Success: true}, AgentVersion: version}); err != nil {
			t.Fatalf("Error writing message: " + err.Error())
			return
		}
	})
	if listenErr != nil {
		panic("error listening: " + listenErr.Error())
	}

	host.Port = uint32(l.Port())
	HostStore.Table.StartWrite(func(tx ds.IReadWriteTransaction) error {
		return tx.Update(*host)
	})

	go func() {
		l.Accept()
	}()

	time.Sleep(100 * time.Millisecond)

	result, serr := host.RunScript(script, nil)
	if serr != nil {
		t.Fatalf("Error triggering action: %s", serr.Error())
	}

	if !result.Result.Success {
		t.Errorf("Unexpected result status: %+v", result.Result)
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
	agentId, _ := otto.NewIdentity()

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

	_, allowFrom, _ := net.ParseCIDR("0.0.0.0/0")
	l, listenErr := otto.SetupListener(&otto.ListenOptions{
		Address:   ip + ":0",
		AllowFrom: []net.IPNet{*allowFrom},
		Identity:  agentId.Signer(),
		GetTrustedPublicKeys: func() []string {
			serverId, err := IdentityStore.Get(host.ID)
			if err != nil {
				panic(err)
			}
			if serverId == nil {
				return []string{}
			}
			return []string{serverId.PublicKeyString()}
		},
	}, func(conn *otto.Connection) {
		if err := conn.WriteMessage(otto.MessageTypeHeartbeatResponse, otto.MessageHeartbeatResponse{AgentVersion: version}); err != nil {
			t.Fatalf("Error writing message: " + err.Error())
		}
		conn.Close()
	})
	if listenErr != nil {
		panic("error listening: " + listenErr.Error())
	}

	host.Port = uint32(l.Port())
	HostStore.Table.StartWrite(func(tx ds.IReadWriteTransaction) error {
		return tx.Update(*host)
	})

	go func() {
		l.Accept()
	}()

	time.Sleep(100 * time.Millisecond)

	host.Ping()

	host = HostCache.ByID(host.ID)

	if host.Trust.UntrustedIdentity != agentId.PublicKeyString() {
		t.Errorf("Unrecognized public key for host. Expected '%s' got '%s'", agentId.PublicKeyString(), host.Trust.UntrustedIdentity)
	}
}

func TestExecuteIncorrectIdentity(t *testing.T) {
	t.Parallel()

	ip := randLocalhostIP()
	version := randomString(4)
	agentId, _ := otto.NewIdentity()
	otherId, _ := otto.NewIdentity()

	group, err := GroupStore.NewGroup(newGroupParameters{
		Name: randomString(6),
	})
	if err != nil {
		t.Fatalf("Error making group: %s", err.Message)
	}

	host, err := HostStore.NewHost(newHostParameters{
		Name:          randomString(6),
		Address:       ip,
		Port:          uint32(secutil.RandomNumber(0, 65535)),
		AgentIdentity: otherId.PublicKeyString(),
		GroupIDs:      []string{group.ID},
	})
	if err != nil {
		t.Fatalf("Error making host: %s", err.Message)
	}

	_, allowFrom, _ := net.ParseCIDR("0.0.0.0/0")
	l, listenErr := otto.SetupListener(&otto.ListenOptions{
		Address:   ip + ":0",
		AllowFrom: []net.IPNet{*allowFrom},
		Identity:  agentId.Signer(),
		GetTrustedPublicKeys: func() []string {
			serverId, err := IdentityStore.Get(host.ID)
			if err != nil {
				panic(err)
			}
			if serverId == nil {
				return []string{}
			}
			return []string{serverId.PublicKeyString()}
		},
	}, func(conn *otto.Connection) {
		if err := conn.WriteMessage(otto.MessageTypeHeartbeatResponse, otto.MessageHeartbeatResponse{AgentVersion: version}); err != nil {
			t.Fatalf("Error writing message: " + err.Error())
		}
		conn.Close()
	})
	if listenErr != nil {
		panic("error listening: " + listenErr.Error())
	}

	host.Port = uint32(l.Port())
	HostStore.Table.StartWrite(func(tx ds.IReadWriteTransaction) error {
		return tx.Update(*host)
	})

	go func() {
		l.Accept()
	}()

	time.Sleep(100 * time.Millisecond)

	if host.Ping() == nil {
		t.Fatalf("No error found when one expected connecting to host with wrong key")
	}
}

func TestExecuteReconnect(t *testing.T) {
	t.Parallel()

	ip := randLocalhostIP()
	version := randomString(4)
	agentId, _ := otto.NewIdentity()

	group, err := GroupStore.NewGroup(newGroupParameters{
		Name: randomString(6),
	})
	if err != nil {
		t.Fatalf("Error making group: %s", err.Message)
	}

	host, err := HostStore.NewHost(newHostParameters{
		Name:          randomString(6),
		Address:       ip,
		Port:          uint32(secutil.RandomNumber(0, 65535)),
		AgentIdentity: agentId.PublicKeyString(),
		GroupIDs:      []string{group.ID},
	})
	if err != nil {
		t.Fatalf("Error making host: %s", err.Message)
	}

	_, allowFrom, _ := net.ParseCIDR("0.0.0.0/0")
	l, listenErr := otto.SetupListener(&otto.ListenOptions{
		Address:   ip + ":0",
		AllowFrom: []net.IPNet{*allowFrom},
		Identity:  agentId.Signer(),
		GetTrustedPublicKeys: func() []string {
			serverId, err := IdentityStore.Get(host.ID)
			if err != nil {
				panic(err)
			}
			if serverId == nil {
				return []string{}
			}
			return []string{serverId.PublicKeyString()}
		},
	}, func(conn *otto.Connection) {
		_, m, _ := conn.ReadMessage()
		message := m.(otto.MessageHeartbeatRequest)
		if err := conn.WriteMessage(otto.MessageTypeHeartbeatResponse, otto.MessageHeartbeatResponse{AgentVersion: version, Nonce: message.Nonce}); err != nil {
			t.Fatalf("Error writing message: " + err.Error())
		}
		conn.Close()
	})
	if listenErr != nil {
		panic("error listening: " + listenErr.Error())
	}

	host.Port = uint32(l.Port())
	HostStore.Table.StartWrite(func(tx ds.IReadWriteTransaction) error {
		return tx.Update(*host)
	})

	go func() {
		l.Accept()
	}()

	time.Sleep(100 * time.Millisecond)

	if err := host.Ping(); err != nil {
		t.Fatalf("Error pinging host: %s", err.Error())
	}

	hb := heartbeatStore.LastHeartbeat(host)
	if hb == nil {
		t.Fatalf("No heartbeat found for host")
	}

	if hb.Version != version {
		t.Errorf("Unexpected version for host. Expected '%s' got '%s'", version, hb.Version)
	}

	time.Sleep(100 * time.Millisecond)

	if err := host.Ping(); err != nil {
		t.Fatalf("Error pinging host: %s", err.Error())
	}

	hb = heartbeatStore.LastHeartbeat(host)
	if hb == nil {
		t.Fatalf("No heartbeat found for host")
	}

	if hb.Version != version {
		t.Errorf("Unexpected version for host. Expected '%s' got '%s'", version, hb.Version)
	}
}

func TestExecuteUpdateIdentity(t *testing.T) {
	t.Parallel()

	ip := randLocalhostIP()
	version := randomString(4)
	agentId, _ := otto.NewIdentity()

	group, err := GroupStore.NewGroup(newGroupParameters{
		Name: randomString(6),
	})
	if err != nil {
		t.Fatalf("Error making group: %s", err.Message)
	}

	host, err := HostStore.NewHost(newHostParameters{
		Name:          randomString(6),
		Address:       ip,
		Port:          uint32(secutil.RandomNumber(0, 65535)),
		AgentIdentity: agentId.PublicKeyString(),
		GroupIDs:      []string{group.ID},
	})
	if err != nil {
		t.Fatalf("Error making host: %s", err.Message)
	}

	_, allowFrom, _ := net.ParseCIDR("0.0.0.0/0")
	l, listenErr := otto.SetupListener(&otto.ListenOptions{
		Address:   ip + ":0",
		AllowFrom: []net.IPNet{*allowFrom},
		Identity:  agentId.Signer(),
		GetTrustedPublicKeys: func() []string {
			serverId, err := IdentityStore.Get(host.ID)
			if err != nil {
				panic(err)
			}
			if serverId == nil {
				return []string{}
			}
			return []string{serverId.PublicKeyString()}
		},
	}, func(conn *otto.Connection) {
		messageType, m, err := conn.ReadMessage()
		if err != nil {
			panic("Error reading message: " + err.Error())
		}
		if messageType != otto.MessageTypeHeartbeatRequest {
			panic("Incorrect message type")
		}
		message := m.(otto.MessageHeartbeatRequest)
		if err := conn.WriteMessage(otto.MessageTypeHeartbeatResponse, otto.MessageHeartbeatResponse{AgentVersion: version, Nonce: message.Nonce}); err != nil {
			t.Fatalf("Error writing message: " + err.Error())
		}
		conn.Close()
	})
	if listenErr != nil {
		panic("error listening: " + listenErr.Error())
	}

	host.Port = uint32(l.Port())
	HostStore.Table.StartWrite(func(tx ds.IReadWriteTransaction) error {
		return tx.Update(*host)
	})

	go func() {
		l.Accept()
	}()

	time.Sleep(100 * time.Millisecond)

	if err := host.Ping(); err != nil {
		t.Fatalf("Error pinging host: %s", err.Error())
	}

	hb := heartbeatStore.LastHeartbeat(host)
	if hb == nil {
		t.Fatalf("No heartbeat found for host")
	}

	if hb.Version != version {
		t.Errorf("Unexpected version for host. Expected '%s' got '%s'", version, hb.Version)
	}

	newId, idErr := otto.NewIdentity()
	if idErr != nil {
		panic(idErr)
	}

	IdentityStore.Set(host.ID, newId)

	if err := host.Ping(); err != nil {
		t.Fatalf("Error pinging host: %s", err.Error())
	}

	hb = heartbeatStore.LastHeartbeat(host)
	if hb == nil {
		t.Fatalf("No heartbeat found for host")
	}

	if hb.Version != version {
		t.Errorf("Unexpected version for host. Expected '%s' got '%s'", version, hb.Version)
	}
}
