package server

import (
	"regexp"

	"github.com/ecnepsnai/ds"
	"github.com/ecnepsnai/limits"
	"github.com/ecnepsnai/otto/shared/otto"
)

// RegisterRuleClause describes a single clause for a register rule
type RegisterRuleClause struct {
	Property string
	Pattern  string
}

func (clause RegisterRuleClause) validate() *Error {
	if !IsRegisterRuleProperty(clause.Property) {
		return ErrorUser("Invalid rule property")
	}

	if _, err := regexp.Compile(clause.Pattern); err != nil {
		return ErrorUser("Invalid pattern regex")
	}

	return nil
}

// RegisterRule describes a register rule
type RegisterRule struct {
	ID      string               `ds:"primary"`
	Name    string               `ds:"unique" min:"1" max:"140"`
	Clauses []RegisterRuleClause `min:"1"`
	GroupID string               `ds:"index"`
}

// Matches does this rule match the given set of host properties
func (rule RegisterRule) Matches(properties otto.RegisterRequestProperties) bool {
	allClausesMatched := true
	for _, clause := range rule.Clauses {
		pattern, err := regexp.Compile(clause.Pattern)
		if err != nil {
			log.Error("Invalid registration rule regex: %s: %s", clause.Pattern, err.Error())
			allClausesMatched = false
			continue
		}

		switch clause.Property {
		case RegisterRulePropertyHostname:
			if !pattern.MatchString(properties.Hostname) {
				allClausesMatched = false
			}
		case RegisterRulePropertyKernelName:
			if !pattern.MatchString(properties.KernelName) {
				allClausesMatched = false
			}
		case RegisterRulePropertyKernelVersion:
			if !pattern.MatchString(properties.KernelVersion) {
				allClausesMatched = false
			}
		case RegisterRulePropertyDistributionName:
			if !pattern.MatchString(properties.DistributionName) {
				allClausesMatched = false
			}
		case RegisterRulePropertyDistributionVersion:
			if !pattern.MatchString(properties.DistributionVersion) {
				allClausesMatched = false
			}
		}

		if !allClausesMatched {
			break
		}
	}

	return allClausesMatched
}

func (s *registerruleStoreObject) AllRules() (rules []RegisterRule) {
	s.Table.StartRead(func(tx ds.IReadTransaction) error {
		rules = s.allRules(tx)
		return nil
	})
	return
}

