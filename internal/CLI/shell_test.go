package cli

import "testing"

func TestShellDisplayOptions(t *testing.T) {
	t.Setenv("SHELL", "/bin/bash")
	t.Setenv("BASH_VERSION", "5.2.15(1)-release")
	if got := GetShell("on", "on"); got != "/bin/bash 5.2.15" {
		t.Fatalf("path/version lost: %q", got)
	}
	if got := GetShell("off", "off"); got != "bash" {
		t.Fatalf("version should be hidden: %q", got)
	}
}
