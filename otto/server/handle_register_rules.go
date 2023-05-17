package server

import (
	"fmt"

	"github.com/ecnepsnai/web"
)

func (h *handle) RegisterRuleList(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
	return RegisterRuleStore.AllRules(), nil, nil
}

func (h *handle) RegisterRuleNew(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
	session := request.UserData.(*Session)

	if !session.User().Permissions.CanModifyAutoregister {
		EventStore.UserPermissionDenied(session.User().Username, "Create register rule")
		return nil, nil, web.ValidationError("Permission denied")
	}

	params := newRegisterRuleParams{}
	if err := request.DecodeJSON(&params); err != nil {
		return nil, nil, err
	}

	rule, err := RegisterRuleStore.NewRule(params)
	if err != nil {
		if err.Server {
			return nil, nil, web.CommonErrors.ServerError
		}
		return nil, nil, web.ValidationError(err.Message)
	}

	EventStore.RegisterRuleAdded(rule, session.Username)

	return rule, nil, nil
}

func (h *handle) RegisterRuleGet(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
	id := request.Parameters["id"]

	rule := RegisterRuleStore.RuleWithID(id)
	if rule == nil {
		return nil, nil, web.ValidationError("No rule with ID %s", id)
	}

	return rule, nil, nil
}

func (h *handle) RegisterRuleEdit(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
	session := request.UserData.(*Session)
	id := request.Parameters["id"]

	if !session.User().Permissions.CanModifyAutoregister {
		EventStore.UserPermissionDenied(session.User().Username, fmt.Sprintf("Modify register rule %s", id))
		return nil, nil, web.ValidationError("Permission denied")
	}

	params := editRegisterRuleParams{}
	if err := request.DecodeJSON(&params); err != nil {
		return nil, nil, err
	}

	rule, err := RegisterRuleStore.EditRule(id, params)
	if err != nil {
		if err.Server {
			return nil, nil, web.CommonErrors.ServerError
		}
		return nil, nil, web.ValidationError(err.Message)
	}

	EventStore.RegisterRuleModified(rule, session.Username)

	return rule, nil, nil
}

func (h *handle) RegisterRuleDelete(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
	session := request.UserData.(*Session)
	id := request.Parameters["id"]

	if !session.User().Permissions.CanModifyAutoregister {
		EventStore.UserPermissionDenied(session.User().Username, fmt.Sprintf("Delete register rule %s", id))
		return nil, nil, web.ValidationError("Permission denied")
	}

	rule, err := RegisterRuleStore.DeleteRule(id)
	if err != nil {
		if err.Server {
			return nil, nil, web.CommonErrors.ServerError
		}
		return nil, nil, web.ValidationError(err.Message)
	}

	EventStore.RegisterRuleDeleted(rule, session.Username)

	return true, nil, nil
}

func (h *handle) AutoRegisterOptionsGet(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
	session := request.UserData.(*Session)
	if !session.User().Permissions.CanModifyAutoregister {
		EventStore.UserPermissionDenied(session.User().Username, "Access auto register settings")
		return nil, nil, web.ValidationError("Permission denied")
	}

	return AutoRegisterOptions, nil, nil
}

func (h *handle) AutoRegisterOptionsUpdate(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
	session := request.UserData.(*Session)
	if !session.User().Permissions.CanModifyAutoregister {
		EventStore.UserPermissionDenied(session.User().Username, "Modify auto register settings")
		return nil, nil, web.ValidationError("Permission denied")
	}

	options := OptionsRegister{}
	if err := request.DecodeJSON(&options); err != nil {
		return nil, nil, err
	}

	if err := options.Validate(); err != nil {
		return nil, nil, web.ValidationError(err.Error())
	}
	options.Save()
	return true, nil, nil
}
