// Package snapshot provides a way to collect information about a running go application
package snapshot

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"runtime/debug"
)

// Snapshot describes a snapshot of a running go program
type Snapshot struct {
	Memory        runtime.MemStats
	GC            debug.GCStats
	Stack         string
	BuildInfo     debug.BuildInfo
	NumGoRoutines int
	Pid           int
	Uid           int
	Gid           int
	Environ       []string
	Wd            string
	Hostname      string
}

// Collect will take a snapshot of useful statistics of your running Go application
func Collect() (s Snapshot) {
	runtime.ReadMemStats(&s.Memory)
	debug.ReadGCStats(&s.GC)
	buildInfo, _ := debug.ReadBuildInfo()
	s.BuildInfo = *buildInfo
	s.Stack = string(debug.Stack())
	s.NumGoRoutines = runtime.NumGoroutine()
	s.Pid = os.Getpid()
	s.Uid = os.Getuid()
	s.Gid = os.Getgid()
	s.Environ = os.Environ()
	wd, _ := os.Getwd()
	s.Wd = wd
	hostname, _ := os.Hostname()
	s.Hostname = hostname

	return
}

// Full will take a full detailed snapshot of your go application, including memory dumps, and save it as a ZIP file at
// the given path. fileName should end with ".zip"
//
// The ZIP file will contain two files: snapshot.json and heap.dump
//
// Warning: the size of the output file will be at most the amount of memory used by the go application.
func Full(fileName string) error {
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return fmt.Errorf("open: %s", err.Error())
	}
	defer f.Close()

	zw := zip.NewWriter(f)
	sn := Collect()

	snapshotFile, err := zw.Create("snapshot.json")
	if err != nil {
		return fmt.Errorf("snapshot: %s", err.Error())
	}

	encoder := json.NewEncoder(snapshotFile)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(sn); err != nil {
		return fmt.Errorf("snapshot: %s", err.Error())
	}

	tmpFile, err := os.CreateTemp("", "dump")
	if err != nil {
		return fmt.Errorf("dump: %s", err.Error())
	}
	debug.WriteHeapDump(tmpFile.Fd())
	tmpFile.Seek(0, 0)

	dumpFile, err := zw.Create("heap.dump")
	if err != nil {
		return fmt.Errorf("dump: %s", err.Error())
	}

	io.Copy(dumpFile, tmpFile)
	tmpFile.Close()
	os.Remove(tmpFile.Name())

	zw.Close()
	return nil
}
