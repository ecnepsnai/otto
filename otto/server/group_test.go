package server

import (
	"testing"
)

func TestAddGetGroup(t *testing.T) {
	name := randomString(6)

	group, err := GroupStore.NewGroup(newGroupParameters{
		Name: name,
	})
	if err != nil {
		t.Fatalf("Error making new group: %s", err.Message)
	}
	if group == nil {
		t.Fatalf("Should return a group")
	}

	if GroupStore.GroupWithID(group.ID) == nil {
		t.Fatalf("Should return a group with an ID")
	}
	if GroupStore.GroupWithName(name) == nil {
		t.Fatalf("Should return a group with an Name")
	}
}

func TestEditGroup(t *testing.T) {
	name := randomString(6)

	group, err := GroupStore.NewGroup(newGroupParameters{
		Name: name,
	})
	if err != nil {
		t.Fatalf("Error making new group: %s", err.Message)
	}
	if group == nil {
		t.Fatalf("Should return a group")
	}

	group, err = GroupStore.EditGroup(group, editGroupParameters{
		Name: randomString(6),
	})
	if err != nil {
		t.Fatalf("Error editing group: %s", err.Message)
	}
	if group == nil {
		t.Fatalf("Should return a group")
	}

	group = GroupStore.GroupWithID(group.ID)
	if group.Name == name {
		t.Fatalf("Should change name")
	}
}

func TestDeleteGroup(t *testing.T) {
	// Make a dummy group since there must be at least one group
	_, err := GroupStore.NewGroup(newGroupParameters{
		Name: randomString(6),
	})
	if err != nil {
		t.Fatalf("Error making new group: %s", err.Message)
	}

	name := randomString(6)
	group, err := GroupStore.NewGroup(newGroupParameters{
		Name: name,
	})
	if err != nil {
		t.Fatalf("Error making new group: %s", err.Message)
	}
	if group == nil {
		t.Fatalf("Should return a group")
	}

	if err := GroupStore.DeleteGroup(group); err != nil {
		t.Fatalf("Error deleting group: %s", err.Message)
	}
	if GroupStore.GroupWithID(group.ID) != nil {
		t.Fatalf("Should not return a group with an ID")
	}
	if GroupStore.GroupWithName(name) != nil {
		t.Fatalf("Should not return a group with an Name")
	}
}

func TestAddDuplicateGroup(t *testing.T) {
	name := randomString(6)

	group, err := GroupStore.NewGroup(newGroupParameters{
		Name: name,
	})
	if err != nil {
		t.Fatalf("Error making new group: %s", err.Message)
	}
	if group == nil {
		t.Fatalf("Should return a group")
	}

	_, err = GroupStore.NewGroup(newGroupParameters{
		Name: name,
	})
	if err == nil {
		t.Fatalf("Should return error")
	}
}

func TestRenameDuplicateGroup(t *testing.T) {
	name := randomString(6)

	groupA, err := GroupStore.NewGroup(newGroupParameters{
		Name: name,
	})
	if err != nil {
		t.Fatalf("Error making new group: %s", err.Message)
	}
	if groupA == nil {
		t.Fatalf("Should return a group")
	}

	groupB, err := GroupStore.NewGroup(newGroupParameters{
		Name: randomString(6),
	})
	if err != nil {
		t.Fatalf("Error making new group: %s", err.Message)
	}
	if groupB == nil {
		t.Fatalf("Should return a group")
	}

	_, err = GroupStore.EditGroup(groupB, editGroupParameters{
		Name: name,
	})
	if err == nil {
		t.Fatalf("Should return error")
	}
}
