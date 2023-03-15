package main

import (
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/ecnepsnai/otto/shared/otto"
	"github.com/ecnepsnai/secutil"
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
	Stats.FilesUploaded++

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

	atomicFilePath := file.Path + "_" + secutil.RandomString(3)
	if err := os.WriteFile(atomicFilePath, file.Data, mode); err != nil {
		log.PError("Error writing file", map[string]interface{}{
			"path":  atomicFilePath,
			"error": err.Error(),
		})
		return err
	}

	// Only chown if we need to
	if !file.Owner.Inherit && !compareUIDandGID(file.Owner.UID, file.Owner.GID) {
		if err := os.Chown(atomicFilePath, int(file.Owner.UID), int(file.Owner.GID)); err != nil {
			log.PError("Error setting file owner", map[string]interface{}{
				"path":  atomicFilePath,
				"error": err.Error(),
			})
			return err
		}
		log.PDebug("Set file owner", map[string]interface{}{
			"path": atomicFilePath,
			"uid":  file.Owner.UID,
			"gid":  file.Owner.GID,
		})
	}
	log.PDebug("Created file", map[string]interface{}{
		"path": atomicFilePath,
	})

	checksum, err := getFileSHA256Checksum(atomicFilePath)
	if err != nil {
		log.PError("Error calculating file checksum", map[string]interface{}{
			"path":  atomicFilePath,
			"error": err.Error(),
		})
		os.Remove(atomicFilePath)
		return err
	}
	if checksum != file.Checksum {
		log.PError("File checksum validation failed!", map[string]interface{}{
			"path":              atomicFilePath,
			"expected_checksum": file.Checksum,
			"actual_checksum":   checksum,
		})
		os.Remove(atomicFilePath)
		return fmt.Errorf("checksum validation failed")
	}

	if err := os.Rename(atomicFilePath, file.Path); err != nil {
		log.PError("Error renaming atomic file", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}

	return nil
}

func intToFileMode(m uint32) (os.FileMode, error) {
	n, err := strconv.ParseUint(fmt.Sprintf("0%d", m), 8, 32)
	if err != nil {
		return os.ModePerm, err
	}
	return os.FileMode(n), nil
}
