package server

import (
	"fmt"
	"testing"

	"github.com/ecnepsnai/otto"
	"github.com/ecnepsnai/web"
)

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

	RegisterRuleStore.Table.DeleteAll()
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

	mockRequest := func(remoteAddr string, params otto.RegisterRequest) web.Request {
		request := web.MockRequest(nil, map[string]string{}, params)
		request.HTTP.Header.Add("X-OTTO-PROTO-VERSION", fmt.Sprintf("%d", otto.ProtocolVersion))
		request.HTTP.RemoteAddr = remoteAddr
		return request
	}

	// Test that a host that doesn't match any rules gets added to the default group
	defaultAddress := randomString(6)
	h := handle{}
	_, webErr := h.Register(mockRequest(defaultAddress, otto.RegisterRequest{
		Key:  Key,
		Port: 12444,
		Properties: otto.RegisterRequestProperties{
			Hostname:            randomString(6),
			KernelName:          randomString(6),
			KernelVersion:       randomString(6),
			DistributionName:    randomString(6),
			DistributionVersion: randomString(6),
		},
	}))
	if webErr != nil {
		t.Errorf("Unexpected error trying to register valid host")
	}
	if HostStore.HostWithAddress(defaultAddress) == nil {
		t.Errorf("Host was not registered")
	}

	// Test that a host that matches a rule gets added to the correct group
	centos8address := randomString(6)
	_, webErr = h.Register(mockRequest(centos8address, otto.RegisterRequest{
		Key:  Key,
		Port: 12444,
		Properties: otto.RegisterRequestProperties{
			Hostname:            randomString(6),
			KernelName:          randomString(6),
			KernelVersion:       randomString(6),
			DistributionName:    "CentOS Linux",
			DistributionVersion: "8",
		},
	}))
	if webErr != nil {
		t.Errorf("Unexpected error trying to register valid host")
	}
	if HostStore.HostWithAddress(centos8address) == nil {
		t.Errorf("Host was not registered")
	}
	if HostStore.HostWithAddress(centos8address).GroupIDs[0] != centos8group.ID {
		t.Errorf("Incorrect group")
	}

	// Test that an incorrect Key does not get registered
	incorrectKeyAddress := randomString(6)
	_, webErr = h.Register(mockRequest(incorrectKeyAddress, otto.RegisterRequest{
		Key:  randomString(6),
		Port: 12444,
		Properties: otto.RegisterRequestProperties{
			Hostname:            randomString(6),
			KernelName:          randomString(6),
			KernelVersion:       randomString(6),
			DistributionName:    "CentOS Linux",
			DistributionVersion: "8",
		},
	}))
	if webErr == nil {
		t.Errorf("No error seen when one expected")
	}
	if HostStore.HostWithAddress(incorrectKeyAddress) != nil {
		t.Errorf("Host was registered when one was not expected")
	}
}
