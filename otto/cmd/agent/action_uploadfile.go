package main

import (
	"fmt"
	"io"
	"os"
	"path"
	"strconv"

	"github.com/ecnepsnai/otto/shared/otto"
	"github.com/ecnepsnai/secutil"
)

func handleTriggerActionUploadFile(conn *otto.Connection, message otto.MessageTriggerActionUploadFile) string {
	err := uploadFile(message.FileInfo, func(f io.Writer) error {
		log.Debug("Telling server we're ready for script data")
		if err := conn.WriteMessage(otto.MessageTypeReadyForData, nil); err != nil {
			log.PError("Error replying to server", map[string]interface{}{
				"error": err.Error(),
			})
			return err
		}

		totalCopied := uint64(0)
		var fileBuffer = make([]byte, 1024)
		for totalCopied < message.Length {
			read, err := conn.ReadData(fileBuffer)
			if err != nil && err != io.EOF {
				return err
			}
			f.Write(fileBuffer[0:read])
			totalCopied += uint64(read)
		}
		return nil
	})
	if err != nil {
		return err.Error()
	}
	return ""
}

func createDirectoryForOttoFile(fileInfo otto.FileInfo) error {
	dirName := path.Dir(fileInfo.Path)
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

	if !fileInfo.Owner.Inherit && !compareUIDandGID(fileInfo.Owner.UID, fileInfo.Owner.GID) {
		if err := os.Chown(dirName, int(fileInfo.Owner.UID), int(fileInfo.Owner.GID)); err != nil {
			log.Error("Error chowning directory: directory='%s' error='%s'", dirName, err.Error())
			return err
		}
	}

	log.Debug("Created directory: directory='%s'", dirName)
	return nil
}

func uploadFile(fileInfo otto.FileInfo, writeFunc func(f io.Writer) error) error {
	Stats.FilesUploaded++

	if err := createDirectoryForOttoFile(fileInfo); err != nil {
		return err
	}

	if fileInfo.Path == "" {
		return fmt.Errorf("invalid file path")
	}
	if fileInfo.Path[0] != '/' {
		return fmt.Errorf("file path must be absolute")
	}

	mode, err := intToFileMode(fileInfo.Mode)
	if err != nil {
		return fmt.Errorf("invalid file permissions: %s", err.Error())
	}

	atomicFilePath := fileInfo.Path + "_" + secutil.RandomString(3)
	f, err := os.OpenFile(atomicFilePath, os.O_CREATE|os.O_RDWR, mode)
	if err != nil {
		log.PError("Error writing file", map[string]interface{}{
			"path":  atomicFilePath,
			"error": err.Error(),
		})
		return err
	}
	if err := writeFunc(f); err != nil {
		log.PError("Error writing file", map[string]interface{}{
			"path":  atomicFilePath,
			"error": err.Error(),
		})
		return err
	}
	f.Close()

	// Only chown if we need to
	if !fileInfo.Owner.Inherit && !compareUIDandGID(fileInfo.Owner.UID, fileInfo.Owner.GID) {
		if err := os.Chown(atomicFilePath, int(fileInfo.Owner.UID), int(fileInfo.Owner.GID)); err != nil {
			log.PError("Error setting file owner", map[string]interface{}{
				"path":  atomicFilePath,
				"error": err.Error(),
			})
			return err
		}
		log.PDebug("Set file owner", map[string]interface{}{
			"path": atomicFilePath,
			"uid":  fileInfo.Owner.UID,
			"gid":  fileInfo.Owner.GID,
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
	if checksum != fileInfo.Checksum {
		log.PError("File checksum validation failed!", map[string]interface{}{
			"path":              atomicFilePath,
			"expected_checksum": fileInfo.Checksum,
			"actual_checksum":   checksum,
		})
		os.Remove(atomicFilePath)
		return fmt.Errorf("checksum validation failed")
	}

	if err := os.Rename(atomicFilePath, fileInfo.Path); err != nil {
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
