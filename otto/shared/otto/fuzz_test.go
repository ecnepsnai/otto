package otto

import (
	"bytes"
	"testing"
	"time"
)

func FuzzReadFrame(f *testing.F) {
	f.Add([]byte(time.Now().String()))
	f.Fuzz(func(t *testing.T, a []byte) {
		buf := bytes.NewBuffer(a)
		frame, err := readFrame(buf)

		if err == nil && len(a) > 0 {
			t.Fatalf("No error seen with fuzzed frame.")
		}
		if len(frame) > 0 {
			t.Fatalf("Frame data returned with fuzzed data.")
		}
	})
}
