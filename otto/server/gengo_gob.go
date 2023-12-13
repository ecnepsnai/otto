package server

// This file is was generated automatically by GenGo v1.13.0
// Do not make changes to this file as they will be lost

import (
	"encoding/gob"

	"github.com/ecnepsnai/otto/server/environ"
)

func gengoGobRegisterType(o interface{}) {
	defer gengoGobPanicRecovery()
	gob.Register(o)
}

func gengoGobPanicRecovery() {
	recover()
}

// gobSetup register gob types
func gobSetup() {

	gengoGobRegisterType(environ.Variable{})

}
