package server

import (
	"sync"
	"time"
)

// Session describes a user session
type Session struct {
	Key      string
	ShortID  string
	Username string
	Partial  bool
	Expires  time.Time
}

type sessionStoreObject struct {
	m map[string]Session
	l *sync.RWMutex
}

// SessionStore describes the session store
var SessionStore = &sessionStoreObject{
	m: map[string]Session{},
	l: &sync.RWMutex{},
}

// NewSessionForUser start a new session for the given user
func (s *sessionStoreObject) NewSessionForUser(user *User) Session {
	session := Session{
		Key:      generateSessionSecret(),
		ShortID:  newPlainID(),
		Username: user.Username,
		Expires:  time.Now().AddDate(0, 0, 1),
		Partial:  user.MustChangePassword,
	}
	log.Info("Started new session: username='%s' session_id='%s'", user.Username, session.ShortID)
	s.l.Lock()
	s.m[session.Key] = session
	log.Debug("Sessions: %+v", s.m)
	s.l.Unlock()
	return session
}

// SessionWithID locate a session with the given ID
func (s *sessionStoreObject) SessionWithID(ID string) *Session {
	s.l.RLock()
	session, present := s.m[ID]
	s.l.RUnlock()
	if !present {
		return nil
	}

	return &session
}

// DeleteSession delete a session with the given ID
func (s *sessionStoreObject) DeleteSession(session *Session) {
	s.l.Lock()
	delete(s.m, session.Key)
	log.Debug("Sessions: %+v", s.m)
	s.l.Unlock()
	log.Info("Ending session: username='%s' session_id='%s'", session.Username, session.ShortID)
}

// SessionForUser locate all sessions for the given user
func (s *sessionStoreObject) SessionForUser(username string) []Session {
	s.l.RLock()
	defer s.l.RUnlock()
	sessions := []Session{}
	for _, session := range s.m {
		if session.Username == username {
			sessions = append(sessions, session)
		}
	}
	return sessions
}

// EndAllForUser end all sessions for user
func (s *sessionStoreObject) EndAllForUser(username string) {
	sessions := s.SessionForUser(username)
	s.l.Lock()
	defer s.l.Unlock()
	for _, session := range sessions {
		delete(s.m, session.Key)
		log.Debug("Sessions: %+v", s.m)
	}
}

// EndAllOtherForUser end all sessions for the user except the current session
func (s *sessionStoreObject) EndAllOtherForUser(username string, current *Session) {
	sessions := s.SessionForUser(username)
	s.l.Lock()
	defer s.l.Unlock()
	for _, session := range sessions {
		if current.Key == session.Key {
			continue
		}

		delete(s.m, session.Key)
		log.Debug("Sessions: %+v", s.m)
	}
}

func (s *sessionStoreObject) UpdateSessionExpiry(sessionKey string) Session {
	s.l.Lock()
	session := s.m[sessionKey]
	session.Expires = time.Now().Add(time.Duration(Options.Authentication.MaxAgeMinutes) * time.Minute)
	s.m[sessionKey] = session
	log.Debug("Sessions: %+v", s.m)
	s.l.Unlock()
	return session
}

func (s *sessionStoreObject) CompletePartialSession(sessionKey string) Session {
	s.l.Lock()
	session := s.m[sessionKey]
	session.Partial = false
	s.m[sessionKey] = session
	log.Debug("Sessions: %+v", s.m)
	s.l.Unlock()
	return session
}

// User get the user object for this session
func (s Session) User() *User {
	return UserStore.UserWithUsername(s.Username)
}

func (s *sessionStoreObject) CleanupSessions() *Error {
	s.l.Lock()
	defer s.l.Unlock()

	sessions := make([]Session, len(s.m))
	i := 0
	for _, session := range s.m {
		sessions[i] = session
		i++
	}

	sessionsCleared := 0
	for _, session := range sessions {
		if time.Since(session.Expires) > 0 {
			delete(s.m, session.Key)
			sessionsCleared++
		}
	}
	log.Info("Removed %d expired sessions", sessionsCleared)
	return nil
}
