package server

import (
	"github.com/ecnepsnai/nanoid"
	uuid "github.com/satori/go.uuid"
)

// NewID generate a new
func NewID() string {
	id, err := nanoid.Generate("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890", 11)
	if err != nil {
		panic(err)
	}
	return id
}

// NewUUID generate a new UUID
func NewUUID() string {
	return uuid.NewV4().String()
}
