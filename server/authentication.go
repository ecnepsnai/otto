package server

import (
	"net/http"
	"time"
)

const (
	ottoSessionCookie = "otto-session"
)

func sessionForHTTPRequest(r *http.Request) *Session {
	sessionCookie, _ := r.Cookie(ottoSessionCookie)
	if sessionCookie == nil {
		return nil
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

	log.Info("HTTP request: session_id='%s' method='%s' url='%s' username='%s'", session.ShortID, r.Method, r.URL.String(), session.Username)
	updatedSession := SessionStore.UpdateSessionExpiry(session.Key)
	return &updatedSession
}

func authenticateUser(username, password string, req *http.Request) *string {
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

	if !user.Enabled {
		log.Warn("Reject login from disabled user: username='%s'", user.Username)
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
	return &session.Key
}
