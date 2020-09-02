package server

import "github.com/ecnepsnai/otto/server/environ"

func staticEnvironment() []environ.Variable {
	return []environ.Variable{
		environ.New("OTTO_VERSION", ServerVersion),
		environ.New("OTTO_URL", Options.General.ServerURL),
	}
}
