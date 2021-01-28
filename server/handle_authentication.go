package server

import (
	"net/http"
	"time"

	"github.com/ecnepsnai/limits"
	"github.com/ecnepsnai/web"
)

func (h *handle) Login(request web.Request) (interface{}, *web.Error) {
	type credentials struct {
		Username string
		Password string
	}

	login := credentials{}
	if err := request.Decode(&login); err != nil {
		return nil, err
	}
	if err := limits.Check(&login); err != nil {
		return nil, web.ValidationError(err.Error())
	}

	sessionKey := authenticateUser(login.Username, login.Password, request.HTTP)
	login.Password = ""
	login = credentials{}
	if sessionKey == nil {
		return nil, web.CommonErrors.Unauthorized
	}

	request.AddCookie(&http.Cookie{
		Name:     ottoSessionCookie,
		Value:    *sessionKey,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		Expires:  time.Now().AddDate(0, 0, 1),
		Secure:   Options.Authentication.SecureOnly,
	})

	return true, nil
}

func (h *handle) Logout(request web.Request) (interface{}, *web.Error) {
	session := request.UserData.(*Session)

	SessionStore.DeleteSession(session)
	request.AddCookie(&http.Cookie{
		Name:    ottoSessionCookie,
		Value:   "",
		Path:    "/",
		Expires: time.Now().AddDate(0, 0, -1),
		Secure:  Options.Authentication.SecureOnly,
	})

	EventStore.UserLoggedOut(session.Username)

	return nil, nil
}
