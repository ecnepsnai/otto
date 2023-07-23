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

		user, ok := UserCache.ByUsername(apiUsername)
		if !ok {
			log.PWarn("API request for non-existant user", map[string]interface{}{
				"username": apiUsername,
			})
			return nil
		}

		if !ShadowStore.Compare("api_"+user.Username, []byte(apiKey)) {
			log.PWarn("API request with incorrect API key", map[string]interface{}{
				"username": apiUsername,
			})
			return nil
		}

		session := Session{
			Key:      generateSessionSecret(),
			ShortID:  newPlainID(),
			Username: user.Username,
		}
		log.PInfo("HTTP api request", map[string]interface{}{
			"method":   r.Method,
			"url":      r.URL.String(),
			"username": session.Username,
		})
		return &session
	}

	sessionKey := sessionCookie.Value
	session := SessionStore.SessionWithID(sessionKey)
	if session == nil {
		log.PDebug("No session found", map[string]interface{}{
			"session_key": sessionKey,
		})
		return nil
	}

	if time.Since(session.Expires) > 0 {
		log.PWarn("Session expired", map[string]interface{}{
			"session_id": session.ShortID,
			"username":   session.Username,
			"expired_on": session.Expires.String(),
		})
		return nil
	}

	user, ok := UserCache.ByUsername(session.Username)
	if !ok {
		log.PWarn("Session for nonexistant user", map[string]interface{}{
			"session_id": session.ShortID,
			"username":   session.Username,
		})
		return nil
	}

	if session.Partial && !allowPartial {
		return nil
	}

	log.PInfo("HTTP request", map[string]interface{}{
		"method":   r.Method,
		"url":      r.URL.String(),
		"username": user.Username,
	})
	updatedSession := SessionStore.UpdateSessionExpiry(session.Key)
	return &updatedSession
}

func authenticateUser(username string, password []byte, req *http.Request) *AuthenticationResult {
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

	if !ShadowStore.Compare(user.Username, password) {
		EventStore.UserIncorrectPassword(username, req.RemoteAddr)
		log.Warn("Reject login with incorrect password: username='%s'", user.Username)
		return nil
	}

	ShadowStore.Upgrade(user.Username, password)
	password = nil

	session := SessionStore.NewSessionForUser(user)
	EventStore.UserLoggedIn(username, req.RemoteAddr)
	return &AuthenticationResult{
		SessionKey:         session.Key,
		MustChangePassword: user.MustChangePassword,
	}
}
