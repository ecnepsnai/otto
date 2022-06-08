package server

import (
	"testing"

	"github.com/ecnepsnai/otto/server/environ"
)

func TestAddGetScript(t *testing.T) {
	name := randomString(6)

	script, err := ScriptStore.NewScript(newScriptParameters{
		Name:       name,
		Executable: "/bin/bash",
		Script:     "#!/bin/bash\necho hello\n",
		Environment: []environ.Variable{
			{
				Key:   "FOO",
				Value: "BAR",
			},
		},
		RunAs: RunAs{
			UID: 0,
			GID: 0,
		},
	})
	if err != nil {
		t.Fatalf("Error making new script: %s", err.Message)
	}
	if script == nil {
		t.Fatalf("Should return a script")
	}

	if ScriptStore.ScriptWithID(script.ID) == nil {
		t.Fatalf("Should return a script with an ID")
	}
	if ScriptStore.ScriptWithName(name) == nil {
		t.Fatalf("Should return a script with an Name")
	}
}

func TestEditScript(t *testing.T) {
	name := randomString(6)

	script, err := ScriptStore.NewScript(newScriptParameters{
		Name:       name,
		Executable: "/bin/bash",
		Script:     "#!/bin/bash\necho hello\n",
		Environment: []environ.Variable{
			{
				Key:   "FOO",
				Value: "BAR",
			},
		},
		RunAs: RunAs{
			UID: 0,
			GID: 0,
		},
	})
	if err != nil {
		t.Fatalf("Error making new script: %s", err.Message)
	}
	if script == nil {
		t.Fatalf("Should return a script")
	}

	script, err = ScriptStore.EditScript(script, editScriptParameters{
		Name:       randomString(6),
		Executable: "/bin/bash",
		Script:     "#!/bin/bash\necho hello\n",
		Environment: []environ.Variable{
			{
				Key:   "FOO",
				Value: "BAR",
			},
		},
		RunAs: RunAs{
			UID: 0,
			GID: 0,
		},
	})
	if err != nil {
		t.Fatalf("Error editing script: %s", err.Message)
	}
	if script == nil {
		t.Fatalf("Should return a script")
	}

	script = ScriptStore.ScriptWithID(script.ID)
	if script.Name == name {
		t.Fatalf("Should change name")
	}
}

func TestDeleteScript(t *testing.T) {
	name := randomString(6)

	script, err := ScriptStore.NewScript(newScriptParameters{
		Name:       name,
		Executable: "/bin/bash",
		Script:     "#!/bin/bash\necho hello\n",
		Environment: []environ.Variable{
			{
				Key:   "FOO",
				Value: "BAR",
			},
		},
		RunAs: RunAs{
			UID: 0,
			GID: 0,
		},
	})
	if err != nil {
		t.Fatalf("Error making new script: %s", err.Message)
	}
	if script == nil {
		t.Fatalf("Should return a script")
	}

	if err := ScriptStore.DeleteScript(script); err != nil {
		t.Fatalf("Error deleting script: %s", err.Message)
	}
	if ScriptStore.ScriptWithID(script.ID) != nil {
		t.Fatalf("Should not return a script with an ID")
	}
	if ScriptStore.ScriptWithName(name) != nil {
		t.Fatalf("Should not return a script with an Name")
	}
}

func TestAddDuplicateScript(t *testing.T) {
	name := randomString(6)

	script, err := ScriptStore.NewScript(newScriptParameters{
		Name:       name,
		Executable: "/bin/bash",
		Script:     "#!/bin/bash\necho hello\n",
		Environment: []environ.Variable{
			{
				Key:   "FOO",
				Value: "BAR",
			},
		},
		RunAs: RunAs{
			UID: 0,
			GID: 0,
		},
	})
	if err != nil {
		t.Fatalf("Error making new script: %s", err.Message)
	}
	if script == nil {
		t.Fatalf("Should return a script")
	}

	_, err = ScriptStore.NewScript(newScriptParameters{
		Name:       name,
		Executable: "/bin/bash",
		Script:     "#!/bin/bash\necho hello\n",
		Environment: []environ.Variable{
			{
				Key:   "FOO",
				Value: "BAR",
			},
		},
		RunAs: RunAs{
			UID: 0,
			GID: 0,
		},
	})
	if err == nil {
		t.Fatalf("Should return error")
	}
}

func TestRenameDuplicateScript(t *testing.T) {
	name := randomString(6)

	scriptA, err := ScriptStore.NewScript(newScriptParameters{
		Name:       name,
		Executable: "/bin/bash",
		Script:     "#!/bin/bash\necho hello\n",
		Environment: []environ.Variable{
			{
				Key:   "FOO",
				Value: "BAR",
			},
		},
		RunAs: RunAs{
			UID: 0,
			GID: 0,
		},
	})
	if err != nil {
		t.Fatalf("Error making new script: %s", err.Message)
	}
	if scriptA == nil {
		t.Fatalf("Should return a script")
	}

	scriptB, err := ScriptStore.NewScript(newScriptParameters{
		Name:       randomString(6),
		Executable: "/bin/bash",
		Script:     "#!/bin/bash\necho hello\n",
		Environment: []environ.Variable{
			{
				Key:   "FOO",
				Value: "BAR",
			},
		},
		RunAs: RunAs{
			UID: 0,
			GID: 0,
		},
	})
	if err != nil {
		t.Fatalf("Error making new script: %s", err.Message)
	}
	if scriptB == nil {
		t.Fatalf("Should return a script")
	}

	_, err = ScriptStore.EditScript(scriptB, editScriptParameters{
		Name:       name,
		Executable: "/bin/bash",
		Script:     "#!/bin/bash\necho hello\n",
		Environment: []environ.Variable{
			{
				Key:   "FOO",
				Value: "BAR",
			},
		},
		RunAs: RunAs{
			UID: 0,
			GID: 0,
		},
	})
	if err == nil {
		t.Fatalf("Should return error")
	}
}
