// Package logo renders the embedded distribution artwork without global state.
package logo

import (
	"embed"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

const Reset = "\x1b[0m"

// Render selects artwork, resolves its color placeholders, and always resets
// styling at line boundaries. Unknown distributions fall back to Linux artwork.
func Render(distro, source string, bold bool, override []int, files embed.FS) ([]string, error) {
	_, value, ok := r.LongestPrefix(distro)
	if !ok {
		_, value, _ = r.LongestPrefix("Linux")
	}
	definition := value.(DistroWithColor)
	var data []byte
	var err error
	if source == "" || source == "auto" || source == "ascii" {
		data, err = files.ReadFile("assets/" + definition.FileName)
		if err != nil {
			_, fallback, _ := r.LongestPrefix("Linux")
			definition = fallback.(DistroWithColor)
			data, err = files.ReadFile("assets/" + definition.FileName)
		}
	} else if stat, statErr := os.Stat(source); statErr == nil && stat.Mode().IsRegular() {
		data, err = os.ReadFile(source)
	} else {
		data = []byte(source)
	}
	if err != nil {
		return nil, fmt.Errorf("read ASCII artwork: %w", err)
	}
	if !utf8.Valid(data) || strings.ContainsRune(string(data), 0) {
		return nil, fmt.Errorf("ASCII source must be UTF-8 text; image files require an upstream image backend")
	}
	colors := definition.Colors
	if len(override) > 0 {
		colors = override
	}
	text := string(data)
	for i := 1; i <= 6; i++ {
		code := Reset
		if i <= len(colors) {
			code = color(colors[i-1])
		}
		if bold {
			code += "\x1b[1m"
		}
		text = strings.ReplaceAll(text, fmt.Sprintf("${c%d}", i), code)
	}
	lines := strings.Split(strings.TrimRight(text, "\n"), "\n")
	for i := range lines {
		lines[i] += Reset
	}
	return lines, nil
}
func color(index int) string {
	if index < 0 || index > 255 {
		return Reset
	}
	return "\x1b[38;5;" + strconv.Itoa(index) + "m"
}

// PrintASCII retains the former internal API for callers using default styling.
func PrintASCII(distro, source string, files embed.FS) []string {
	lines, _ := Render(distro, source, true, nil, files)
	return lines
}
