package cli

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ffgan/tintfetch/internal/utils"
)

var programVersionPattern = regexp.MustCompile(`[0-9]+(?:\.[0-9]+)+(?:[a-zA-Z0-9.-]*)`)

func GetShell(shellPath, shellVersion string) string {
	executable := os.Getenv("SHELL")
	if executable == "" {
		return ""
	}
	name := filepath.Base(executable)
	display := name
	if shellPath == utils.ON {
		display = executable
	}
	if shellVersion != utils.ON {
		return display
	}
	var version string
	switch name {
	case "sh", "ash", "dash", "es":
		return display
	case "bash":
		version = os.Getenv("BASH_VERSION")
	case "osh":
		version = os.Getenv("OIL_VERSION")
	}
	if version == "" {
		out, err := utils.CommandCombinedOutput(executable, "--version")
		if err == nil {
			version = strings.Split(string(out), "\n")[0]
		}
	}
	if parsed := programVersionPattern.FindString(version); parsed != "" {
		display += " " + parsed
	}
	return display
}
