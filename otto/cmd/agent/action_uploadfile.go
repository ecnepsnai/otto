package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"

	"github.com/ecnepsnai/otto/shared/otto"
)

func createDirectoryForOttoFile(file otto.File) error {
	dirName := path.Dir(file.Path)
	info, err := os.Stat(dirName)
	if err == nil && info.IsDir() {
		return nil
	}

	if !os.IsNotExist(err) {
		log.Error("Error performing stat on directory: directory='%s' error='%s'", dirName, err.Error())
		return err
	}

	if err := os.MkdirAll(dirName, os.ModePerm); err != nil {
		log.Error("Error creating directory: directory='%s' error='%s'", dirName, err.Error())
		return err
	}

	if !file.Owner.Inherit && !compareUIDandGID(file.Owner.UID, file.Owner.GID) {
		if err := os.Chown(dirName, int(file.Owner.UID), int(file.Owner.GID)); err != nil {
			log.Error("Error chowning directory: directory='%s' error='%s'", dirName, err.Error())
			return err
		}
	}

	log.Debug("Created directory: directory='%s'", dirName)
	return nil
}

func uploadFile(file otto.File) error {
	if err := createDirectoryForOttoFile(file); err != nil {
		return err
	}

	if file.Path == "" {
		return fmt.Errorf("invalid file path")
	}
	if file.Path[0] != '/' {
		return fmt.Errorf("file path must be absolute")
	}

	mode, err := intToFileMode(file.Mode)
	if err != nil {
		return fmt.Errorf("invalid file permissions: %s", err.Error())
	}

	f, err := os.OpenFile(file.Path, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, mode)
	if err != nil {
		log.Error("Error opening file: path='%s' error='%s'", file.Path, err.Error())
		return err
	}
	defer f.Close()

	n, err := io.Copy(f, bytes.NewReader(file.Data))
	if err != nil {
		log.Error("Error writing to file: path='%s' error='%s'", file.Path, err.Error())
		return err
	}
	// Only chown if we need to
	if !file.Owner.Inherit && !compareUIDandGID(file.Owner.UID, file.Owner.GID) {
		if err := f.Chown(int(file.Owner.UID), int(file.Owner.GID)); err != nil {
			log.Error("Error chowning file: path='%s' error='%s'", file.Path, err.Error())
			return err
		}
		log.Debug("Chowned file: path='%s' uid=%d gid=%d", file.Path, file.Owner.UID, file.Owner.GID)
	}

	log.Debug("Wrote %d bytes to '%s'", n, file.Path)

	return nil
}

func intToFileMode(m uint32) (os.FileMode, error) {
	n, err := strconv.ParseUint(fmt.Sprintf("0%d", m), 8, 32)
	if err != nil {
		return os.ModePerm, err
	}
	return os.FileMode(n), nil
}
