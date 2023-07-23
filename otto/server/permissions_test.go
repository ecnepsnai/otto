package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/ecnepsnai/ds"
	"github.com/ecnepsnai/otto/server/environ"
	"github.com/ecnepsnai/web"
)

func TestPermissionsCanModifyHosts(t *testing.T) {
	host, err := HostStore.NewHost(newHostParameters{
		Name:    randomString(3),
		Address: randomString(3),
		Port:    12444,
		Environment: []environ.Variable{
			{
				Key:    "key",
				Value:  "value",
				Secret: true,
			},
		},
	})
	if err != nil {
		panic(err)
	}

	user, err := UserStore.NewUser(newUserParameters{
		Username:    randomString(3),
		Password:    randomString(6),
		Permissions: UserPermissionsMin(),
	})
	if err != nil {
		panic(err)
	}

	session := SessionStore.NewSessionForUser(user)
	h := handle{}

	// Ensure user cannot create new host
	data, _, werr := h.HostNew(web.MockRequest(web.MockRequestParameters{UserData: &session}))
	if werr == nil {
		t.Fatalf("No error seen when one expected")
	}
	if data != nil {
		t.Fatalf("Data returned when none expected")
	}

	// Ensure user can read existing host
	data, _, werr = h.HostGet(web.MockRequest(web.MockRequestParameters{UserData: &session, Parameters: map[string]string{"id": host.ID}}))
	if werr != nil {
		t.Fatalf("Unexpected error: %s", werr.Message)
	}
	if data == nil {
		t.Fatalf("No data returned when some expected")
	}
	// ...but they cannot see hidden environment variables
	if h, ok := data.(Host); ok {
		if h.Environment[0].Value != "" {
			t.Fatalf("Secret environment variable not hidden")
		}
	}

	_, err = UserStore.EditUser(user, editUserParameters{
		Permissions: UserPermissions{
			CanModifyHosts: true,
		},
	})
	if err != nil {
		panic(err)
	}

	// Ensure user can edit host
	data, _, werr = h.HostEdit(web.MockRequest(web.MockRequestParameters{UserData: &session, Parameters: map[string]string{"id": host.ID}, JSONBody: *host}))
	if werr != nil {
		t.Fatalf("Unexpected error: %s", werr.Message)
	}
	if data == nil {
		t.Fatalf("No data returned when some expected")
	}
}

func TestPermissionsCanModifyGroups(t *testing.T) {
	group, err := GroupStore.NewGroup(newGroupParameters{
		Name: randomString(3),
		Environment: []environ.Variable{
			{
				Key:    "key",
				Value:  "value",
				Secret: true,
			},
		},
	})
	if err != nil {
		panic(err)
	}

	user, err := UserStore.NewUser(newUserParameters{
		Username:    randomString(3),
		Password:    randomString(6),
		Permissions: UserPermissionsMin(),
	})
	if err != nil {
		panic(err)
	}

	session := SessionStore.NewSessionForUser(user)
	h := handle{}

	// Ensure user cannot create new group
	data, _, werr := h.GroupNew(web.MockRequest(web.MockRequestParameters{UserData: &session}))
	if werr == nil {
		t.Fatalf("No error seen when one expected")
	}
	if data != nil {
		t.Fatalf("Data returned when none expected")
	}

	// Ensure user can read existing group
	data, _, werr = h.GroupGet(web.MockRequest(web.MockRequestParameters{UserData: &session, Parameters: map[string]string{"id": group.ID}}))
	if werr != nil {
		t.Fatalf("Unexpected error: %s", werr.Message)
	}
	if data == nil {
		t.Fatalf("No data returned when some expected")
	}
	// ...but they cannot see hidden environment variables
	if h, ok := data.(Group); ok {
		if h.Environment[0].Value != "" {
			t.Fatalf("Secret environment variable not hidden")
		}
	}

	_, err = UserStore.EditUser(user, editUserParameters{
		Permissions: UserPermissions{
			CanModifyGroups: true,
		},
	})
	if err != nil {
		panic(err)
	}

	// Ensure user can edit group
	data, _, werr = h.GroupEdit(web.MockRequest(web.MockRequestParameters{UserData: &session, Parameters: map[string]string{"id": group.ID}, JSONBody: *group}))
	if werr != nil {
		t.Fatalf("Unexpected error: %s", werr.Message)
	}
	if data == nil {
		t.Fatalf("No data returned when some expected")
	}
}

