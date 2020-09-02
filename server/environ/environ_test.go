package environ_test

import (
	"testing"

	"github.com/ecnepsnai/otto/server/environ"
)

func TestEnviron(t *testing.T) {
	vars := []environ.Variable{
		environ.New("1", "1"),
		environ.New("2", "1"),
		environ.New("3", "1"),
	}

	moreVars := []environ.Variable{
		environ.New("3", "2"),
		environ.New("4", "2"),
		environ.New("5", "2"),
	}

	finalVars := environ.Map(environ.Merge(vars, moreVars))
	if finalVars["3"] != "2" {
		t.Fatalf("Incorrect variable value after merge")
	}

	if finalVars["5"] != "2" {
		t.Fatalf("Missing variable after merge")
	}
}
