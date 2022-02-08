package server

import (
	"fmt"
)

// Error describes an error object
type Error struct {
	Server  bool
	Message string
	Error   error
}

// ErrorUser create a new user-facing error
func ErrorUser(format string, a ...interface{}) *Error {
	message := fmt.Sprintf(format, a...)
	return &Error{
		Server:  false,
		Message: message,
		Error:   fmt.Errorf(message),
	}
}

// ErrorServer create a new server-side error
func ErrorServer(format string, a ...interface{}) *Error {
	message := fmt.Sprintf(format, a...)
	return &Error{
		Server:  true,
		Message: message,
		Error:   fmt.Errorf(message),
	}
}

// ErrorFrom create a new server-side error from the given error
func ErrorFrom(err error) *Error {
	return &Error{
		Server:  true,
		Message: err.Error(),
		Error:   err,
	}
}
