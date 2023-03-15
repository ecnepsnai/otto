package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

func addressFromSocketString(s string) string {
	// Remove the port first
	portIdx := -1
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == ':' {
			portIdx = i
			break
		}
	}

	s = s[:portIdx]

	if s[0] == '[' && s[len(s)-1] == ']' {
		s = s[1 : len(s)-1]
	}

	return s
}

func fileExists(pathname string) bool {
	_, err := os.Stat(pathname)
	if err != nil {
		return !os.IsNotExist(err)
	}
	return true
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
