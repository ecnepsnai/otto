package server

// This file is was generated automatically by Codegen v1.8.1
// Do not make changes to this file as they will be lost

import (
	"encoding/gob"

	"github.com/ecnepsnai/otto/server/environ"
)

func cbgenGobRegisterType(o interface{}) {
	defer cbgenGobPanicRecovery()
	gob.Register(o)
}

func cbgenGobPanicRecovery() {
	recover()
}

// gobSetup register gob types
func gobSetup() {

	cbgenGobRegisterType(environ.Variable{})

}
