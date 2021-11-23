package server

import (
	"testing"

	"github.com/ecnepsnai/otto/server/environ"
)

func TestAddGetScheduleGroup(t *testing.T) {
	script, err := ScriptStore.NewScript(newScriptParameters{
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
	if script == nil {
		t.Fatalf("Should return a script")
	}

	group, err := GroupStore.NewGroup(newGroupParameters{
		Name:      randomString(6),
		ScriptIDs: []string{script.ID},
	})
	if err != nil {
		t.Fatalf("Error making new group: %s", err.Message)
	}
	if group == nil {
		t.Fatalf("Should return a group")
	}

	host, err := HostStore.NewHost(newHostParameters{
		Name:     randomString(6),
		Address:  randomString(6),
		Port:     12444,
		GroupIDs: []string{group.ID},
	})
	if err != nil {
		t.Fatalf("Error making new host: %s", err.Message)
	}
	if host == nil {
		t.Fatalf("Should return a host")
	}

	name := randomString(6)
	schedule, err := ScheduleStore.NewSchedule(newScheduleParameters{
		ScriptID: script.ID,
		Name:     name,
		Scope: ScheduleScope{
			GroupIDs: []string{group.ID},
		},
		Pattern: "* * * * *",
	})
	if err != nil {
		t.Fatalf("Error making new schedule: %s", err.Message)
	}
	if schedule == nil {
		t.Fatalf("Should return a schedule")
	}

	if ScheduleStore.ScheduleWithID(schedule.ID) == nil {
		t.Fatalf("Should return a schedule with an ID")
	}
	if ScheduleStore.ScheduleWithName(name) == nil {
		t.Fatalf("Should return a schedule with an Name")
	}
}

func TestAddGetScheduleHost(t *testing.T) {
	script, err := ScriptStore.NewScript(newScriptParameters{
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
	if script == nil {
		t.Fatalf("Should return a script")
	}

	host, err := HostStore.NewHost(newHostParameters{
		Name:    randomString(6),
		Address: randomString(6),
		Port:    12444,
	})
	if err != nil {
		t.Fatalf("Error making new host: %s", err.Message)
	}
	if host == nil {
		t.Fatalf("Should return a host")
	}

	name := randomString(6)
	schedule, err := ScheduleStore.NewSchedule(newScheduleParameters{
		ScriptID: script.ID,
		Name:     name,
		Scope: ScheduleScope{
			HostIDs: []string{host.ID},
		},
		Pattern: "* * * * *",
	})
	if err != nil {
		t.Fatalf("Error making new schedule: %s", err.Message)
	}
	if schedule == nil {
		t.Fatalf("Should return a schedule")
	}

	if ScheduleStore.ScheduleWithID(schedule.ID) == nil {
		t.Fatalf("Should return a schedule with an ID")
	}
	if ScheduleStore.ScheduleWithName(name) == nil {
		t.Fatalf("Should return a schedule with an Name")
	}
}

func TestEditSchedule(t *testing.T) {
	script, err := ScriptStore.NewScript(newScriptParameters{
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
	if script == nil {
		t.Fatalf("Should return a script")
	}

	group, err := GroupStore.NewGroup(newGroupParameters{
		Name:      randomString(6),
		ScriptIDs: []string{script.ID},
	})
	if err != nil {
		t.Fatalf("Error making new group: %s", err.Message)
	}
	if group == nil {
		t.Fatalf("Should return a group")
	}

	host, err := HostStore.NewHost(newHostParameters{
		Name:     randomString(6),
		Address:  randomString(6),
		Port:     12444,
		GroupIDs: []string{group.ID},
	})
	if err != nil {
		t.Fatalf("Error making new host: %s", err.Message)
	}
	if host == nil {
		t.Fatalf("Should return a host")
	}

	name := randomString(6)
	schedule, err := ScheduleStore.NewSchedule(newScheduleParameters{
		ScriptID: script.ID,
		Name:     name,
		Scope: ScheduleScope{
			GroupIDs: []string{group.ID},
		},
		Pattern: "* * * * *",
	})
	if err != nil {
		t.Fatalf("Error making new schedule: %s", err.Message)
	}
	if schedule == nil {
		t.Fatalf("Should return a schedule")
	}

	_, err = ScheduleStore.EditSchedule(schedule, editScheduleParameters{
		Name: randomString(6),
		Scope: ScheduleScope{
			GroupIDs: []string{group.ID},
		},
		Pattern: "* * * * *",
		Enabled: true,
	})
	if err != nil {
		t.Fatalf("Error editing schedule: %s", err.Message)
	}
	if ScheduleStore.ScheduleWithID(schedule.ID).Name == name {
		t.Fatalf("Should change name")
	}

}

func TestDeleteSchedule(t *testing.T) {
	script, err := ScriptStore.NewScript(newScriptParameters{
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
	if script == nil {
		t.Fatalf("Should return a script")
	}

	group, err := GroupStore.NewGroup(newGroupParameters{
		Name:      randomString(6),
		ScriptIDs: []string{script.ID},
	})
	if err != nil {
		t.Fatalf("Error making new group: %s", err.Message)
	}
	if group == nil {
		t.Fatalf("Should return a group")
	}

	host, err := HostStore.NewHost(newHostParameters{
		Name:     randomString(6),
		Address:  randomString(6),
		Port:     12444,
		GroupIDs: []string{group.ID},
	})
	if err != nil {
		t.Fatalf("Error making new host: %s", err.Message)
	}
	if host == nil {
		t.Fatalf("Should return a host")
	}

	name := randomString(6)
	schedule, err := ScheduleStore.NewSchedule(newScheduleParameters{
		ScriptID: script.ID,
		Name:     name,
		Scope: ScheduleScope{
			GroupIDs: []string{group.ID},
		},
		Pattern: "* * * * *",
	})
	if err != nil {
		t.Fatalf("Error making new schedule: %s", err.Message)
	}
	if schedule == nil {
		t.Fatalf("Should return a schedule")
	}

	if err := ScheduleStore.DeleteSchedule(schedule); err != nil {
		t.Fatalf("Error deleting schedule: %s", err.Message)
	}
	if ScheduleStore.ScheduleWithID(schedule.ID) != nil {
		t.Fatalf("Should not return a schedule with an ID")
	}
}

func TestAddDuplicateSchedule(t *testing.T) {
	script, err := ScriptStore.NewScript(newScriptParameters{
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
	if script == nil {
		t.Fatalf("Should return a script")
	}

	group, err := GroupStore.NewGroup(newGroupParameters{
		Name:      randomString(6),
		ScriptIDs: []string{script.ID},
	})
	if err != nil {
		t.Fatalf("Error making new group: %s", err.Message)
	}
	if group == nil {
		t.Fatalf("Should return a group")
	}

	host, err := HostStore.NewHost(newHostParameters{
		Name:     randomString(6),
		Address:  randomString(6),
		Port:     12444,
		GroupIDs: []string{group.ID},
	})
	if err != nil {
		t.Fatalf("Error making new host: %s", err.Message)
	}
	if host == nil {
		t.Fatalf("Should return a host")
	}

	name := randomString(6)
	schedule, err := ScheduleStore.NewSchedule(newScheduleParameters{
		ScriptID: script.ID,
		Name:     name,
		Scope: ScheduleScope{
			GroupIDs: []string{group.ID},
		},
		Pattern: "* * * * *",
	})
	if err != nil {
		t.Fatalf("Error making new schedule: %s", err.Message)
	}
	if schedule == nil {
		t.Fatalf("Should return a schedule")
	}
	_, err = ScheduleStore.NewSchedule(newScheduleParameters{
		ScriptID: script.ID,
		Name:     name,
		Scope: ScheduleScope{
			GroupIDs: []string{group.ID},
		},
		Pattern: "* * * * *",
	})
	if err == nil {
		t.Fatalf("Should return an error")
	}
}

func TestRenameDuplicateSchedule(t *testing.T) {
	script, err := ScriptStore.NewScript(newScriptParameters{
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
	if script == nil {
		t.Fatalf("Should return a script")
	}

	group, err := GroupStore.NewGroup(newGroupParameters{
		Name:      randomString(6),
		ScriptIDs: []string{script.ID},
	})
	if err != nil {
		t.Fatalf("Error making new group: %s", err.Message)
	}
	if group == nil {
		t.Fatalf("Should return a group")
	}

	host, err := HostStore.NewHost(newHostParameters{
		Name:     randomString(6),
		Address:  randomString(6),
		Port:     12444,
		GroupIDs: []string{group.ID},
	})
	if err != nil {
		t.Fatalf("Error making new host: %s", err.Message)
	}
	if host == nil {
		t.Fatalf("Should return a host")
	}

	name := randomString(6)
	scheduleA, err := ScheduleStore.NewSchedule(newScheduleParameters{
		ScriptID: script.ID,
		Name:     name,
		Scope: ScheduleScope{
			GroupIDs: []string{group.ID},
		},
		Pattern: "* * * * *",
	})
	if err != nil {
		t.Fatalf("Error making new schedule: %s", err.Message)
	}
	if scheduleA == nil {
		t.Fatalf("Should return a schedule")
	}

	scheduleB, err := ScheduleStore.NewSchedule(newScheduleParameters{
		ScriptID: script.ID,
		Name:     randomString(6),
		Scope: ScheduleScope{
			GroupIDs: []string{group.ID},
		},
		Pattern: "* * * * *",
	})
	if err != nil {
		t.Fatalf("Error making new schedule: %s", err.Message)
	}
	if scheduleB == nil {
		t.Fatalf("Should return a schedule")
	}

	_, err = ScheduleStore.EditSchedule(scheduleB, editScheduleParameters{
		Name: name,
		Scope: ScheduleScope{
			GroupIDs: []string{group.ID},
		},
		Pattern: "* * * * *",
	})
	if err == nil {
		t.Fatalf("Should return an error")
	}
}
