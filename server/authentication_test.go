package server

import (
	"net/http"
	"net/url"
	"testing"
	"time"
)

func mockHTTPRequest(urlString string, sessionCookieValue string) *http.Request {
	req, err := http.NewRequest("GET", urlString, nil)
	if err != nil {
		panic("invalid request")
	}

	req.AddCookie(&http.Cookie{
		Name:  ottoSessionCookie,
		Value: sessionCookieValue,
	})

	return req
}

func TestAuthenticationLogin(t *testing.T) {
	username := randomString(6)
	password := randomString(6)
	if _, err := UserStore.NewUser(newUserParameters{
		Username: username,
		Email:    randomString(6),
		Password: password,
	}); err != nil {
		t.Fatalf("Error making user: %s", err.Message)
	}

	authenticationResult := authenticateUser(username, password, &http.Request{RemoteAddr: randomString(6)})
	if authenticationResult == nil {
		t.Fatalf("Should return a session key")
	}

	session := sessionForHTTPRequest(mockHTTPRequest("/", authenticationResult.SessionKey), false)
	if session == nil {
		t.Fatalf("Empty session for valid cookie")
	}
}

func TestAuthenticationIncorrectPassword(t *testing.T) {
	username := randomString(6)
	password := randomString(6)
	if _, err := UserStore.NewUser(newUserParameters{
		Username: username,
		Email:    randomString(6),
		Password: password,
	}); err != nil {
		t.Fatalf("Error making user: %s", err.Message)
	}

	sessionKey := authenticateUser(username, randomString(6), &http.Request{RemoteAddr: randomString(6)})
	if sessionKey != nil {
		t.Fatalf("Should not return authentication result")
	}
}

func TestAuthenticationUnknownUsername(t *testing.T) {
	username := randomString(6)
	password := randomString(6)
	if _, err := UserStore.NewUser(newUserParameters{
		Username: username,
		Email:    randomString(6),
		Password: password,
	}); err != nil {
		t.Fatalf("Error making user: %s", err.Message)
	}

	sessionKey := authenticateUser(randomString(6), randomString(6), &http.Request{RemoteAddr: randomString(6)})
	if sessionKey != nil {
		t.Fatalf("Should not return authentication result")
	}
}

func TestAuthenticationIllegalParams(t *testing.T) {
	// Empty username
	sessionKey := authenticateUser("", randomString(6), &http.Request{RemoteAddr: randomString(6)})
	if sessionKey != nil {
		t.Fatalf("Should not return authentication result")
	}

	// Empty password
	sessionKey = authenticateUser(randomString(6), "", &http.Request{RemoteAddr: randomString(6)})
	if sessionKey != nil {
		t.Fatalf("Should not return authentication result")
	}

	// Too long username
	sessionKey = authenticateUser(randomString(32), randomString(6), &http.Request{RemoteAddr: randomString(6)})
	if sessionKey != nil {
		t.Fatalf("Should not return authentication result")
	}

	// Too long password
	sessionKey = authenticateUser(randomString(6), randomString(256), &http.Request{RemoteAddr: randomString(6)})
	if sessionKey != nil {
		t.Fatalf("Should not return authentication result")
	}
}

func TestAuthenticationDisabledUser(t *testing.T) {
	username := randomString(6)
	password := randomString(6)
	user, erro := UserStore.NewUser(newUserParameters{
		Username: username,
		Email:    randomString(6),
		Password: password,
	})
	if erro != nil {
		t.Fatalf("Error making user: %s", erro.Message)
	}

	if _, err := UserStore.EditUser(user, editUserParameters{
		Email:    user.Email,
		CanLogIn: false,
	}); err != nil {
		t.Fatalf("Error updating user: %s", err.Message)
	}

	sessionKey := authenticateUser(username, password, &http.Request{RemoteAddr: randomString(6)})
	if sessionKey != nil {
		t.Fatalf("Should not return authentication result")
	}
}

func TestAuthenticationDeletedUser(t *testing.T) {
	username := randomString(6)
	password := randomString(6)
	user, erro := UserStore.NewUser(newUserParameters{
		Username: username,
		Email:    randomString(6),
		Password: password,
	})
	if erro != nil {
		t.Fatalf("Error making user: %s", erro.Message)
	}

	authenticationResult := authenticateUser(username, password, &http.Request{RemoteAddr: randomString(6)})
	if authenticationResult == nil {
		t.Fatalf("Should return a session key")
	}

	session := sessionForHTTPRequest(mockHTTPRequest("/", authenticationResult.SessionKey), false)
	if session == nil {
		t.Fatalf("Empty session for valid cookie")
	}

	UserStore.DeleteUser(user)
	session = sessionForHTTPRequest(mockHTTPRequest("/", authenticationResult.SessionKey), false)
	if session != nil {
		t.Fatalf("Should not return session for deleted user")
	}
}

