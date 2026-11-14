package configs

import (
	"context"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestConfigRoundTrip(t *testing.T) {
	original := DefaultConfig()
	original.Separator = "it's # fine"
	original.InfoList = [][]string{{"My operating system", "distro"}, {"title"}}
	loaded := DefaultConfig()
	if err := ParseConfig(Format(&original), &loaded); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(original, loaded) {
		t.Fatalf("round trip changed configuration")
	}
}
func TestDeclarativeConfig(t *testing.T) {
	c := DefaultConfig()
	err := ParseConfig("cpu_cores=physical\nseparator=' # ' # comment\ndisk_show=('/' '/home')\nprint_info() {\n info \"Operating System\" distro\n}\n", &c)
	if err != nil {
		t.Fatal(err)
	}
	if c.CPUCores != "physical" || c.Separator != " # " || c.DiskShow != "('/' '/home')" || !reflect.DeepEqual(c.InfoList, [][]string{{"Operating System", "distro"}}) {
		t.Fatalf("configuration not applied: %+v", c)
	}
}
func TestRejectShellAndInvalidConfig(t *testing.T) {
	for _, input := range []string{"unknown=x", "cpu_cores=$(touch /tmp/never)", "cpu_cores=\"`id`\"", "separator='broken", "bar_length=999999999", "print_info() {\n info title", "print_info() {\n echo hi\n}", "cpu_cores=x;id"} {
		t.Run(input, func(t *testing.T) {
			c := DefaultConfig()
			if err := ParseConfig(input, &c); err == nil {
				t.Fatal("invalid input accepted")
			}
		})
	}
}
func TestShellConfigOptIn(t *testing.T) {
	if _, err := exec.LookPath("bash"); err != nil {
		t.Skip("Bash not installed")
	}
	path := filepath.Join(t.TempDir(), "config.conf")
	if err := os.WriteFile(path, []byte("cpu_cores=$(printf physical)\ndisk_show=('/mnt/with spaces')\nprint_info() { info 'Custom OS' distro; }\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadConfig(path); err == nil {
		t.Fatal("ordinary parser executed shell config")
	}
	c, err := LoadShellConfig(context.Background(), path, io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	words, wordErr := Words(strings.Trim(c.DiskShow, "()"))
	if wordErr != nil || !reflect.DeepEqual(words, []string{"/mnt/with spaces"}) {
		t.Fatalf("shell array boundaries lost: %q %v", c.DiskShow, wordErr)
	}
	if c.CPUCores != "physical" || !reflect.DeepEqual(c.InfoList, [][]string{{"Custom OS", "distro"}}) {
		t.Fatalf("bad shell output: %+v", c)
	}
}
func FuzzParseConfig(f *testing.F) {
	f.Add("separator=': '\nprint_info() {\n info title\n}\n")
	f.Add("print_info() {\n")
	f.Fuzz(func(t *testing.T, input string) {
		if len(input) > 65536 {
			t.Skip()
		}
		c := DefaultConfig()
		_ = ParseConfig(strings.ReplaceAll(input, "\x00", ""), &c)
	})
}
