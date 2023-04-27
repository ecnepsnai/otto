package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

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
