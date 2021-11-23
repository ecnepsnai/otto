package server

import "github.com/ecnepsnai/otto"

// Get will save the given identity for the clientID
func (s *identityStoreObject) Get(clientID string) otto.Identity {
	return s.Store.Get(clientID)
}

// Set will save the given identity for the clientID
func (s *identityStoreObject) Set(clientID string, identity otto.Identity) {
	s.Store.Write(clientID, identity)
}

// Delete will remove any saved identity for the clientID
func (s *identityStoreObject) Delete(clientID string) {
	s.Store.Delete(clientID)
}
