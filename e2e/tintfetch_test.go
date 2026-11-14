package e2e_test

import (
	"strings"
	"testing"
)

func TestTintfetchAliases(t *testing.T) {
	binary := buildTintfetch(t)
	text, _ := runTintfetch(t, binary, "--no-config", "--plain")
	if strings.Contains(text, "\x1b") {
		t.Fatal("plain alias emitted ANSI")
	}
	presets, _ := runTintfetch(t, binary, "--list-presets")
	if !strings.Contains(presets, "rainbow") {
		t.Fatal("preset catalog was not listed")
	}
	config, _ := runTintfetch(t, binary, "--print-config")
	if !strings.Contains(config, "print_info()") {
		t.Fatal("print-config alias did not print configuration")
	}
}
