package server

import (
	"sort"

	"github.com/ecnepsnai/limits"
	"github.com/ecnepsnai/web"
)

func (h *handle) UserList(request web.Request) (interface{}, *web.Error) {
	users := UserStore.AllUsers()
	sort.Slice(users, func(i int, j int) bool {
		return users[i].Username < users[j].Username
	})

	return users, nil
}

func (h *handle) UserGet(request web.Request) (interface{}, *web.Error) {
	username := request.Parameters["username"]

	user := UserStore.UserWithUsername(username)
	if user == nil {
		return nil, web.ValidationError("No user with Username %s", username)
	}

	return user, nil
}

func (h *handle) UserNew(request web.Request) (interface{}, *web.Error) {
	session := request.UserData.(*Session)

	params := newUserParameters{}
	if err := request.DecodeJSON(&params); err != nil {
		return nil, err
	}

	if params.Username == "system" {
		return nil, web.ValidationError("Username is reserved")
	}

	user, err := UserStore.NewUser(params)
	if err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}

	EventStore.UserAdded(user, session.Username)

	return user, nil
}

func (h *handle) UserEdit(request web.Request) (interface{}, *web.Error) {
	session := request.UserData.(*Session)

	username := request.Parameters["username"]

	user := UserStore.UserWithUsername(username)
	if user == nil {
		return nil, web.ValidationError("No user with Username %s", username)
	}

	params := editUserParameters{}
	if err := request.DecodeJSON(&params); err != nil {
		return nil, err
	}

	user, err := UserStore.EditUser(user, params)
	if err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}

	if params.Password != "" {
		// End all other sessions if the user changes their own password
		if user.Username == session.Username {
			SessionStore.EndAllOtherForUser(user.Username, session)
		} else {
			// End all sessions if somebody else changes a users password
			SessionStore.EndAllForUser(user.Username)
		}
	}

	EventStore.UserModified(user.Username, session.Username)

	return user, nil
}

func (h *handle) UserResetAPIKey(request web.Request) (interface{}, *web.Error) {
	session := request.UserData.(*Session)

	username := request.Parameters["username"]

	apiKey, err := UserStore.ResetAPIKey(username)
	if err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}

	EventStore.UserResetAPIKey(username, session.Username)
	return *apiKey, nil
}

func (h *handle) UserResetPassword(request web.Request) (interface{}, *web.Error) {
	session := request.UserData.(*Session)

	type changePasswordParameters struct {
		Password string `min:"1"`
	}

	params := changePasswordParameters{}
	if err := request.DecodeJSON(&params); err != nil {
		return nil, err
	}

	if err := limits.Check(params); err != nil {
		return nil, web.ValidationError(err.Error())
	}

	user, err := UserStore.ResetPassword(session.Username, []byte(params.Password))
	if err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}

	EventStore.UserResetPassword(session.Username)
	SessionStore.CompletePartialSession(session.Key)
	return user, nil
}

func (h *handle) UserDelete(request web.Request) (interface{}, *web.Error) {
	username := request.Parameters["username"]
	session := request.UserData.(*Session)
	if username == session.Username {
		return nil, web.ValidationError("Cannot delete own user")
	}

	user := UserStore.UserWithUsername(username)
	if user == nil {
		return nil, web.ValidationError("No user with Username %s", username)
	}

	if err := UserStore.DeleteUser(user); err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}

	SessionStore.EndAllForUser(user.Username)

	EventStore.UserDeleted(user.Username, session.Username)

	return true, nil
}
