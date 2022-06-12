package server

// This file is was generated automatically by Codegen v1.10.0
// Do not make changes to this file as they will be lost

import (
	"bytes"
	"encoding/gob"
	"sync"

	"github.com/ecnepsnai/store"
)

type cbgenStateObject struct {
	store *store.Store
	locks map[string]*sync.RWMutex
}

// State the global state object
var State *cbgenStateObject

// stateSetup load the saved state
func stateSetup() {
	s, err := store.New(Directories.Data, "state", nil)
	if err != nil {
		log.Fatal("Error opening state store: %s", err.Error())
	}
	state := cbgenStateObject{
		store: s,
		locks: map[string]*sync.RWMutex{
			"TableVersion": {},
		},
	}
	State = &state
}

// Close closes the state session
func (s *cbgenStateObject) Close() {
	s.store.Close()
}

// GetAll will return a map of all current state values
func (s *cbgenStateObject) GetAll() map[string]interface{} {
	return map[string]interface{}{
		"TableVersion": s.GetTableVersion(),
	}
}

// GetTableVersion get the TableVersion value
func (s *cbgenStateObject) GetTableVersion() int {
	s.locks["TableVersion"].RLock()
	defer s.locks["TableVersion"].RUnlock()

	d := s.store.Get("TableVersion")
	if d == nil {
		return 0
	}
	v, err := cbgenStateDecodeint(d)
	if err != nil {
		log.Error("Error decoding %s value for %s: %s", "int", "TableVersion", err.Error())
		return 0
	}
	log.Debug("state: key='state.TableVersion' current='%v'", v)
	return *v
}

// SetTableVersion set the TableVersion value
func (s *cbgenStateObject) SetTableVersion(value int) {
	s.locks["TableVersion"].Lock()
	defer s.locks["TableVersion"].Unlock()

	b, err := cbgenStateEncodeint(value)
	if err != nil {
		log.Error("Error encoding %s value for %s: %s", "int", "TableVersion", err.Error())
		return
	}
	log.Debug("state: key='state.TableVersion' new='%v'", value)
	s.store.Write("TableVersion", b)
}

// DefaultTableVersion get the default value for TableVersion
func (s *cbgenStateObject) DefaultTableVersion() int {
	return 0
}

// ResetTableVersion resets TableVersion to the default value
func (s *cbgenStateObject) ResetTableVersion() {
	s.SetTableVersion(s.DefaultTableVersion())
}

func cbgenStateEncodeint(o int) ([]byte, error) {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(o)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func cbgenStateDecodeint(data []byte) (*int, error) {
	w := new(int)
	reader := bytes.NewReader(data)
	dec := gob.NewDecoder(reader)
	if err := dec.Decode(&w); err != nil {
		return nil, err
	}
	return w, nil
}
