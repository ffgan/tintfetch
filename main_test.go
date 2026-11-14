package main

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ffgan/tintfetch/cmd"
	"github.com/ffgan/tintfetch/configs"
)

func invoke(args ...string) (int, string, string) {
	var out, diagnostics bytes.Buffer
	code := cmd.Execute(context.Background(), args, strings.NewReader(""), &out, &diagnostics, assets)
	return code, out.String(), diagnostics.String()
}
func TestMetadataIgnoresBrokenConfig(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bad.conf")
	if err := os.WriteFile(path, []byte("not a config"), 0600); err != nil {
		t.Fatal(err)
	}
	for _, flag := range []string{"--version", "--help", "--print-config", "--list-presets"} {
		code, _, err := invoke("--config", path, flag)
		if code != 0 {
			t.Fatalf("%s: %d %s", flag, code, err)
		}
	}
}
func TestCLIConfigPrecedenceAndJSON(t *testing.T) {
	c := configs.DefaultConfig()
	c.InfoList = [][]string{{"Operating system", "distro"}}
	c.OSArch = "off"
	c.DistroShorthand = "tiny"
	path := filepath.Join(t.TempDir(), "config.conf")
	if err := os.WriteFile(path, []byte(configs.Format(&c)), 0600); err != nil {
		t.Fatal(err)
	}
	code, out, err := invoke("--config", path, "--os_arch=on", "--json")
	if code != 0 {
		t.Fatalf("exit %d: %s", code, err)
	}
	var data struct {
		Entries []struct{ Label, Module, Value string }
	}
	if err := json.Unmarshal([]byte(out), &data); err != nil {
		t.Fatal(err)
	}
	if len(data.Entries) != 1 || data.Entries[0].Label != "Operating system" || strings.Contains(out, "\x1b") {
		t.Fatalf("unexpected output %s", out)
	}
	_, without, _ := invoke("--config", path, "--json")
	if out == without {
		t.Fatal("explicit os_arch flag did not override configuration")
	}
}
func TestInvalidCLI(t *testing.T) {
	for _, args := range [][]string{{"--no-config", "--bar_length=-1"}, {"--no-config", "--gpu_type=integrated"}, {"--no-config", "--backend=kitty"}, {"--no-config", "--config=anything"}, {"--no-config", "--shell-config"}, {"--preset=unknown"}, {"--engine=unknown"}, {"--no-config", "--cpu_cores=invalid"}} {
		code, out, err := invoke(args...)
		if code != 2 || out != "" || err == "" {
			t.Fatalf("%v: code=%d stdout=%q stderr=%q", args, code, out, err)
		}
	}
}
func TestNoImplicitShell(t *testing.T) {
	dir := t.TempDir()
	marker := filepath.Join(dir, "executed")
	path := filepath.Join(dir, "config.conf")
	if err := os.WriteFile(path, []byte("separator=$(touch '"+marker+"')\n"), 0600); err != nil {
		t.Fatal(err)
	}
	code, _, _ := invoke("--config", path)
	if code != 2 {
		t.Fatalf("expected invalid config, got %d", code)
	}
	if _, err := os.Stat(marker); !os.IsNotExist(err) {
		t.Fatal("shell code was executed")
	}
}
func TestArtworkFallback(t *testing.T) {
	c := configs.DefaultConfig()
	c.InfoList = [][]string{{"OS", "distro"}}
	c.ASCIIDistro = "unknown distro"
	c.Stdout = "off"
	path := filepath.Join(t.TempDir(), "config.conf")
	if err := os.WriteFile(path, []byte(configs.Format(&c)), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("NO_COLOR", "")
	code, out, err := invoke("--config", path, "--preset=trans")
	if code != 0 || !strings.Contains(out, "\x1b[38;2;") {
		t.Fatalf("fallback: %d %q %s", code, out, err)
	}
	t.Setenv("NO_COLOR", "1")
	code, out, err = invoke("--config", path, "--preset=trans")
	if code != 0 || strings.Contains(out, "\x1b") {
		t.Fatalf("NO_COLOR: %d %q %s", code, out, err)
	}
}
