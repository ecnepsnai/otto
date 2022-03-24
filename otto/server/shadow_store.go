package server

import "github.com/ecnepsnai/secutil"

// Compare will return true if the saved hash for the given username matches the provided raw password
func (s *shadowStoreObject) Compare(username string, raw []byte) bool {
	hashedPwBytes := s.Store.Get(username)
	if len(hashedPwBytes) == 0 {
		log.PWarn("No shadow entry found for user", map[string]interface{}{
			"username": username,
		})
		UserStore.DisableUser(username)
		EventStore.UserModified(username, systemUsername)
		return false
	}

	hashedPw := secutil.HashedPassword(hashedPwBytes)
	return hashedPw.Compare(raw)
}

// Set will save the given hash for the username
func (s *shadowStoreObject) Set(username string, hash secutil.HashedPassword) {
	s.Store.Write(username, hash)
}

// Delete will remove any saved hash for the username
func (s *shadowStoreObject) Delete(username string) {
	s.Store.Delete(username)
}

// Upgrade will upgrade this password with an improved algorithm if needed
func (s *shadowStoreObject) Upgrade(username string, raw []byte) {
	hashedPwBytes := s.Store.Get(username)
	if len(hashedPwBytes) == 0 {
		log.PWarn("No shadow entry found for user", map[string]interface{}{
			"username": username,
		})
		UserStore.DisableUser(username)
		EventStore.UserModified(username, systemUsername)
		return
	}

	hashedPw := secutil.HashedPassword(hashedPwBytes)
	newHashedPw := hashedPw.Upgrade(raw)
	if newHashedPw == nil {
		return
	}

	s.Set(username, *newHashedPw)
}