func TestPermissionsCanModifyScripts(t *testing.T) {
	script, err := ScriptStore.NewScript(newScriptParameters{
		Name:       randomString(3),
		Executable: randomString(3),
		Script:     randomString(3),
		Environment: []environ.Variable{
			{
				Key:    "key",
				Value:  "value",
				Secret: true,
			},
		},
		RunLevel: ScriptRunLevelReadOnly,
	})
	if err != nil {
		panic(err)
	}

	user, err := UserStore.NewUser(newUserParameters{
		Username:    randomString(3),
		Password:    randomString(6),
		Permissions: UserPermissionsMin(),
	})
	if err != nil {
		panic(err)
	}

	session := SessionStore.NewSessionForUser(user)
	h := handle{}

	// Ensure user cannot create new script
	data, _, werr := h.ScriptNew(web.MockRequest(web.MockRequestParameters{UserData: &session}))
	if werr == nil {
		t.Fatalf("No error seen when one expected")
	}
	if data != nil {
		t.Fatalf("Data returned when none expected")
	}

	// Ensure user can read existing script
	data, _, werr = h.ScriptGet(web.MockRequest(web.MockRequestParameters{UserData: &session, Parameters: map[string]string{"id": script.ID}}))
	if werr != nil {
		t.Fatalf("Unexpected error: %s", werr.Message)
	}
	if data == nil {
		t.Fatalf("No data returned when some expected")
	}
	// ...but they cannot see hidden environment variables
	if h, ok := data.(Script); ok {
		if h.Environment[0].Value != "" {
			t.Fatalf("Secret environment variable not hidden")
		}
	}

	_, err = UserStore.EditUser(user, editUserParameters{
		Permissions: UserPermissions{
			CanModifyScripts: true,
		},
	})
	if err != nil {
		panic(err)
	}

	// Ensure user can edit script
	data, _, werr = h.ScriptEdit(web.MockRequest(web.MockRequestParameters{UserData: &session, Parameters: map[string]string{"id": script.ID}, JSONBody: *script}))
	if werr != nil {
		t.Fatalf("Unexpected error: %s", werr.Message)
	}
	if data == nil {
		t.Fatalf("No data returned when some expected")
	}
}

func TestPermissionsCanModifySchedules(t *testing.T) {
	script, err := ScriptStore.NewScript(newScriptParameters{
		Name:       randomString(3),
		Executable: randomString(3),
		Script:     randomString(3),
		Environment: []environ.Variable{
			{
				Key:    "key",
				Value:  "value",
				Secret: true,
			},
		},
		RunLevel: ScriptRunLevelReadOnly,
	})
	if err != nil {
		panic(err)
	}
	group, err := GroupStore.NewGroup(newGroupParameters{
		Name:      randomString(3),
		ScriptIDs: []string{script.ID},
		Environment: []environ.Variable{
			{
				Key:    "key",
				Value:  "value",
				Secret: true,
			},
		},
	})
	if err != nil {
		panic(err)
	}

	schedule, err := ScheduleStore.NewSchedule(newScheduleParameters{
		Name:     randomString(3),
		ScriptID: script.ID,
		Pattern:  "* * * * *",
		Scope: ScheduleScope{
			GroupIDs: []string{group.ID},
		},
	})
	if err != nil {
		panic(err)
	}

	user, err := UserStore.NewUser(newUserParameters{
		Username:    randomString(3),
		Password:    randomString(6),
		Permissions: UserPermissionsMin(),
	})
	if err != nil {
		panic(err)
	}

	session := SessionStore.NewSessionForUser(user)
	h := handle{}

	// Ensure user cannot create new schedule
	data, _, werr := h.ScheduleNew(web.MockRequest(web.MockRequestParameters{UserData: &session}))
	if werr == nil {
		t.Fatalf("No error seen when one expected")
	}
	if data != nil {
		t.Fatalf("Data returned when none expected")
	}

	// Ensure user can read existing schedule
	data, _, werr = h.ScheduleGet(web.MockRequest(web.MockRequestParameters{UserData: &session, Parameters: map[string]string{"id": schedule.ID}}))
	if werr != nil {
		t.Fatalf("Unexpected error: %s", werr.Message)
	}
	if data == nil {
		t.Fatalf("No data returned when some expected")
	}

	_, err = UserStore.EditUser(user, editUserParameters{
		Permissions: UserPermissions{
			CanModifySchedules: true,
		},
	})
	if err != nil {
		panic(err)
	}

	// Ensure user can edit schedule
	data, _, werr = h.ScheduleEdit(web.MockRequest(web.MockRequestParameters{UserData: &session, Parameters: map[string]string{"id": schedule.ID}, JSONBody: *schedule}))
	if werr != nil {
		t.Fatalf("Unexpected error: %s", werr.Message)
	}
	if data == nil {
		t.Fatalf("No data returned when some expected")
	}
}

