package server

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"time"

	"github.com/ecnepsnai/security"
)

// GenerateSessionSecret generate a sutable secret for a user session
func GenerateSessionSecret() string {
	return security.RandomString(8)
}

// FormatByte format the given byte number to a string
func FormatByte(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}

// SanitizePath sanitize the path component
func SanitizePath(part string) string {
	p := part

	// Need to remove all NTFS characters,
	// as well as a couple more annoynances
	naughty := map[string]string{
		"<":    "",
		">":    "",
		":":    "",
		"\"":   "",
		"/":    "",
		"\\":   "",
		"|":    "",
		"?":    "",
		"*":    "",
		" ":    "_",
		",":    "",
		"#":    "",
		"\000": "",
	}

	for bad, good := range naughty {
		p = strings.Replace(p, bad, good, -1)
	}

	// Don't allow UNIX "hidden" files
	if p[0] == '.' {
		p = "_" + p
	}

	return p
}

// CurrentYear return the current year as a string
func CurrentYear() string {
	return fmt.Sprintf("%d", time.Now().Year())
}

// TouchDirectory update the Ctime of a directory
func TouchDirectory(directory string) {
	name := path.Join(directory, "."+security.RandomString(12))

	file, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Error("Error touching directory '%s': %s", directory, err.Error())
		return
	}
	file.Close()
	if err := os.Remove(name); err != nil {
		log.Error("Error touching directory '%s': %s", directory, err.Error())
	}
}

// TimeEquals do the given times match
func TimeEquals(a time.Time, b time.Time) bool {
	return a.UnixNano() == b.UnixNano()
}

// Hostname get the system hostname
func Hostname() string {
	h, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	return h
}

// StringSliceContains does this slice of strings contain n?
func StringSliceContains(n string, h []string) bool {
	for _, s := range h {
		if s == n {
			return true
		}
	}
	return false
}

// FilterStringSlice remove any occurrence of `r` from `s`, returning a new slice
func FilterStringSlice(r string, s []string) []string {
	sl := []string{}
	for _, i := range s {
		if i == r {
			continue
		}
		sl = append(sl, i)
	}
	return sl
}

func sliceFirst(s []string) string {
	if len(s) == 0 {
		return ""
	}
	return s[0]
}

func hashFile(filePath string) (string, error) {
	h := sha256.New()

	f, err := os.OpenFile(filePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
