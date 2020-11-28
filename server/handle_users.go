package server

import (
	"github.com/ecnepsnai/web"
)

func (h *handle) UserList(request web.Request) (interface{}, *web.Error) {
	users, err := UserStore.AllUsers()
	if err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}

	return users, nil
}

func (h *handle) UserGet(request web.Request) (interface{}, *web.Error) {
	username := request.Params.ByName("username")

	user, err := UserStore.UserWithUsername(username)
	if err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}
	if user == nil {
		return nil, web.ValidationError("No user with Username %s", username)
	}

	return user, nil
}

func (h *handle) UserNew(request web.Request) (interface{}, *web.Error) {
	session := request.UserData.(*Session)

	params := newUserParameters{}
	if err := request.Decode(&params); err != nil {
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

	username := request.Params.ByName("username")

	user, err := UserStore.UserWithUsername(username)
	if err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}
	if user == nil {
		return nil, web.ValidationError("No user with Username %s", username)
	}

	params := editUserParameters{}
	if err := request.Decode(&params); err != nil {
		return nil, err
	}

	user, err = UserStore.EditUser(user, params)
	if err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}

	EventStore.UserModified(user.Username, session.Username)

	return user, nil
}

func (h *handle) UserDelete(request web.Request) (interface{}, *web.Error) {
	username := request.Params.ByName("username")
	session := request.UserData.(*Session)
	if username == session.Username {
		return nil, web.ValidationError("Cannot delete own user")
	}

	user, err := UserStore.UserWithUsername(username)
	if err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}
	if user == nil {
		return nil, web.ValidationError("No user with Username %s", username)
	}

	if err := UserStore.DeleteUser(user); err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}

	EventStore.UserDeleted(user.Username, session.Username)

	return true, nil
}
