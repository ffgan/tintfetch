package cli

import (
	"os"
	"strings"

	"github.com/ffgan/tintfetch/internal/utils"
)

func getHostname() string {
	if h, err := os.Hostname(); err == nil {
		return h
	}
	if b, err := utils.CommandOutput("hostname"); err == nil {
		return strings.TrimSpace(string(b))
	}
	return "unknown"
}
