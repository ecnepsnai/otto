package server

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/ecnepsnai/logtic"
	"github.com/ecnepsnai/security"
)

var tmpDir *string
var verbose bool
var offline bool

// Perform all startup actions that are typically done during Start(), except don't start
// the http router.
func testSetup() {
	tmp, err := ioutil.TempDir("", "otto")
	if err != nil {
		panic(err)
	}
	tmpDir = &tmp

	// Overwrite the operating directory with the temporary directory
	operatingDirectory = *tmpDir
	Directories = apiDirectories{
		Base:   operatingDirectory,
		Logs:   path.Join(operatingDirectory, "logs"),
		Data:   path.Join(operatingDirectory, "data"),
		Static: path.Join(operatingDirectory, "static"),
	}

	os.Mkdir(Directories.Logs, os.ModePerm)
	os.Mkdir(Directories.Data, os.ModePerm)
	os.Mkdir(Directories.Static, os.ModePerm)

	if verbose {
		initLogtic(true)
	}

	GobSetup()
	StateSetup()
	DataStoreSetup()
	LoadOptions()
}

// Close everything and delete the operating directory
func testTeardown() {
	State.Close()
	DataStoreTeardown()
	logtic.Close()
	if tmpDir != nil {
		os.RemoveAll(*tmpDir)
	}
}

func TestMain(m *testing.M) {

	for _, arg := range os.Args {
		if arg == "-test.v=true" {
			verbose = true
		} else if arg == "-test.short=true" {
			offline = true
		}
	}

	testSetup()
	retCode := m.Run()
	testTeardown()
	os.Exit(retCode)
}

func randomString(length uint16) string {
	return security.RandomString(length)
}
