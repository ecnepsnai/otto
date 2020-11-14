package server

import (
	"net/http"
	"time"

	"github.com/ecnepsnai/limits"
	"github.com/ecnepsnai/web"
)

func (h *handle) Login(request web.Request) (interface{}, *web.Error) {
	login := Credentials{}
	if err := request.Decode(&login); err != nil {
		return nil, err
	}
	if err := limits.Check(&login); err != nil {
		return nil, web.ValidationError(err.Error())
	}

	result, err := AuthenticateUser(login, request.HTTP)
	if err != nil {
		return nil, err
	}

	request.AddCookie(&http.Cookie{
		Name:     ottoSessionCookie,
		Value:    result.CookieValue,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		Expires:  time.Now().AddDate(0, 0, 1),
	})

	return result.Session, nil
}

func (h *handle) Logout(request web.Request) (interface{}, *web.Error) {
	session := request.UserData.(*Session)

	SessionStore.DeleteSession(session)
	request.AddCookie(&http.Cookie{
		Name:    ottoSessionCookie,
		Value:   "",
		Path:    "/",
		Expires: time.Now().AddDate(0, 0, -1),
	})

	return nil, nil
}
