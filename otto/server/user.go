package server

// User describes a user object
type User struct {
	Username           string `ds:"primary" max:"32" min:"1"`
	CanLogIn           bool
	MustChangePassword bool
}
