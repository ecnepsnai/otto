package main

import "testing"

func TestAddressFromSocketString(t *testing.T) {
	test := func(in, expected string) {
		actual := addressFromSocketString(in)
		if actual != expected {
			t.Errorf("Incorrect value for addressFromSocketString: given '%s' expected '%s' got '%s'", in, expected, actual)
		}
	}

	test("127.0.0.1:8080", "127.0.0.1")
	test("[::1]:8080", "::1")
}
