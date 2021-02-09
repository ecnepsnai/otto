package server

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"

	"github.com/ecnepsnai/secutil"
	nanoid "github.com/matoous/go-nanoid"
)

func generateSessionSecret() string {
	return secutil.RandomString(64)
}

// stringSliceContains does this slice of strings contain n?
func stringSliceContains(n string, h []string) bool {
	for _, s := range h {
		if s == n {
			return true
		}
	}
	return false
}

// filterStringSlice remove any occurrence of `r` from `s`, returning a new slice
func filterStringSlice(r string, s []string) []string {
	sl := []string{}
	for _, i := range s {
		if i == r {
			continue
		}
		sl = append(sl, i)
	}
	return sl
}

func sliceFirst(s []string) string {
	if len(s) == 0 {
		return ""
	}
	return s[0]
}

func hashFile(filePath string) (string, error) {
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
