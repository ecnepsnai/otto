package server

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/ecnepsnai/otto"
)

// File describes a file for a script
type File struct {
	ID   string `ds:"primary"`
	Path string
	UID  int
	GID  int
	Mode uint32
}

// OttoFile return an otto common file
func (file File) OttoFile() (*otto.File, error) {
	f, err := os.OpenFile(file.FilePath(), os.O_RDONLY, 0644)
	if err != nil {
		log.Error("Error opening script file '%s': %s", file.ID, err.Error())
		return nil, err
	}
	defer f.Close()
	fileData, err := ioutil.ReadAll(f)
	if err != nil {
		log.Error("Error reading file '%s': %s", file.ID, err.Error())
		return nil, err
	}

	ottoFile := otto.File{
		Path: file.Path,
		UID:  file.UID,
		GID:  file.GID,
		Mode: file.Mode,
		Data: fileData,
	}
	return &ottoFile, nil
}

// FilePath returns the absolute path for this file
func (file File) FilePath() string {
	return path.Join(Directories.Files, file.ID)
}
