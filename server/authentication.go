package server

import (
	"net/http"
	"time"
)

const (
	ottoSessionCookie     = "otto-session"
	ottoAPIUsernameHeader = "X-OTTO-USERNAME"
	ottoAPIKeyheader      = "X-OTTO-API-KEY"
)

// AuthenticationResult describes an authentication result
type AuthenticationResult struct {
	SessionKey         string
	MustChangePassword bool
}

func sessionForHTTPRequest(r *http.Request, allowPartial bool) *Session {
	sessionCookie, _ := r.Cookie(ottoSessionCookie)
	if sessionCookie == nil {
		if len(r.URL.Path) < 4 || r.URL.Path[0:4] != "/api" {
			return nil
		}

		apiUsername := r.Header.Get(ottoAPIUsernameHeader)
		if apiUsername == "" {
			return nil
		}
		apiKey := r.Header.Get(ottoAPIKeyheader)
		if apiKey == "" {
			return nil
		}

		user := UserStore.UserWithUsername(apiUsername)
		if user == nil {
			log.Warn("API request for non-existant user: username='%s'", apiUsername)
			return nil
		}

		if !user.APIKey.Compare([]byte(apiKey)) {
			log.Warn("API request with incorrect API key: username='%s'", apiUsername)
			return nil
		}

		session := Session{
			Key:      generateSessionSecret(),
			ShortID:  newPlainID(),
			Username: user.Username,
		}
		log.Info("HTTP request: api_session method='%s' url='%s' username='%s'", r.Method, r.URL.String(), session.Username)
		return &session
	}

	sessionKey := sessionCookie.Value
	session := SessionStore.SessionWithID(sessionKey)
	if session == nil {
		log.Debug("No session found: session_key='%s'", sessionKey)
		return nil
	}

	if time.Since(session.Expires) > 0 {
		log.Warn("Session expired: session_id='%s' username='%s' expired_on='%s'", session.ShortID, session.Username, session.Expires)
		return nil
	}

	if user := UserStore.UserWithUsername(session.Username); user == nil {
		log.Warn("Session for nonexistant user: session_id='%s' username='%s'", session.ShortID, session.Username)
		return nil
	}

	if session.Partial && !allowPartial {
		return nil
	}

	log.Info("HTTP request: session_id='%s' method='%s' url='%s' username='%s'", session.ShortID, r.Method, r.URL.String(), session.Username)
	updatedSession := SessionStore.UpdateSessionExpiry(session.Key)
	return &updatedSession
}

func authenticateUser(username, password string, req *http.Request) *AuthenticationResult {
	usernameLen := len(username)
	passwordLen := len(password)
	if usernameLen == 0 || usernameLen > 32 || passwordLen == 0 || passwordLen > 256 {
		log.Debug("Reject login with illegal parameters")
		return nil
	}

	user := UserStore.UserWithUsername(username)
	if user == nil {
		log.Warn("Reject login for unknown user: username='%s'", username)
		return nil
	}

	if !user.CanLogIn {
		log.Warn("User prohibited from accessing system: username='%s'", user.Username)
		return nil
	}

	if !user.PasswordHash.Compare([]byte(password)) {
		EventStore.UserIncorrectPassword(username, req.RemoteAddr)
		log.Warn("Reject login with incorrect password: username='%s'", user.Username)
		return nil
	}

	if upgradedPassword := user.PasswordHash.Upgrade([]byte(password)); upgradedPassword != nil {
		user.PasswordHash = *upgradedPassword
		if err := UserStore.Table.Update(*user); err != nil {
			log.Error("Error upgrading user password: username='%s' error='%s'", user.Username, err.Error())
		} else {
			log.Info("Upgraded user password: username='%s'", user.Username)
		}
	}
	password = ""

	session := SessionStore.NewSessionForUser(user)
	EventStore.UserLoggedIn(username, req.RemoteAddr)
	return &AuthenticationResult{
		SessionKey:         session.Key,
		MustChangePassword: user.MustChangePassword,
	}
}
