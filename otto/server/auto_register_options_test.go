package server

import "testing"

func TestValidateOptionsRegister(t *testing.T) {
	options := &OptionsRegister{
		Enabled: false,
	}

	if options.Validate() != nil {
		t.Errorf("Error seen when one not expected")
	}

	options.Enabled = true

	if options.Validate() == nil {
		t.Errorf("No error seen when one expected")
	}

	options.Key = "invalid\"key#"

	if options.Validate() == nil {
		t.Errorf("No error seen when one expected")
	}

	options.Key = "hunter2"

	if options.Validate() != nil {
		t.Errorf("Error seen when one not expected")
	}
}
