package server

import (
	"os"
	"path"
	"testing"
	"time"
)

// Test the byte formatter
func TestUtilFormatByte(t *testing.T) {
	test := func(in uint64, expected string) {
		result := FormatByte(in)
		if result != expected {
			t.Errorf("Unexpected result from FormatByte. Expected '%s' got '%s'", expected, result)
		}
	}

	test(1024, "1.0 KiB")
	test(10240, "10.0 KiB")
	test(102400, "100.0 KiB")
	test(1024000, "1000.0 KiB")
	test(438143210, "417.8 MiB")
	test(57435943275, "53.5 GiB")
	test(2482587438925, "2.3 TiB")
	test(957183938585752, "870.6 TiB")
	test(65718393858575225, "58.4 PiB")
}

// Test the path sanitizer
func TestUtilSanitizePath(t *testing.T) {
	test := func(in, expected string) {
		result := SanitizePath(in)
		if result != expected {
			t.Errorf("Unexpected result from Sanitized Path. Expected '%s' got '%s'", expected, result)
		}
	}

	test("../../../../../../etc/passwd", "_............etcpasswd")
	test("\\\\127.0.0.1\\foo", "127.0.0.1foo")
	test("blah\000", "blah")
}

// Test that the touch directory method updates the MTime for directories
func TestUtilTouchDirectory(t *testing.T) {
	dir := path.Join(*tmpDir, randomString(12))
	subdir := path.Join(dir, randomString(12))
	if err := os.MkdirAll(subdir, os.ModePerm); err != nil {
		t.Fatalf("Error making directory: %s", err.Error())
	}

	getMTime := func(dir string) time.Time {
		dirInfo, err := os.Stat(dir)
		if err != nil {
			t.Fatalf("Error statting directory: %s", err.Error())
		}
		return dirInfo.ModTime()
	}

	dirMtime := getMTime(dir)
	subdirMtime := getMTime(subdir)

	TouchDirectory(subdir)
	if TimeEquals(subdirMtime, getMTime(subdir)) {
		t.Errorf("Mtime did not change for sub directory when expected")
	}

	if !TimeEquals(dirMtime, getMTime(dir)) {
		t.Errorf("Unexpected mtime change for parent directory")
	}

	TouchDirectory(dir)
	if TimeEquals(dirMtime, getMTime(dir)) {
		t.Errorf("Mtime did not change for parent directory when expected")
	}
}
