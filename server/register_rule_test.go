package server

import "testing"

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
		Property: RegisterRulePropertyHostname,
		Pattern:  randomString(6),
		GroupID:  group.ID,
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
		Property: RegisterRulePropertyHostname,
		Pattern:  randomString(6),
		GroupID:  group.ID,
	})
	if err != nil {
		t.Fatalf("Error making new rule: %s", err.Message)
	}
	if rule == nil {
		t.Fatalf("No rule returned")
	}

	_, err = RegisterRuleStore.EditRule(rule.ID, editRegisterRuleParams{
		Property: RegisterRulePropertyHostname,
		Pattern:  randomString(6),
		GroupID:  group.ID,
	})
	if err != nil {
		t.Fatalf("Error editing rule: %s", err.Message)
	}

	if RegisterRuleStore.RuleWithID(rule.ID).Pattern == rule.Pattern {
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
		Property: RegisterRulePropertyHostname,
		Pattern:  randomString(6),
		GroupID:  group.ID,
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
