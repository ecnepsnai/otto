package server

import (
	"github.com/ecnepsnai/nanoid"
)

func newID() string {
	id, err := nanoid.Generate("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890._-", 12)
	if err != nil {
		panic(err)
	}
	return id
}

func newPlainID() string {
	id, err := nanoid.Generate("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890", 12)
	if err != nil {
		panic(err)
	}
	return id
}
