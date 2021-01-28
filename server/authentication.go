package server

import (
	"net/http"
	"strings"
	"time"

	"github.com/ecnepsnai/security"
	"github.com/ecnepsnai/web"
)

// Credentials describes credentiuals
type Credentials struct {
	Username string `limits:"32"`
	Password string `limits:"256"`
}

// AuthenticationResult describes the result for authentication
type AuthenticationResult struct {
	Session     Session
	CookieValue string
}

const (
	ottoSessionCookie = "otto-session"
)

// IsAuthenticated is there a valid user session for the given HTTP request.
// Returns a populated session object if valid, nil if invalid
func IsAuthenticated(r *http.Request) *Session {
	sessionCookie, _ := r.Cookie(ottoSessionCookie)
	return authenticateUser(sessionCookie)
}

func authenticateUser(sessionCookie *http.Cookie) *Session {
	if sessionCookie == nil {
		log.Warn("Invalid or missing otto session cookie")
		return nil
	}

	cookieComponents := strings.Split(sessionCookie.Value, "$")
	if len(cookieComponents) != 2 {
		log.Warn("Invalid otto session cookie")
		return nil
	}

	sessionID := cookieComponents[0]
	sessionHash := cookieComponents[1]

	session, err := SessionStore.SessionWithID(sessionID)
	if err != nil {
		log.Error("Error fetching session '%s': %s", sessionID, err.Message)
		return nil
	}
	if session == nil {
		log.Warn("No session with ID '%s' found", sessionID)
		return nil
	}

	if time.Now().Unix() >= session.Expires {
		log.Warn("Session expired: %d >= %d", time.Now().Unix(), session.Expires)
		return nil
	}

	trustedHash := security.HashSHA256String(session.Secret + session.Username)
	if trustedHash != sessionHash {
		log.Warn("Invalid otto session hash")
		log.Debug("'%s' != '%s'", trustedHash, sessionHash)
		return nil
	}

	if user := UserStore.UserWithUsername(session.Username); user == nil {
		log.Warn("Session for nonexistant user: session_id='%s' username='%s'", session.ID, session.Username)
		return nil
	}

	// Update expires timestamp
	session.Expires = time.Now().Unix() + 7200
	SessionStore.SaveSession(session)

	return session
}

// AuthenticateUser authenticate a user
func AuthenticateUser(credentials Credentials, req *http.Request) (*AuthenticationResult, *web.Error) {
	user := UserStore.UserWithUsername(credentials.Username)
	if user == nil {
		return nil, web.CommonErrors.Unauthorized
	}

	if !user.Enabled {
		log.Warn("Attempted login from disabled user: '%s'", user.Username)
		return nil, web.CommonErrors.Unauthorized
	}

	if !user.PasswordHash.Compare([]byte(credentials.Password)) {
		EventStore.UserIncorrectPassword(credentials.Username, req.RemoteAddr)
		log.Warn("Incorrect password provided for user: '%s'", user.Username)
		return nil, web.CommonErrors.Unauthorized
	}

	if upgradedPassword := user.PasswordHash.Upgrade([]byte(credentials.Password)); upgradedPassword != nil {
		user.PasswordHash = *upgradedPassword
		if err := UserStore.Table.Update(*user); err != nil {
			log.Error("Error upgrading user password: %s", err.Error())
		} else {
			log.Info("Upgraded password for user '%s'", user.Username)
		}
	}

	session, cookieValue, err := SessionStore.NewSessionForUser(user)
	if err != nil {
		if err.Server {
			log.Error("Error starting new session for user '%s': %s", user.Username, err.Message)
			return nil, web.CommonErrors.ServerError
		}

		return nil, web.ValidationError(err.Message)
	}

	result := AuthenticationResult{
		Session:     session,
		CookieValue: cookieValue,
	}

	EventStore.UserLoggedIn(credentials.Username, req.RemoteAddr)

	return &result, nil
}
