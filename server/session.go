package server

import (
	"time"

	"github.com/ecnepsnai/security"
)

// Session describes a user session
type Session struct {
	// ID the session ID
	ID string `ds:"primary"`
	// Secret the salt used when hashing the users email for the session auth cookie value
	Secret string `json:"-" structs:"-" ds:"unique"`
	// Username the username of the user for this session
	Username string `ds:"index"`
	// Expires when this session expires (unix timestamp)
	Expires int64
}

// NewSessionForUser start a new session for the given user
func (s sessionStoreObject) NewSessionForUser(user *User) (Session, string, *Error) {
	session := Session{
		ID:       NewUUID(),
		Secret:   GenerateSessionSecret(),
		Username: user.Username,
		Expires:  time.Now().AddDate(0, 0, 1).Unix(),
	}
	sessionHash := security.HashString(session.Secret + session.Username)
	sessionCookie := session.ID + "$" + sessionHash
	log.Info("Started new session for user: '%s' with session ID: '%s'", user.Username, session.ID)

	if err := s.Table.Add(session); err != nil {
		log.Error("Error adding new session for user '%s': %s", user.Username, err.Error())
		return session, sessionCookie, ErrorFrom(err)
	}

	return session, sessionCookie, nil
}

// SessionWithID locate a session with the given ID
func (s sessionStoreObject) SessionWithID(ID string) (*Session, *Error) {
	object, err := s.Table.Get(ID)
	if err != nil {
		log.Error("Error getting session '%s': %s", ID, err.Error())
		return nil, ErrorFrom(err)
	}
	if object == nil {
		log.Warn("Session with ID: '%s' not found", ID)
		return nil, nil
	}

	session := object.(Session)
	return &session, nil
}

// DeleteSession delete a session with the given ID
func (s sessionStoreObject) DeleteSession(session *Session) {
	log.Info("Ending session for user: '%s' with session ID: '%s'", session.Username, session.ID)
	s.Table.Delete(*session)
}

// SessionForUser locate all sessions for the given user
func (s sessionStoreObject) SessionForUser(username string) ([]Session, *Error) {
	objects, err := s.Table.GetIndex("Username", username, nil)
	if err != nil {
		log.Error("Error getting sessions for user '%s': %s", username, err.Error())
		return nil, ErrorFrom(err)
	}
	sessions := make([]Session, len(objects))
	for i, object := range objects {
		sessions[i] = object.(Session)
	}
	return sessions, nil
}

// EndAllForUser end all sessions for user
func (s sessionStoreObject) EndAllForUser(user User) {
	s.Table.DeleteAllIndex("Username", user.Username)
}

// SaveSession save the given session (new or current)
func (s sessionStoreObject) SaveSession(session *Session) *Error {
	if err := s.Table.Update(*session); err != nil {
		log.Error("Error updating session: %s", err.Error())
		return ErrorFrom(err)
	}
	return nil
}

// User get the user object for this session
func (s Session) User() *User {
	user, _ := UserStore.UserWithUsername(s.Username)
	return user
}

func (s sessionStoreObject) CleanupSessions() *Error {
	objects, err := s.Table.GetAll(nil)
	if err != nil {
		log.Error("Error getting all sessions: %s", err.Error())
		return ErrorFrom(err)
	}
	count := len(objects)
	if count == 0 {
		log.Debug("No sessions in table")
		return nil
	}
	sessions := make([]Session, count)
	for i, object := range objects {
		session, k := object.(Session)
		if !k {
			log.Error("Object is not of type Session")
			return ErrorServer("invalid type")
		}
		sessions[i] = session
	}

	sessionsCleared := 0
	for _, session := range sessions {
		if time.Now().Unix() >= session.Expires {
			s.DeleteSession(&session)
			sessionsCleared++
		}
	}
	log.Info("Removed %d expired sessions", sessionsCleared)
	return nil
}
