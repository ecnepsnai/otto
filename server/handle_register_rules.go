package server

import "github.com/ecnepsnai/web"

func (h *handle) RegisterRuleList(request web.Request) (interface{}, *web.Error) {
	return RegisterRuleStore.AllRules(), nil
}

func (h *handle) RegisterRuleNew(request web.Request) (interface{}, *web.Error) {
	session := request.UserData.(*Session)

	params := newRegisterRuleParams{}
	if err := request.Decode(&params); err != nil {
		return nil, err
	}

	rule, err := RegisterRuleStore.NewRule(params)
	if err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}

	EventStore.RegisterRuleAdded(rule, session.Username)

	return rule, nil
}

func (h *handle) RegisterRuleGet(request web.Request) (interface{}, *web.Error) {
	id := request.Params.ByName("id")

	rule := RegisterRuleStore.RuleWithID(id)
	if rule == nil {
		return nil, web.ValidationError("No rule with ID %s", id)
	}

	return rule, nil
}

func (h *handle) RegisterRuleEdit(request web.Request) (interface{}, *web.Error) {
	session := request.UserData.(*Session)

	id := request.Params.ByName("id")
	params := editRegisterRuleParams{}
	if err := request.Decode(&params); err != nil {
		return nil, err
	}

	rule, err := RegisterRuleStore.EditRule(id, params)
	if err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}

	EventStore.RegisterRuleModified(rule, session.Username)

	return rule, nil
}

func (h *handle) RegisterRuleDelete(request web.Request) (interface{}, *web.Error) {
	session := request.UserData.(*Session)

	id := request.Params.ByName("id")
	rule, err := RegisterRuleStore.DeleteRule(id)
	if err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}

	EventStore.RegisterRuleDeleted(rule, session.Username)

	return true, nil
}
