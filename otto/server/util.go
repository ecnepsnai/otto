package server

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ecnepsnai/secutil"
	nanoid "github.com/matoous/go-nanoid"
)

func generateSessionSecret() string {
	return secutil.RandomString(64)
}

// sliceContains does slice h contain n?
func sliceContains[T string](n T, h []T) bool {
	for _, s := range h {
		if s == n {
			return true
		}
	}
	return false
}

// sliceContainsFold does this slice of strings contain n? (cast insensitive)
func sliceContainsFold(n string, h []string) bool {
	for _, s := range h {
		if strings.EqualFold(s, n) {
			return true
		}
	}
	return false
}

// filterSlice remove any occurrence of `r` from `s`, returning a new slice
func filterSlice[T string](r T, s []T) []T {
	sl := []T{}
	for _, i := range s {
		if i == r {
			continue
		}
		sl = append(sl, i)
	}
	return sl
}

func getFileSHA256Checksum(filePath string) (string, error) {
	h := sha256.New()

	f, err := os.OpenFile(filePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

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

// newAPIKey returns a 64-character string suitable for an otto API key. Keys are always prefixed with `otto_`.
func newAPIKey() string {
	id, err := nanoid.Generate("BCDFGHJKLMNPQRSTVWXYZbcdfghjklmnpqrstvwxyz1234567890", 59)
	if err != nil {
		panic(err)
	}
	return "otto_" + id
}
