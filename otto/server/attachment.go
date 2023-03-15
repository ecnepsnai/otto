package server

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"github.com/ecnepsnai/otto/shared/otto"
)

// Attachment describes a file for a script
type Attachment struct {
	ID          string `ds:"primary"`
	Path        string
	Name        string
	MimeType    string
	Owner       RunAs
	Created     time.Time
	Modified    time.Time
	Mode        uint32
	Size        uint64
	AfterScript bool
	Checksum    string
}

// GetChecksum get the real checksum for the file path
func (attachment Attachment) GetChecksum() (string, error) {
	checksum, err := getFileSHA256Checksum(attachment.FilePath())
	if err != nil {
		log.PError("Error calculating attachment checksum", map[string]interface{}{
			"attachment": attachment.ID,
			"error":      err.Error(),
		})
		return "", err
	}
	return checksum, nil
}

// OttoFile return an otto common file
func (attachment Attachment) OttoFile() (*otto.File, error) {
	f, err := os.Open(attachment.FilePath())
	if err != nil {
		log.Error("Error opening script file '%s': %s", attachment.ID, err.Error())
		return nil, err
	}
	defer f.Close()
	fileData, err := io.ReadAll(f)
	if err != nil {
		log.Error("Error reading file '%s': %s", attachment.ID, err.Error())
		return nil, err
	}

	checksum := fmt.Sprintf("%x", sha256.Sum256(fileData))
	if attachment.Checksum != checksum {
		log.PError("Attachment file checksum verification failed", map[string]interface{}{
			"attachment":        attachment.ID,
			"expected_checksum": attachment.Checksum,
			"actual_checksum":   checksum,
		})
		return nil, fmt.Errorf("file verification failed: %s", attachment.ID)
	}

	ottoFile := otto.File{
		Path: attachment.Path,
		Mode: attachment.Mode,
		Owner: otto.RunAs{
			UID:     attachment.Owner.UID,
			GID:     attachment.Owner.GID,
			Inherit: attachment.Owner.Inherit,
		},
		Data:        fileData,
		Checksum:    checksum,
		AfterScript: attachment.AfterScript,
	}
	return &ottoFile, nil
}

// FilePath returns the absolute path for this attachment
func (attachment Attachment) FilePath() string {
	return path.Join(Directories.Attachments, attachment.ID)
}

// AtomicFilePath returns the absolute path for this attachment
func (attachment Attachment) AtomicFilePath() string {
	return path.Join(Directories.Attachments, ".atomic_"+attachment.ID)
}
