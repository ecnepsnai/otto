package server

import "github.com/ecnepsnai/otto/server/environ"

func staticEnvironment() []environ.Variable {
	return []environ.Variable{
		environ.New("OTTO_SERVER_VERSION", ServerVersion),
		environ.New("OTTO_SERVER_URL", Options.General.ServerURL),
	}
}
