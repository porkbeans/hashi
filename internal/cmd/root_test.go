package cmd

import "testing"

func TestExecute(t *testing.T) {
	if Execute(nil) != 0 {
		t.Errorf("exit code must be 0")
	}
}

func TestExecuteUnknownSubcommand(t *testing.T) {
	if Execute([]string{"unknown"}) != 1 {
		t.Errorf("exit code must be 1")
	}
}
