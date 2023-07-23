package server

import (
	"testing"
)

func TestAddGetHost(t *testing.T) {
	name := randomString(6)
	address := randomString(5)

	host, err := HostStore.NewHost(newHostParameters{
		Name:    name,
		Address: address,
		Port:    12444,
	})
	if err != nil {
		t.Fatalf("Error making new host: %s", err.Message)
	}
	if host == nil {
		t.Fatalf("Should return a host")
	}

	if HostStore.HostWithID(host.ID) == nil {
		t.Fatalf("Should return a host with an ID")
	}
	if HostStore.HostWithName(name) == nil {
		t.Fatalf("Should return a host with an Name")
	}
	if HostStore.HostWithAddress(address) == nil {
		t.Fatalf("Should return a host with an Address")
	}
}

func TestEditHost(t *testing.T) {
	name := randomString(6)
	address := randomString(5)

	host, err := HostStore.NewHost(newHostParameters{
		Name:    name,
		Address: address,
		Port:    12444,
	})
	if err != nil {
		t.Fatalf("Error making new host: %s", err.Message)
	}
	if host == nil {
		t.Fatalf("Should return a host")
	}

	host, err = HostStore.EditHost(host, editHostParameters{
		Name:    randomString(6),
		Address: address,
		Port:    12444,
		Enabled: true,
	})
	if err != nil {
		t.Fatalf("Error editing host: %s", err.Message)
	}
	if host == nil {
		t.Fatalf("Should return a host")
	}

	host = HostStore.HostWithID(host.ID)
	if host.Name == name {
		t.Fatalf("Should change name")
	}
}

func TestDeleteHost(t *testing.T) {
	name := randomString(6)
	address := randomString(5)

	host, err := HostStore.NewHost(newHostParameters{
		Name:    name,
		Address: address,
		Port:    12444,
	})
	if err != nil {
		t.Fatalf("Error making new host: %s", err.Message)
	}
	if host == nil {
		t.Fatalf("Should return a host")
	}

	if err := HostStore.DeleteHost(host); err != nil {
		t.Fatalf("Error deleting host: %s", err.Message)
	}
	if HostStore.HostWithID(host.ID) != nil {
		t.Fatalf("Should not return a host with an ID")
	}
	if HostStore.HostWithName(name) != nil {
		t.Fatalf("Should not return a host with an Name")
	}
}

func TestAddDuplicateHost(t *testing.T) {
	name := randomString(6)
	address := randomString(5)

	host, err := HostStore.NewHost(newHostParameters{
		Name:    name,
		Address: address,
		Port:    12444,
	})
	if err != nil {
		t.Fatalf("Error making new host: %s", err.Message)
	}
	if host == nil {
		t.Fatalf("Should return a host")
	}

	_, err = HostStore.NewHost(newHostParameters{
		Name:    name,
		Address: randomString(6),
		Port:    12444,
	})
	if err == nil {
		t.Fatalf("Should return error")
	}

	_, err = HostStore.NewHost(newHostParameters{
		Name:    randomString(6),
		Address: address,
		Port:    12444,
	})
	if err == nil {
		t.Fatalf("Should return error")
	}
}

func TestRenameDuplicateHost(t *testing.T) {
	name := randomString(6)
	address := randomString(5)

	hostA, err := HostStore.NewHost(newHostParameters{
		Name:    name,
		Address: address,
		Port:    12444,
	})
	if err != nil {
		t.Fatalf("Error making new host: %s", err.Message)
	}
	if hostA == nil {
		t.Fatalf("Should return a host")
	}

	hostB, err := HostStore.NewHost(newHostParameters{
		Name:    randomString(6),
		Address: randomString(6),
		Port:    12444,
	})
	if err != nil {
		t.Fatalf("Error making new host: %s", err.Message)
	}
	if hostB == nil {
		t.Fatalf("Should return a host")
	}

	_, err = HostStore.EditHost(hostB, editHostParameters{
		Name:    name,
		Address: randomString(6),
		Port:    12444,
	})
	if err == nil {
		t.Fatalf("Should return error")
	}

	_, err = HostStore.EditHost(hostB, editHostParameters{
		Name:    randomString(6),
		Address: address,
		Port:    12444,
	})
	if err == nil {
		t.Fatalf("Should return error")
	}
}

func TestHostRemoveWithSchedule(t *testing.T) {
	script, err := ScriptStore.NewScript(newScriptParameters{
		Name:       "script",
		Executable: "a",
		Script:     "a",
		RunLevel:   ScriptRunLevelReadOnly,
	})
	if err != nil {
		t.Fatalf("Error making new script: %s", err.Message)
	}

	group, err := GroupStore.NewGroup(newGroupParameters{
		Name:      "group",
		ScriptIDs: []string{script.ID},
	})
	if err != nil {
		t.Fatalf("Error making new group: %s", err.Message)
	}

	host, err := HostStore.NewHost(newHostParameters{
		Name:     "host",
		Address:  "host",
		Port:     12444,
		GroupIDs: []string{group.ID},
	})
	if err != nil {
		t.Fatalf("Error making new host: %s", err.Message)
	}

	_, err = ScheduleStore.NewSchedule(newScheduleParameters{
		Name:     "schedule",
		Pattern:  "* * * * *",
		ScriptID: script.ID,
		Scope: ScheduleScope{
			HostIDs: []string{
				host.ID,
			},
		},
	})
	if err != nil {
		t.Fatalf("Error making new schedule: %s", err.Message)
	}

	if err := HostStore.DeleteHost(host); err == nil {
		t.Errorf("No error seen when trying to delete host associated with schedule")
	}
}
