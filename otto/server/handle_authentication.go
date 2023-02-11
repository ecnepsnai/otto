package server

import (
	"net/http"
	"time"

	"github.com/ecnepsnai/limits"
	"github.com/ecnepsnai/web"
)

func (h *handle) Login(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
	type credentials struct {
		Username string
		Password string
	}

	login := credentials{}
	if err := request.DecodeJSON(&login); err != nil {
		return nil, nil, err
	}
	if err := limits.Check(&login); err != nil {
		return nil, nil, web.ValidationError(err.Error())
	}

	authenticationResult := authenticateUser(login.Username, []byte(login.Password), request.HTTP)
	login = credentials{}
	if authenticationResult == nil {
		return nil, nil, web.CommonErrors.Unauthorized
	}

	response := web.APIResponse{
		Cookies: []http.Cookie{
			{
				Name:     ottoSessionCookie,
				Value:    authenticationResult.SessionKey,
				SameSite: http.SameSiteStrictMode,
				Path:     "/",
				Expires:  time.Now().AddDate(0, 0, 1),
				Secure:   Options.Authentication.SecureOnly,
			},
		},
	}

	var statusCode = 0
	if authenticationResult.MustChangePassword {
		statusCode = 1
	}

	return statusCode, &response, nil
}

func (h *handle) Logout(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
	session := request.UserData.(*Session)

	SessionStore.DeleteSession(session)

	response := web.APIResponse{
		Cookies: []http.Cookie{
			{
				Name:    ottoSessionCookie,
				Value:   "",
				Path:    "/",
				Expires: time.Now().AddDate(0, 0, -1),
				Secure:  Options.Authentication.SecureOnly,
			},
		},
	}

	EventStore.UserLoggedOut(session.Username)

	return nil, &response, nil
}
