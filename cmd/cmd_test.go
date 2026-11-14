package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func TestEngineArgs(t *testing.T) {
	engine, args, err := engineArgs([]string{"--engine", "hyfetch", "--preset", "trans", "--", "--engine=native"})
	if err != nil || engine != "hyfetch" || !reflect.DeepEqual(args, []string{"--preset", "trans", "--", "--engine=native"}) {
		t.Fatalf("%q %v %v", engine, args, err)
	}
	if _, _, err := engineArgs([]string{"--engine"}); err == nil {
		t.Fatal("missing engine accepted")
	}
}

// The current test executable doubles as a dependency-free fake upstream engine.
func init() {
	if os.Getenv("TINTFETCH_ENGINE_TEST_HELPER") == "1" {
		_ = json.NewEncoder(os.Stdout).Encode(os.Args[1:])
		fmt.Fprintln(os.Stderr, "upstream diagnostic")
		os.Exit(7)
	}
}
func TestUpstreamForwarding(t *testing.T) {
	binary, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	content, err := os.ReadFile(binary)
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	name := "hyfetch"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	if err := os.WriteFile(filepath.Join(dir, name), content, 0755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir)
	t.Setenv("TINTFETCH_ENGINE_TEST_HELPER", "1")
	var out, diagnostics bytes.Buffer
	args := []string{"--preset", "trans", "--unknown-upstream-option", "a value"}
	code := runUpstream(context.Background(), "hyfetch", args, strings.NewReader(""), &out, &diagnostics)
	if code != 7 {
		t.Fatalf("exit status not preserved: %d %s", code, diagnostics.String())
	}
	var forwarded []string
	if err := json.Unmarshal(out.Bytes(), &forwarded); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(forwarded, args) || !strings.Contains(diagnostics.String(), "upstream diagnostic") {
		t.Fatalf("forwarding changed: %v %s", forwarded, diagnostics.String())
	}
}
func TestMissingUpstream(t *testing.T) {
	t.Setenv("PATH", t.TempDir())
	var out, diagnostics bytes.Buffer
	code := runUpstream(context.Background(), "hyfetch", nil, strings.NewReader(""), &out, &diagnostics)
	if code != 1 || out.Len() != 0 || !strings.Contains(diagnostics.String(), "install hyfetch") {
		t.Fatalf("%d %s %s", code, out.String(), diagnostics.String())
	}
}
