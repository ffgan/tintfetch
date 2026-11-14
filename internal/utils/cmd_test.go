package utils

import (
	"os"
	"testing"
	"time"
)

func init() {
	if os.Getenv("TINTFETCH_TIMEOUT_TEST_HELPER") == "1" {
		time.Sleep(30 * time.Second)
		os.Exit(0)
	}
}
func TestCommandTimeout(t *testing.T) {
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("TINTFETCH_TIMEOUT_TEST_HELPER", "1")
	started := time.Now()
	if _, err := CommandOutput(executable); err == nil {
		t.Fatal("hung probe completed without a timeout")
	}
	if elapsed := time.Since(started); elapsed > 6*time.Second {
		t.Fatalf("probe timeout was not bounded: %v", elapsed)
	}
}
