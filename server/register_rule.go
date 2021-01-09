package server

import (
	"regexp"

	"github.com/ecnepsnai/ds"
)

// RegisterRule describes a register rule
type RegisterRule struct {
	ID       string `ds:"primary"`
	Property string `ds:"index"`
	Pattern  string
	GroupID  string `ds:"index"`
}

func (s *registerruleStoreObject) AllRules() []RegisterRule {
	objects, err := s.Table.GetAll(&ds.GetOptions{Sorted: true, Ascending: true})
	if err != nil {
		log.Error("Error getting registration rules: %s", err.Error())
		return []RegisterRule{}
	}
	count := len(objects)
	if count == 0 {
		return []RegisterRule{}
	}
	rules := make([]RegisterRule, count)
	for i, object := range objects {
		rule, ok := object.(RegisterRule)
		if !ok {
			log.Error("Invalid object type for RegisterRule")
			return []RegisterRule{}
		}
		rules[i] = rule
	}
	return rules
}

func (s *registerruleStoreObject) RuleWithID(id string) *RegisterRule {
	object, err := s.Table.Get(id)
	if err != nil {
		log.Error("Error getting registration rule: %s", err.Error())
		return nil
	}
	if object == nil {
		log.Warn("No registration rule with ID: %s", id)
		return nil
	}

	rule, ok := object.(RegisterRule)
	if !ok {
		log.Error("Invalid object type for RegisterRule")
		return nil
	}

	return &rule
}

func (s *registerruleStoreObject) RulesForGroup(groupID string) []RegisterRule {
	objects, err := s.Table.GetIndex("GroupID", groupID, &ds.GetOptions{Sorted: true, Ascending: true})
	if err != nil {
		log.Error("Error getting registration rules: %s", err.Error())
		return []RegisterRule{}
	}
	count := len(objects)
	if count == 0 {
		return []RegisterRule{}
	}
	rules := make([]RegisterRule, count)
	for i, object := range objects {
		rule, ok := object.(RegisterRule)
		if !ok {
			log.Error("Invalid object type for RegisterRule")
			return []RegisterRule{}
		}
		rules[i] = rule
	}
	return rules
}

type newRegisterRuleParams struct {
	Property string
	Pattern  string
	GroupID  string
}

func (s *registerruleStoreObject) NewRule(params newRegisterRuleParams) (*RegisterRule, *Error) {
	if !IsRegisterRuleProperty(params.Property) {
		return nil, ErrorUser("Invalid rule property")
	}

	if _, err := regexp.Compile(params.Pattern); err != nil {
		return nil, ErrorUser("Invalid pattern regex")
	}

	if group, _ := GroupStore.GroupWithID(params.GroupID); group == nil {
		return nil, ErrorUser("Unknown group ID")
	}

	rule := RegisterRule{
		ID:       newID(),
		Property: params.Property,
		Pattern:  params.Pattern,
		GroupID:  params.GroupID,
	}

	if err := s.Table.Add(rule); err != nil {
		log.Error("Error adding new rule: %s", err.Error())
		return nil, ErrorFrom(err)
	}

	return &rule, nil
}

type editRegisterRuleParams struct {
	Property string
	Pattern  string
	GroupID  string
}

func (s *registerruleStoreObject) EditRule(id string, params editRegisterRuleParams) (*RegisterRule, *Error) {
	rule := s.RuleWithID(id)
	if rule == nil {
		return nil, ErrorUser("No rule with ID %s", id)
	}

	if !IsRegisterRuleProperty(params.Property) {
		return nil, ErrorUser("Invalid rule property")
	}

	if _, err := regexp.Compile(params.Pattern); err != nil {
		return nil, ErrorUser("Invalid pattern regex")
	}

	if group, _ := GroupStore.GroupWithID(params.GroupID); group == nil {
		return nil, ErrorUser("Unknown group ID")
	}

	rule.Property = params.Property
	rule.Pattern = params.Pattern
	rule.GroupID = params.GroupID
	if err := s.Table.Update(*rule); err != nil {
		log.Error("Error updating rule '%s': %s", rule.ID, err.Error())
		return nil, ErrorFrom(err)
	}

	return rule, nil
}

func (s *registerruleStoreObject) DeleteRule(id string) (*RegisterRule, *Error) {
	rule := s.RuleWithID(id)
	if rule == nil {
		return nil, ErrorUser("No rule with ID %s", id)
	}

	if err := s.Table.Delete(*rule); err != nil {
		log.Error("Error deleting group '%s': %s", id, err.Error())
		return nil, ErrorFrom(err)
	}

	return rule, nil
}
