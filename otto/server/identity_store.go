package server

import "github.com/ecnepsnai/otto"

// Get will save the given identity for the host ID
func (s *identityStoreObject) Get(hostID string) otto.Identity {
	return s.Store.Get(hostID)
}

// Set will save the given identity for the host ID
func (s *identityStoreObject) Set(hostID string, identity otto.Identity) {
	s.Store.Write(hostID, identity)
}

// Delete will remove any saved identity for the host ID
func (s *identityStoreObject) Delete(hostID string) {
	s.Store.Delete(hostID)
}
