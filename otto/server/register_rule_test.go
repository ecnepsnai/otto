package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/ecnepsnai/ds"
	"github.com/ecnepsnai/otto/shared/otto"
	"github.com/ecnepsnai/secutil"
	"github.com/ecnepsnai/web"
)

func encryptRegistrationRequest(request otto.RegisterRequest, key string) ([]byte, error) {
	decryptedData, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	return secutil.Encryption.AES_256_GCM.Encrypt(decryptedData, key)
}

func TestAddGetRegisterRule(t *testing.T) {
	group, err := GroupStore.NewGroup(newGroupParameters{
		Name:      randomString(6),
		ScriptIDs: []string{},
	})
	if err != nil {
		t.Fatalf("Error making new group: %s", err.Message)
	}
	if group == nil {
		t.Fatalf("Should return a group")
	}

	rule, err := RegisterRuleStore.NewRule(newRegisterRuleParams{
		Name: randomString(12),
		Clauses: []RegisterRuleClause{
			{
				Property: RegisterRulePropertyHostname,
				Pattern:  randomString(6),
			},
		},
		GroupID: group.ID,
	})
	if err != nil {
		t.Fatalf("Error making new rule: %s", err.Message)
	}
	if rule == nil {
		t.Fatalf("No rule returned")
	}

	if RegisterRuleStore.RuleWithID(rule.ID) == nil {
		t.Fatalf("Should return rule")
	}
	if len(RegisterRuleStore.RulesForGroup(group.ID)) == 0 {
		t.Fatalf("Should return rule")
	}
}

func TestEditRegisterRule(t *testing.T) {
	group, err := GroupStore.NewGroup(newGroupParameters{
		Name:      randomString(6),
		ScriptIDs: []string{},
	})
	if err != nil {
		t.Fatalf("Error making new group: %s", err.Message)
	}
	if group == nil {
		t.Fatalf("Should return a group")
	}

	rule, err := RegisterRuleStore.NewRule(newRegisterRuleParams{
		Name: randomString(12),
		Clauses: []RegisterRuleClause{
			{
				Property: RegisterRulePropertyHostname,
				Pattern:  randomString(6),
			},
		},
		GroupID: group.ID,
	})
	if err != nil {
		t.Fatalf("Error making new rule: %s", err.Message)
	}
	if rule == nil {
		t.Fatalf("No rule returned")
	}

	_, err = RegisterRuleStore.EditRule(rule.ID, editRegisterRuleParams{
		Name: rule.Name,
		Clauses: []RegisterRuleClause{
			{
				Property: RegisterRulePropertyHostname,
				Pattern:  randomString(6),
			},
		},
		GroupID: group.ID,
	})
	if err != nil {
		t.Fatalf("Error editing rule: %s", err.Message)
	}

	if RegisterRuleStore.RuleWithID(rule.ID).Clauses[0].Pattern == rule.Clauses[0].Pattern {
		t.Fatalf("Should change pattern")
	}
}

func TestDeleteRegisterRule(t *testing.T) {
	group, err := GroupStore.NewGroup(newGroupParameters{
		Name:      randomString(6),
		ScriptIDs: []string{},
	})
	if err != nil {
		t.Fatalf("Error making new group: %s", err.Message)
	}
	if group == nil {
		t.Fatalf("Should return a group")
	}

	rule, err := RegisterRuleStore.NewRule(newRegisterRuleParams{
		Name: randomString(12),
		Clauses: []RegisterRuleClause{
			{
				Property: RegisterRulePropertyHostname,
				Pattern:  randomString(6),
			},
		},
		GroupID: group.ID,
	})
	if err != nil {
		t.Fatalf("Error making new rule: %s", err.Message)
	}
	if rule == nil {
		t.Fatalf("No rule returned")
	}

	if _, err := RegisterRuleStore.DeleteRule(rule.ID); err != nil {
		t.Fatalf("Error deleing rule: %s", err.Message)
	}

	if RegisterRuleStore.RuleWithID(rule.ID) != nil {
		t.Fatalf("Should not return rule")
	}
}

