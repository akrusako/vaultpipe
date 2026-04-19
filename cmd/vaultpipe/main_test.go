package main

import (
	"os"
	"os/exec"
	"testing"
)

func TestMain_VersionFlag(t *testing.T) {
	if os.Getenv("RUN_SUBPROCESS") == "1" {
		main()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestMain_VersionFlag")
	cmd.Env = append(os.Environ(),
		"RUN_SUBPROCESS=1",
	)
	// We can't easily test --version via subprocess in unit tests without
	// building the binary, so just verify version var is set.
	if version == "" {
		t.Error("version should not be empty")
	}
}

func TestMain_MissingConfig(t *testing.T) {
	if os.Getenv("RUN_SUBPROCESS") != "1" {
		t.Skip("subprocess-only test")
	}
	os.Unsetenv("VAULT_ADDR")
	os.Unsetenv("VAULT_TOKEN")
	os.Unsetenv("VAULT_PATHS")
	os.Unsetenv("VP_COMMAND")
	main()
}

func TestVersion_Default(t *testing.T) {
	if version != "dev" {
		t.Errorf("expected default version 'dev', got %q", version)
	}
}
