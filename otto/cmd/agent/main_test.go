package main

import (
	"os"
	"testing"

	"github.com/ecnepsnai/logtic"
)

func TestMain(m *testing.M) {
	for _, arg := range os.Args {
		if arg == "-test.v=true" {
			logtic.Log.Level = logtic.LevelDebug
			logtic.Log.Open()
			log = logtic.Log.Connect("otto")
		}
	}

	retCode := m.Run()
	os.Exit(retCode)
}
