package server

import (
	"github.com/ecnepsnai/nanoid"
)

// NewID generate a new
func NewID() string {
	id, err := nanoid.Generate("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890", 11)
	if err != nil {
		panic(err)
	}
	return id
}
