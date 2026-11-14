package configs

import (
	"context"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

func TestExplicitShellConfiguration(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Bash integration is exercised on Unix CI")
	}
	if _, err := exec.LookPath("bash"); err != nil {
		t.Skip("Bash is not installed")
	}
	path := filepath.Join(t.TempDir(), "trusted.conf")
	if err := os.WriteFile(path, []byte("cpu_cores=$(printf physical)\nprint_info() { info 'My OS' distro; }\n"), 0600); err != nil {
		t.Fatal(err)
	}
	config, err := LoadShellConfig(context.Background(), path, io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	if config.CPUCores != "physical" || !reflect.DeepEqual(config.InfoList, [][]string{{"My OS", "distro"}}) {
		t.Fatalf("shell configuration was not captured: %+v", config)
	}
}
