package server

import (
	"bytes"

	"github.com/ecnepsnai/otto/shared/otto"
)

// Get will save the given identity for the host ID
func (s *identityStoreObject) Get(hostID string) (*otto.Identity, error) {
	return otto.ParseIdentity(s.Store.Get(hostID))
}

// Set will save the given identity for the host ID
func (s *identityStoreObject) Set(hostID string, identity *otto.Identity) {
	buf := &bytes.Buffer{}
	identity.Write(buf)
	s.Store.Write(hostID, buf.Bytes())
}

// Delete will remove any saved identity for the host ID
func (s *identityStoreObject) Delete(hostID string) {
	s.Store.Delete(hostID)
}
