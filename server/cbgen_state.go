package server

// This file is was generated automatically by Codegen v1.6.0
// Do not make changes to this file as they will be lost

import (
	"bytes"
	"encoding/gob"

	"github.com/ecnepsnai/store"
)

type stateObject struct {
	store *store.Store
}

// State the global state object
var State *stateObject

// stateSetup load the saved state
func stateSetup() {
	s, err := store.New(Directories.Data, "state", nil)
	if err != nil {
		log.Fatal("Error opening state store: %s", err.Error())
	}
	state := stateObject{
		store: s,
	}
	State = &state
}

// Close closes the state session
func (s *stateObject) Close() {
	s.store.Close()
}

// GetTableVersion get the TableVersion value
func (s *stateObject) GetTableVersion() int {
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
func (s *stateObject) SetTableVersion(value int) {
	b, err := cbgenStateEncodeint(value)
	if err != nil {
		log.Error("Error encoding %s value for %s: %s", "int", "TableVersion", err.Error())
		return
	}
	log.Debug("state: key='state.TableVersion' new='%v'", value)
	s.store.Write("TableVersion", b)
}

// DefaultTableVersion get the default value for TableVersion
func (s *stateObject) DefaultTableVersion() int {
	return 0
}

// ResetTableVersion resets TableVersion to the default value
func (s *stateObject) ResetTableVersion() {
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
