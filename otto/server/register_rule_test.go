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

func decryptRegistrationResponse(data []byte, key string) (*otto.RegisterResponse, error) {
	plainText, err := secutil.Encryption.AES_256_GCM.Decrypt(data, key)
	if err != nil {
		return nil, err
	}
	response := otto.RegisterResponse{}
	if err := json.Unmarshal(plainText, &response); err != nil {
		return nil, err
	}

	return &response, nil
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
	o := AutoRegisterOptions
	o.Enabled = true
	o.DefaultGroupID = defaultGroup.ID
	o.Key = Key
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
	nonce := secutil.RandomString(6)
	defaultAddress := "127.0.0.1"
	v := view{}
	webReply := v.Register(mockRequest(defaultAddress, Key, otto.RegisterRequest{
		AgentIdentity: mustIdentity().PublicKeyString(),
		Port:          12444,
		Nonce:         nonce,
		Properties: otto.RegisterRequestProperties{
			Hostname:            randomString(6),
			KernelName:          randomString(6),
			KernelVersion:       randomString(6),
			DistributionName:    randomString(6),
			DistributionVersion: randomString(6),
		},
	}))
	if webReply.Status != 200 {
		t.Fatalf("[default] Unexpected error trying to register valid host: HTTP %d", webReply.Status)
	}
	encryptedReply, erro := io.ReadAll(webReply.Reader)
	if erro != nil {
		t.Errorf("[default] error reading response: %s", erro.Error())
	}
	registerResponse, erro := decryptRegistrationResponse(encryptedReply, Key)
	if erro != nil {
		t.Errorf("[default] error reading response: %s", erro.Error())
	}
	if registerResponse.Nonce != nonce {
		t.Errorf("[default] bad nonce")
	}
	if HostStore.HostWithAddress(defaultAddress) == nil {
		t.Errorf("[default] Host was not registered")
	}

	// Test that a host that matches a rule gets added to the correct group
	centos8address := "127.0.0.2"
	if webReply := v.Register(mockRequest(centos8address, Key, otto.RegisterRequest{
		AgentIdentity: mustIdentity().PublicKeyString(),
		Port:          12444,
		Nonce:         secutil.RandomString(6),
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
		Nonce:         secutil.RandomString(6),
		Properties: otto.RegisterRequestProperties{
			Hostname:            randomString(6),
			KernelName:          randomString(6),
			KernelVersion:       randomString(6),
			DistributionName:    "CentOS Linux",
			DistributionVersion: "8",
		},
	})); webReply.Status == 200 {
		t.Fatalf("[incorrect] Unexpected HTTP status trying to register host with invalid key: HTTP %d", webReply.Status)
	}
	if HostStore.HostWithAddress(incorrectKeyAddress) != nil {
		t.Errorf("[incorrect] Host was registered when one was not expected")
	}

	// Test that a reused nonce does not get registered
	duplicateNonceAddress := "127.0.0.4"
	if webReply := v.Register(mockRequest(duplicateNonceAddress, Key, otto.RegisterRequest{
		AgentIdentity: mustIdentity().PublicKeyString(),
		Port:          12444,
		Nonce:         nonce,
		Properties: otto.RegisterRequestProperties{
			Hostname:            randomString(6),
			KernelName:          randomString(6),
			KernelVersion:       randomString(6),
			DistributionName:    "CentOS Linux",
			DistributionVersion: "8",
		},
	})); webReply.Status == 200 {
		t.Fatalf("[duplicate nonce] Unexpected HTTP status trying to register host with duplicate nonce: HTTP %d", webReply.Status)
	}
	if HostStore.HostWithAddress(duplicateNonceAddress) != nil {
		t.Errorf("[duplicate nonce] Host was registered when one was not expected")
	}

	// Test that a reused nonce does not get registered
	noNonceAddress := "127.0.0.4"
	if webReply := v.Register(mockRequest(noNonceAddress, Key, otto.RegisterRequest{
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
		t.Fatalf("[no nonce] Unexpected HTTP status trying to register host with no nonce: HTTP %d", webReply.Status)
	}
	if HostStore.HostWithAddress(noNonceAddress) != nil {
		t.Errorf("[no nonce] Host was registered when one was not expected")
	}
}