func (s *registerruleStoreObject) allRules(tx ds.IReadTransaction) []RegisterRule {
	objects, err := tx.GetAll(&ds.GetOptions{Sorted: true, Ascending: true})
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

func (s *registerruleStoreObject) RuleWithID(id string) (rule *RegisterRule) {
	s.Table.StartRead(func(tx ds.IReadTransaction) error {
		rule = s.ruleWithID(tx, id)
		return nil
	})
	return
}

func (s *registerruleStoreObject) ruleWithID(tx ds.IReadTransaction, id string) *RegisterRule {
	object, err := tx.Get(id)
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

func (s *registerruleStoreObject) RuleWithName(name string) (rule *RegisterRule) {
	s.Table.StartRead(func(tx ds.IReadTransaction) error {
		rule = s.ruleWithName(tx, name)
		return nil
	})
	return
}

func (s *registerruleStoreObject) ruleWithName(tx ds.IReadTransaction, name string) *RegisterRule {
	object, err := tx.GetUnique("Name", name)
	if err != nil {
		log.Error("Error getting registration rule: %s", err.Error())
		return nil
	}
	if object == nil {
		return nil
	}

	rule, ok := object.(RegisterRule)
	if !ok {
		log.Error("Invalid object type for RegisterRule")
		return nil
	}

	return &rule
}

func (s *registerruleStoreObject) RulesForGroup(groupID string) (rules []RegisterRule) {
	s.Table.StartRead(func(tx ds.IReadTransaction) error {
		rules = s.rulesForGroup(tx, groupID)
		return nil
	})
	return
}

func (s *registerruleStoreObject) rulesForGroup(tx ds.IReadTransaction, groupID string) []RegisterRule {
	objects, err := tx.GetIndex("GroupID", groupID, &ds.GetOptions{Sorted: true, Ascending: true})
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
	Name    string
	Clauses []RegisterRuleClause
	GroupID string
}

func (s *registerruleStoreObject) NewRule(params newRegisterRuleParams) (rule *RegisterRule, err *Error) {
	s.Table.StartWrite(func(tx ds.IReadWriteTransaction) error {
		rule, err = s.newRule(tx, params)
		return nil
	})
	return
}

func (s *registerruleStoreObject) newRule(tx ds.IReadWriteTransaction, params newRegisterRuleParams) (*RegisterRule, *Error) {
	if len(params.Clauses) <= 0 {
		return nil, ErrorUser("Must include at least one clause")
	}

	for _, clause := range params.Clauses {
		if err := clause.validate(); err != nil {
			return nil, err
		}
	}

	if s.ruleWithName(tx, params.Name) != nil {
		return nil, ErrorUser("Rule with name %s already exists", params.Name)
	}

	if group := GroupCache.ByID(params.GroupID); group == nil {
		return nil, ErrorUser("Unknown group ID")
	}

	rule := RegisterRule{
		ID:      newID(),
		Name:    params.Name,
		Clauses: params.Clauses,
		GroupID: params.GroupID,
	}
	if err := limits.Check(rule); err != nil {
		return nil, ErrorUser(err.Error())
	}

	if err := tx.Add(rule); err != nil {
		log.Error("Error adding new rule: %s", err.Error())
		return nil, ErrorFrom(err)
	}

	return &rule, nil
}

type editRegisterRuleParams struct {
	Name    string
	Clauses []RegisterRuleClause
	GroupID string
}

func (s *registerruleStoreObject) EditRule(id string, params editRegisterRuleParams) (rule *RegisterRule, err *Error) {
	s.Table.StartWrite(func(tx ds.IReadWriteTransaction) error {
		rule, err = s.editRule(tx, id, params)
		return nil
	})
	return
}

func (s *registerruleStoreObject) editRule(tx ds.IReadWriteTransaction, id string, params editRegisterRuleParams) (*RegisterRule, *Error) {
	rule := s.ruleWithID(tx, id)
	if rule == nil {
		return nil, ErrorUser("No rule with ID %s", id)
	}

	if len(params.Clauses) <= 0 {
		return nil, ErrorUser("Must include at least one clause")
	}

	for _, clause := range params.Clauses {
		if err := clause.validate(); err != nil {
			return nil, err
		}
	}

	if existing := s.ruleWithName(tx, params.Name); existing != nil && existing.ID != id {
		return nil, ErrorUser("Rule with name %s already exists", params.Name)
	}

	if group := GroupCache.ByID(params.GroupID); group == nil {
		return nil, ErrorUser("Unknown group ID")
	}

	rule.Name = params.Name
	rule.Clauses = params.Clauses
	rule.GroupID = params.GroupID
	if err := limits.Check(rule); err != nil {
		return nil, ErrorUser(err.Error())
	}

	if err := tx.Update(*rule); err != nil {
		log.Error("Error updating rule '%s': %s", rule.ID, err.Error())
		return nil, ErrorFrom(err)
	}

	return rule, nil
}

func (s *registerruleStoreObject) DeleteRule(id string) (rule *RegisterRule, rerr *Error) {
	s.Table.StartWrite(func(tx ds.IReadWriteTransaction) error {
		rule = s.ruleWithID(tx, id)
		if rule == nil {
			rerr = ErrorUser("No rule with ID %s", id)
			return nil
		}

		if err := tx.Delete(*rule); err != nil {
			log.Error("Error deleting register rule '%s': %s", id, err.Error())
			rerr = ErrorFrom(err)
			return nil
		}
		log.Info("Deleted register rule '%s': %s", id, rule.Name)
		return nil
	})
	return
}