func TestPermissionsCanAccessAuditLog(t *testing.T) {
	user, err := UserStore.NewUser(newUserParameters{
		Username:    randomString(3),
		Password:    randomString(6),
		Permissions: UserPermissionsMin(),
	})
	if err != nil {
		panic(err)
	}

	session := SessionStore.NewSessionForUser(user)
	h := handle{}
	u, _ := url.Parse("https://example.localhost/blah?c=5")

	data, _, werr := h.EventsGet(web.MockRequest(web.MockRequestParameters{UserData: &session, Request: &http.Request{URL: u}}))
	if werr == nil {
		t.Fatalf("No error seen when one expected")
	}
	if data != nil {
		t.Fatalf("Data returned when none expected")
	}

	_, err = UserStore.EditUser(user, editUserParameters{
		Permissions: UserPermissions{
			CanAccessAuditLog: true,
		},
	})
	if err != nil {
		panic(err)
	}

	// Ensure the user cannot access the audit log
	data, _, werr = h.EventsGet(web.MockRequest(web.MockRequestParameters{UserData: &session, Request: &http.Request{URL: u}}))
	if werr != nil {
		t.Fatalf("Unexpected error: %s", werr.Message)
	}
	if data == nil {
		t.Fatalf("No data returned when some expected")
	}
}

func TestPermissionsCanModifyUsers(t *testing.T) {
	// Make other user with full permissions to avoid lockout errors
	if _, err := UserStore.NewUser(newUserParameters{
		Username:    randomString(3),
		Password:    randomString(6),
		Permissions: UserPermissionsMax(),
	}); err != nil {
		panic(err)
	}

	user, err := UserStore.NewUser(newUserParameters{
		Username:    randomString(3),
		Password:    randomString(6),
		Permissions: UserPermissionsMin(),
	})
	if err != nil {
		panic(err)
	}

	session := SessionStore.NewSessionForUser(user)
	h := handle{}

	// Ensure user cannot create new user
	data, _, werr := h.UserNew(web.MockRequest(web.MockRequestParameters{UserData: &session}))
	if werr == nil {
		t.Fatalf("No error seen when one expected")
	}
	if data != nil {
		t.Fatalf("Data returned when none expected")
	}

	// Ensure user can change their own password
	data, _, werr = h.UserEdit(web.MockRequest(web.MockRequestParameters{UserData: &session, Parameters: map[string]string{"username": user.Username}, JSONBody: editUserParameters{Password: randomString(4), Permissions: user.Permissions}}))
	if werr != nil {
		t.Fatalf("Unexpected error: %s", werr.Message)
	}
	if data == nil {
		t.Fatalf("No data returned when some expected")
	}

	// Ensure user cannot change their own permissions
	data, _, werr = h.UserEdit(web.MockRequest(web.MockRequestParameters{UserData: &session, Parameters: map[string]string{"username": user.Username}, JSONBody: editUserParameters{Password: randomString(4), Permissions: UserPermissionsMax()}}))
	if werr == nil {
		t.Fatalf("No error seen when one expected")
	}
	if data != nil {
		t.Fatalf("Data returned when none expected")
	}

	_, err = UserStore.EditUser(user, editUserParameters{
		Permissions: UserPermissions{
			CanModifyUsers: true,
		},
	})
	if err != nil {
		panic(err)
	}

	data, _, werr = h.UserEdit(web.MockRequest(web.MockRequestParameters{UserData: &session, Parameters: map[string]string{"username": user.Username}, JSONBody: editUserParameters{Password: randomString(4), Permissions: UserPermissionsMax()}}))
	if werr != nil {
		t.Fatalf("Unexpected error: %s", werr.Message)
	}
	if data == nil {
		t.Fatalf("No data returned when some expected")
	}
}