func TestRegisterRuleEndToEnd(t *testing.T) {
	Key := randomString(6)

	defaultGroup, err := GroupStore.NewGroup(newGroupParameters{
		Name:      randomString(6),
		ScriptIDs: []string{},
	})
	if err != nil {
		t.Fatalf("Error making new group: %s", err.Message)
	}
	centos7group, err := GroupStore.NewGroup(newGroupParameters{
		Name:      randomString(6),
		ScriptIDs: []string{},
	})
	if err != nil {
		t.Fatalf("Error making new group: %s", err.Message)
	}
	centos8group, err := GroupStore.NewGroup(newGroupParameters{
		Name:      randomString(6),
		ScriptIDs: []string{},
	})
	if err != nil {
		t.Fatalf("Error making new group: %s", err.Message)
	}

	RegisterRuleStore.Table.StartWrite(func(tx ds.IReadWriteTransaction) error {
		return tx.DeleteAll()
	})
	if _, err := RegisterRuleStore.NewRule(newRegisterRuleParams{
		Name: "CentOS 7 Hosts",
		Clauses: []RegisterRuleClause{
			{
				Property: RegisterRulePropertyDistributionName,
				Pattern:  "CentOS Linux",
			},
			{
				Property: RegisterRulePropertyDistributionVersion,
				Pattern:  "7",
			},
		},
		GroupID: centos7group.ID,
	}); err != nil {
		t.Fatalf("Error making new rule: %s", err.Message)
	}
	if _, err := RegisterRuleStore.NewRule(newRegisterRuleParams{
		Name: "CentOS 8 Hosts",
		Clauses: []RegisterRuleClause{
			{
				Property: RegisterRulePropertyDistributionName,
				Pattern:  "CentOS Linux",
			},
			{
				Property: RegisterRulePropertyDistributionVersion,
				Pattern:  "8",
			},
		},
		GroupID: centos8group.ID,
	}); err != nil {
		t.Fatalf("Error making new rule: %s", err.Message)
	}
	o := Options
	o.Register.Enabled = true
	o.Register.DefaultGroupID = defaultGroup.ID
	o.Register.Key = Key
	o.Save()

	mockRequest := func(remoteAddr string, key string, params otto.RegisterRequest) web.Request {
		encryptedData, err := encryptRegistrationRequest(params, key)
		if err != nil {
			panic(err)
		}

		h := http.Header{}
		h.Set("User-Agent", "go test")
		h.Set("X-OTTO-PROTO-VERSION", fmt.Sprintf("%d", otto.ProtocolVersion))
		h.Set("X-Real-Ip", remoteAddr)

		return web.Request{
			HTTP: &http.Request{
				Body:       io.NopCloser(bytes.NewReader(encryptedData)),
				Header:     h,
				RemoteAddr: remoteAddr,
			},
		}
	}

	mustIdentity := func() *otto.Identity {
		id, err := otto.NewIdentity()
		if err != nil {
			panic(err)
		}
		return id
	}

	// Test that a host that doesn't match any rules gets added to the default group
	defaultAddress := "127.0.0.1"
	v := view{}
	if webReply := v.Register(mockRequest(defaultAddress, Key, otto.RegisterRequest{
		AgentIdentity: mustIdentity().PublicKeyString(),
		Port:          12444,
		Properties: otto.RegisterRequestProperties{
			Hostname:            randomString(6),
			KernelName:          randomString(6),
			KernelVersion:       randomString(6),
			DistributionName:    randomString(6),
			DistributionVersion: randomString(6),
		},
	})); webReply.Status != 200 {
		t.Fatalf("[default] Unexpected error trying to register valid host: HTTP %d", webReply.Status)
	}
	if HostStore.HostWithAddress(defaultAddress) == nil {
		t.Errorf("[default] Host was not registered")
	}

	// Test that a host that matches a rule gets added to the correct group
	centos8address := "127.0.0.2"
	if webReply := v.Register(mockRequest(centos8address, Key, otto.RegisterRequest{
		AgentIdentity: mustIdentity().PublicKeyString(),
		Port:          12444,
		Properties: otto.RegisterRequestProperties{
			Hostname:            randomString(6),
			KernelName:          randomString(6),
			KernelVersion:       randomString(6),
			DistributionName:    "CentOS Linux",
			DistributionVersion: "8",
		},
	})); webReply.Status != 200 {
		t.Fatalf("[centos 8] Unexpected error trying to register valid host: HTTP %d", webReply.Status)
	}
	if HostStore.HostWithAddress(centos8address) == nil {
		t.Errorf("[centos 8] Host was not registered")
	}
	if HostStore.HostWithAddress(centos8address).GroupIDs[0] != centos8group.ID {
		t.Errorf("[centos 8] Incorrect group")
	}

	// Test that an incorrect Key does not get registered
	incorrectKeyAddress := "127.0.0.3"
	if webReply := v.Register(mockRequest(incorrectKeyAddress, randomString(12), otto.RegisterRequest{
		AgentIdentity: mustIdentity().PublicKeyString(),
		Port:          12444,
		Properties: otto.RegisterRequestProperties{
			Hostname:            randomString(6),
			KernelName:          randomString(6),
			KernelVersion:       randomString(6),
			DistributionName:    "CentOS Linux",
			DistributionVersion: "8",
		},
	})); webReply.Status == 200 {
		t.Fatalf("[incorrect] Unexpected error trying to register valid host: HTTP %d", webReply.Status)
	}
	if HostStore.HostWithAddress(incorrectKeyAddress) != nil {
		t.Errorf("[incorrect] Host was registered when one was not expected")
	}
}
