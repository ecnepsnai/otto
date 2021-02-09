package server

import "github.com/ecnepsnai/secutil"

// User describes a user object
type User struct {
	Username           string                 `ds:"primary" max:"32" min:"1"`
	Email              string                 `ds:"unique" max:"128" min:"1"`
	PasswordHash       secutil.HashedPassword `json:"-"`
	CanLogIn           bool
	MustChangePassword bool
}
