package configs

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// LoadShellConfig intentionally executes trusted user code. Callers must require
// an explicit opt-in and filename; ordinary configuration never uses Bash.
func LoadShellConfig(parent context.Context, path string, diagnostics io.Writer) (*Config, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	config := DefaultConfig()
	fields := Fields(&config)
	keys := make([]string, 0, len(fields))
	for key := range fields {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var script strings.Builder
	script.WriteString("set -e\n")
	for _, key := range keys {
		fmt.Fprintf(&script, "%s=%s\n", key, shellQuote(*fields[key]))
	}
	script.WriteString("source \"$1\" >&2\n")
	for _, key := range keys {
		if key == "disk_show" || key == "local_ip_interface" {
			fmt.Fprintf(&script, "if [[ $(declare -p %s) == 'declare -a '* ]]; then printf -v tintfetch_value '%%q ' \"${%s[@]}\"; tintfetch_value=\"(${tintfetch_value})\"; else tintfetch_value=\"${%s}\"; fi\n", key, key, key)
			fmt.Fprintf(&script, "printf 'F\\0%%s\\0%%s\\0' %s \"$tintfetch_value\"\n", shellQuote(key))
		} else {
			fmt.Fprintf(&script, "printf 'F\\0%%s\\0%%s\\0' %s \"${%s[*]}\"\n", shellQuote(key), key)
		}
	}
	script.WriteString("info() { if (( $# == 1 )); then printf 'I\\0\\0%s\\0' \"$1\"; elif (( $# == 2 )); then printf 'I\\0%s\\0%s\\0' \"$1\" \"$2\"; else return 1; fi; }\n")
	script.WriteString("if declare -F print_info >/dev/null; then printf 'R\\0\\0\\0'; print_info; fi\n")
	ctx, cancel := context.WithTimeout(parent, 5*time.Second)
	defer cancel()
	command := exec.CommandContext(ctx, "bash", "--noprofile", "--norc", "-c", script.String(), "tintfetch-config", filepath.ToSlash(path))
	for _, entry := range os.Environ() {
		if !strings.HasPrefix(entry, "BASH_ENV=") && !strings.HasPrefix(entry, "ENV=") {
			command.Env = append(command.Env, entry)
		}
	}
	var output bytes.Buffer
	command.Stdout = &output
	command.Stderr = diagnostics
	command.WaitDelay = 100 * time.Millisecond
	if err := command.Run(); err != nil {
		return nil, fmt.Errorf("execute shell config: %w", err)
	}
	parts := bytes.Split(output.Bytes(), []byte{0})
	if len(parts) == 0 || len(parts[len(parts)-1]) != 0 || (len(parts)-1)%3 != 0 {
		return nil, fmt.Errorf("invalid shell config output")
	}
	for i := 0; i < len(parts)-1; i += 3 {
		kind, key, value := string(parts[i]), string(parts[i+1]), string(parts[i+2])
		switch kind {
		case "F":
			field, ok := fields[key]
			if !ok {
				return nil, fmt.Errorf("unknown shell config key %q", key)
			}
			*field = value
		case "R":
			config.InfoList = nil
		case "I":
			if key == "" {
				config.InfoList = append(config.InfoList, []string{value})
			} else {
				config.InfoList = append(config.InfoList, []string{key, value})
			}
		default:
			return nil, fmt.Errorf("unsupported output from shell print_info; use only info calls")
		}
	}
	return &config, Validate(&config)
}
