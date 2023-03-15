package main

import (
	"crypto/sha256"
	"fmt"
	"path"
	"testing"

	"github.com/ecnepsnai/otto/shared/otto"
	"github.com/ecnepsnai/secutil"
)

func TestUploadFile(t *testing.T) {
	fileDir := t.TempDir()
	filePath := path.Join(fileDir, "test", "example.test")
	data := secutil.RandomBytes(32)
	checksum := fmt.Sprintf("%x", sha256.Sum256(data))

	file := otto.File{
		Path: filePath,
		Owner: otto.RunAs{
			Inherit: true,
		},
		Mode:     644,
		Data:     data,
		Checksum: checksum,
	}

	if err := uploadFile(file); err != nil {
		t.Fatalf("Error uploading otto file: %s", err.Error())
	}

	compareChecksum, err := getFileSHA256Checksum(filePath)
	if err != nil {
		t.Fatalf("Error getting file checksum: %s", err.Error())
	}
	if compareChecksum != checksum {
		t.Fatalf("Unexpected file checksum: %s != %s", compareChecksum, checksum)
	}

	// Overwrite the file
	filePath = path.Join(fileDir, "test", "example.test")
	data = secutil.RandomBytes(32)
	checksum = fmt.Sprintf("%x", sha256.Sum256(data))

	file = otto.File{
		Path: filePath,
		Owner: otto.RunAs{
			Inherit: true,
		},
		Mode:     644,
		Data:     data,
		Checksum: checksum,
	}

	if err := uploadFile(file); err != nil {
		t.Fatalf("Error uploading otto file: %s", err.Error())
	}

	compareChecksum, err = getFileSHA256Checksum(filePath)
	if err != nil {
		t.Fatalf("Error getting file checksum: %s", err.Error())
	}
	if compareChecksum != checksum {
		t.Fatalf("Unexpected file checksum: %s != %s", compareChecksum, checksum)
	}
}

func TestUploadFileBadChecksum(t *testing.T) {
	fileDir := t.TempDir()
	filePath := path.Join(fileDir, "test", "example.test")
	data := secutil.RandomBytes(32)
	checksum := secutil.RandomString(32)

	file := otto.File{
		Path: filePath,
		Owner: otto.RunAs{
			Inherit: true,
		},
		Mode:     644,
		Data:     data,
		Checksum: checksum,
	}

	if err := uploadFile(file); err == nil {
		t.Fatalf("No error seen when one expected for uploading file with bad checksum")
	}

	file = otto.File{
		Path: filePath,
		Owner: otto.RunAs{
			Inherit: true,
		},
		Mode:     644,
		Data:     data,
		Checksum: "",
	}

	if err := uploadFile(file); err == nil {
		t.Fatalf("No error seen when one expected for uploading file with no checksum")
	}
}
