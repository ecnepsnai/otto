package server

import (
	"github.com/ecnepsnai/security"
)

// User describes a user object
type User struct {
	Username           string                  `ds:"primary" max:"32" min:"1"`
	Email              string                  `ds:"unique" max:"128" min:"1"`
	PasswordHash       security.HashedPassword `json:"-"`
	CanLogIn           bool
	MustChangePassword bool
}