func TestPermissionsCanModifyAutoregister(t *testing.T) {
	group, err := GroupStore.NewGroup(newGroupParameters{
		Name: randomString(3),
		Environment: []environ.Variable{
			{
				Key:    "key",
				Value:  "value",
				Secret: true,
			},
		},
	})
	if err != nil {
		panic(err)
	}

	user, err := UserStore.NewUser(newUserParameters{
		Username:    randomString(3),
		Password:    randomString(6),
		Permissions: UserPermissionsMin(),
	})
	if err != nil {
		panic(err)
	}

	rule, err := RegisterRuleStore.NewRule(newRegisterRuleParams{
		Name: randomString(3),
		Clauses: []RegisterRuleClause{
			{
				Property: RegisterRulePropertyHostname,
				Pattern:  ".*",
			},
		},
		GroupID: group.ID,
	})
	if err != nil {
		panic(err)
	}

	session := SessionStore.NewSessionForUser(user)
	h := handle{}

	// Ensure user cannot create new rule
	data, _, werr := h.RegisterRuleNew(web.MockRequest(web.MockRequestParameters{UserData: &session}))
	if werr == nil {
		t.Fatalf("No error seen when one expected")
	}
	if data != nil {
		t.Fatalf("Data returned when none expected")
	}

	// Ensure user can read existing rule
	data, _, werr = h.RegisterRuleGet(web.MockRequest(web.MockRequestParameters{UserData: &session, Parameters: map[string]string{"id": rule.ID}}))
	if werr != nil {
		t.Fatalf("Unexpected error: %s", werr.Message)
	}
	if data == nil {
		t.Fatalf("No data returned when some expected")
	}

	_, err = UserStore.EditUser(user, editUserParameters{
		Permissions: UserPermissions{
			CanModifyAutoregister: true,
		},
	})
	if err != nil {
		panic(err)
	}

	data, _, werr = h.RegisterRuleEdit(web.MockRequest(web.MockRequestParameters{UserData: &session, Parameters: map[string]string{"id": rule.ID}, JSONBody: *rule}))
	if werr != nil {
		t.Fatalf("Unexpected error: %s", werr.Message)
	}
	if data == nil {
		t.Fatalf("No data returned when some expected")
	}
}

func TestPermissionsCanModifySystem(t *testing.T) {
	user, err := UserStore.NewUser(newUserParameters{
		Username:    randomString(3),
		Password:    randomString(6),
		Permissions: UserPermissionsMin(),
	})
	if err != nil {
		panic(err)
	}

	session := SessionStore.NewSessionForUser(user)
	h := handle{}

	// Ensure user cannot access system settings
	data, _, werr := h.OptionsGet(web.MockRequest(web.MockRequestParameters{UserData: &session}))
	if werr == nil {
		t.Fatalf("No error seen when one expected")
	}
	if data != nil {
		t.Fatalf("Data returned when none expected")
	}
	data, _, werr = h.State(web.MockRequest(web.MockRequestParameters{UserData: &session}))
	if werr != nil {
		t.Fatalf("Unexpected error: %s", werr.Message)
	}
	if data == nil {
		t.Fatalf("No data returned when some expected")
	}
	d, _ := json.Marshal(data)
	if !bytes.ContainsAny(d, "\"Options\":null") {
		t.Fatalf("Options data returned when it shouldnt")
	}

	// Ensure user cannot modify system settings
	data, _, werr = h.OptionsSet(web.MockRequest(web.MockRequestParameters{UserData: &session, JSONBody: *Options}))
	if werr == nil {
		t.Fatalf("No error seen when one expected")
	}
	if data != nil {
		t.Fatalf("Data returned when none expected")
	}

	_, err = UserStore.EditUser(user, editUserParameters{
		Permissions: UserPermissions{
			CanModifySystem: true,
		},
	})
	if err != nil {
		panic(err)
	}

	data, _, werr = h.OptionsSet(web.MockRequest(web.MockRequestParameters{UserData: &session, JSONBody: *Options}))
	if werr != nil {
		t.Fatalf("Unexpected error: %s", werr.Message)
	}
	if data == nil {
		t.Fatalf("No data returned when some expected")
	}
}

func TestPermissionsCantRemoveUserPermissions(t *testing.T) {
	for _, user := range UserCache.All() {
		user.Permissions.CanModifyUsers = false
		UserStore.Table.StartWrite(func(tx ds.IReadWriteTransaction) error {
			return tx.Update(user)
		})
	}

	user, err := UserStore.NewUser(newUserParameters{
		Username:    randomString(3),
		Password:    randomString(6),
		Permissions: UserPermissionsMax(),
	})
	if err != nil {
		panic(err)
	}
	user.Permissions.CanModifyUsers = false
	params := editUserParameters{
		CanLogIn:           true,
		MustChangePassword: false,
		Permissions:        user.Permissions,
	}

	session := SessionStore.NewSessionForUser(user)
	h := handle{}

	data, _, werr := h.UserEdit(web.MockRequest(web.MockRequestParameters{UserData: &session, JSONBody: params}))
	if werr == nil {
		t.Fatalf("No error seen when one expected")
	}
	if data != nil {
		t.Fatalf("Data returned when none expected")
	}
}