func TestAuthenticationExpiredSession(t *testing.T) {
	username := randomString(6)
	password := randomString(6)
	if _, err := UserStore.NewUser(newUserParameters{
		Username: username,
		Email:    randomString(6),
		Password: password,
	}); err != nil {
		t.Fatalf("Error making user: %s", err.Message)
	}

	authenticationResult := authenticateUser(username, password, &http.Request{RemoteAddr: randomString(6)})
	if authenticationResult == nil {
		t.Fatalf("Should return a session key")
	}

	// Expire the session
	SessionStore.l.Lock()
	s := SessionStore.m[authenticationResult.SessionKey]
	s.Expires = time.Now().AddDate(0, 0, -1)
	SessionStore.m[s.Key] = s
	SessionStore.l.Unlock()

	session := sessionForHTTPRequest(mockHTTPRequest("/", authenticationResult.SessionKey), false)
	if session != nil {
		t.Fatalf("Should not return an expired session")
	}
}

func TestAuthenticationUnknownCookie(t *testing.T) {
	session := sessionForHTTPRequest(mockHTTPRequest("/", randomString(6)), false)
	if session != nil {
		t.Fatalf("Should not return an unknown session")
	}
}

func TestAuthenticationNoCookie(t *testing.T) {
	session := sessionForHTTPRequest(&http.Request{
		URL: &url.URL{},
	}, false)
	if session != nil {
		t.Fatalf("Should not return an unknown session")
	}
}

func TestAuthenticationPartialSession(t *testing.T) {
	username := randomString(6)
	password := randomString(6)
	if _, err := UserStore.NewUser(newUserParameters{
		Username:           username,
		Email:              randomString(6),
		Password:           password,
		MustChangePassword: true,
	}); err != nil {
		t.Fatalf("Error making user: %s", err.Message)
	}

	authenticationResult := authenticateUser(username, password, &http.Request{RemoteAddr: randomString(6)})
	if authenticationResult == nil {
		t.Fatalf("Should return a session key")
	}

	if !authenticationResult.MustChangePassword {
		t.Fatalf("Must return partial session")
	}

	session := sessionForHTTPRequest(mockHTTPRequest("/", authenticationResult.SessionKey), true)
	if session == nil {
		t.Fatalf("Should return a partial session")
	}

	session = sessionForHTTPRequest(mockHTTPRequest("/", authenticationResult.SessionKey), false)
	if session != nil {
		t.Fatalf("Should not return a partial session")
	}
}

func TestAuthenticationAPIKey(t *testing.T) {
	username := randomString(6)
	password := randomString(6)
	if _, err := UserStore.NewUser(newUserParameters{
		Username: username,
		Email:    randomString(6),
		Password: password,
	}); err != nil {
		t.Fatalf("Error making user: %s", err.Message)
	}

	a, err := UserStore.ResetAPIKey(username)
	if err != nil {
		t.Fatalf("Error resetting API key: %s", err.Message)
	}
	if a == nil {
		t.Fatalf("must return API key")
	}
	apiKey := *a

	req, erro := http.NewRequest("GET", "/api/blah", nil)
	if erro != nil {
		panic("invalid request")
	}
	req.Header.Add(ottoAPIUsernameHeader, username)
	req.Header.Add(ottoAPIKeyheader, apiKey)

	if sessionForHTTPRequest(req, false) == nil {
		t.Fatalf("Should return a session for an API key")
	}
}

func TestAuthenticationMissingAPIHeaders(t *testing.T) {
	req, erro := http.NewRequest("GET", "/api/blah", nil)
	if erro != nil {
		panic("invalid request")
	}

	if sessionForHTTPRequest(req, false) != nil {
		t.Fatalf("Should not return a session with missing API headers")
	}
}

func TestAuthenticationIncorrectAPIUsername(t *testing.T) {
	req, erro := http.NewRequest("GET", "/api/blah", nil)
	if erro != nil {
		panic("invalid request")
	}
	req.Header.Add(ottoAPIUsernameHeader, randomString(6))
	req.Header.Add(ottoAPIKeyheader, randomString(6))

	if sessionForHTTPRequest(req, false) != nil {
		t.Fatalf("Should not return a session for API request with an unknown username")
	}
}

func TestAuthenticationIncorrectAPIKey(t *testing.T) {
	username := randomString(6)
	password := randomString(6)
	if _, err := UserStore.NewUser(newUserParameters{
		Username: username,
		Email:    randomString(6),
		Password: password,
	}); err != nil {
		t.Fatalf("Error making user: %s", err.Message)
	}

	a, err := UserStore.ResetAPIKey(username)
	if err != nil {
		t.Fatalf("Error resetting API key: %s", err.Message)
	}
	if a == nil {
		t.Fatalf("must return API key")
	}

	req, erro := http.NewRequest("GET", "/api/blah", nil)
	if erro != nil {
		panic("invalid request")
	}
	req.Header.Add(ottoAPIUsernameHeader, username)
	req.Header.Add(ottoAPIKeyheader, randomString(6))

	if sessionForHTTPRequest(req, false) != nil {
		t.Fatalf("Should not return a session for an API request with an unknown API key")
	}
}
