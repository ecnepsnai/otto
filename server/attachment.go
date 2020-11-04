package server

import (
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/ecnepsnai/otto"
)

// Attachment describes a file for a script
type Attachment struct {
	ID       string `ds:"primary"`
	Path     string
	UID      int
	GID      int
	Created  time.Time
	Modified time.Time
	Mode     uint32
}

// OttoFile return an otto common file
func (attachment Attachment) OttoFile() (*otto.File, error) {
	f, err := os.OpenFile(attachment.FilePath(), os.O_RDONLY, 0644)
	if err != nil {
		log.Error("Error opening script file '%s': %s", attachment.ID, err.Error())
		return nil, err
	}
	defer f.Close()
	fileData, err := ioutil.ReadAll(f)
	if err != nil {
		log.Error("Error reading file '%s': %s", attachment.ID, err.Error())
		return nil, err
	}

	ottoFile := otto.File{
		Path: attachment.Path,
		UID:  attachment.UID,
		GID:  attachment.GID,
		Mode: attachment.Mode,
		Data: fileData,
	}
	return &ottoFile, nil
}

// FilePath returns the absolute path for this attachment
func (attachment Attachment) FilePath() string {
	return path.Join(Directories.Attachments, attachment.ID)
}
