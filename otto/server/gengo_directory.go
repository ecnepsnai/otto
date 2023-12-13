package server

// This file is was generated automatically by GenGo v1.13.0
// Do not make changes to this file as they will be lost

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
)

func getAPIOperatingDir() string {
	ex, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to determine working directory: %s\n"+err.Error())
		os.Exit(1)
	}
	return filepath.Dir(ex)
}

var operatingDirectory = getAPIOperatingDir()
var dataDirectory = getAPIOperatingDir()

type apiDirectories struct {
	Base        string
	Agents      string
	Data        string
	Attachments string
	Logs        string
	Static      string
}

// Directories absolute paths of API related directires.
var Directories = apiDirectories{}

func fsSetup() {
	Directories = apiDirectories{
		Base: operatingDirectory,

		Agents: path.Join(operatingDirectory, "agents"),

		Data: path.Join(dataDirectory, "data"),

		Attachments: path.Join(dataDirectory, "data", "attachments"),

		Logs: path.Join(dataDirectory, "logs"),

		Static: path.Join(operatingDirectory, "static"),
	}

	MakeDirectoryIfNotExist(Directories.Agents)

	MakeDirectoryIfNotExist(Directories.Data)

	MakeDirectoryIfNotExist(Directories.Attachments)

	MakeDirectoryIfNotExist(Directories.Logs)

	MakeDirectoryIfNotExist(Directories.Static)

}

// DirectoryExists does the given directory exist (and is it a directory)
func DirectoryExists(directoryPath string) bool {
	stat, err := os.Stat(directoryPath)
	return err == nil && stat.IsDir()
}

// MakeDirectoryIfNotExist make the given directory if it does not exist
func MakeDirectoryIfNotExist(directoryPath string) error {
	if !DirectoryExists(directoryPath) {
		return os.MkdirAll(directoryPath, 0755)
	}
	return nil
}

// FileExists does the given file exist
func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	fmt.Fprintf(os.Stderr, "Error stat-ing file '%s': %s", filePath, err.Error())
	return false
}
