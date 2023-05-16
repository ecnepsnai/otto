package server

import "reflect"

// User describes a user object
type User struct {
	Username           string `ds:"primary" max:"32" min:"1"`
	CanLogIn           bool
	MustChangePassword bool
	Permissions        UserPermissions
}

// UserPermissions describes permissions for this user
type UserPermissions struct {
	ScriptRunLevel        int
	CanModifyHosts        bool
	CanModifyGroups       bool
	CanModifyScripts      bool
	CanModifySchedules    bool
	CanAccessAuditLog     bool
	CanModifyUsers        bool
	CanModifyAutoregister bool
	CanModifySystem       bool
}

func (p UserPermissions) EqualTo(o UserPermissions) bool {
	return reflect.DeepEqual(p, o)
}

func UserPermissionsMax() UserPermissions {
	return UserPermissions{
		ScriptRunLevel:        ScriptRunLevelReadWrite,
		CanModifyHosts:        true,
		CanModifyGroups:       true,
		CanModifyScripts:      true,
		CanModifySchedules:    true,
		CanAccessAuditLog:     true,
		CanModifyUsers:        true,
		CanModifyAutoregister: true,
		CanModifySystem:       true,
	}
}

func UserPermissionsMin() UserPermissions {
	return UserPermissions{
		ScriptRunLevel:        ScriptRunLevelNone,
		CanModifyHosts:        false,
		CanModifyGroups:       false,
		CanModifyScripts:      false,
		CanModifySchedules:    false,
		CanAccessAuditLog:     false,
		CanModifyUsers:        false,
		CanModifyAutoregister: false,
		CanModifySystem:       false,
	}
}
