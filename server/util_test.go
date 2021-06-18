package server

import "testing"

func TestStripPortFromRemoteAddr(t *testing.T) {
	check := func(in, expected string) {
		result := stripPortFromRemoteAddr(in)
		if result != expected {
			t.Errorf("Unexpected result. Expected '%s' got '%s'", expected, result)
		}
	}

	check("127.0.0.1:1234", "127.0.0.1")
	check("[fe80::1]:1234", "fe80::1")
}
