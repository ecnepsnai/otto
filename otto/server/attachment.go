package server

import (
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
}

// OttoFile return an otto common file
func (attachment Attachment) OttoFile() (*otto.File, error) {
	f, err := os.OpenFile(attachment.FilePath(), os.O_RDONLY, 0644)
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

	ottoFile := otto.File{
		Path: attachment.Path,
		Mode: attachment.Mode,
		Owner: otto.RunAs{
			UID:     attachment.Owner.UID,
			GID:     attachment.Owner.GID,
			Inherit: attachment.Owner.Inherit,
		},
		Data:        fileData,
		AfterScript: attachment.AfterScript,
	}
	return &ottoFile, nil
}

// FilePath returns the absolute path for this attachment
func (attachment Attachment) FilePath() string {
	return path.Join(Directories.Attachments, attachment.ID)
}
