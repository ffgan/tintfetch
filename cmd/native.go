package cmd

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/ffgan/tintfetch/configs"
)

var nativeOptions = map[string]bool{
	"package_managers": true, "package_separate": true, "package_minimal": true, "speed_type": true,
	"title_fqdn": true, "os_arch": true, "distro_shorthand": true, "kernel_shorthand": true,
	"uptime_shorthand": true, "shell_path": true, "shell_version": true, "editor_path": true, "editor_version": true,
	"cpu_cores": true, "cpu_speed": true, "speed_shorthand": true, "cpu_brand": true, "cpu_temp": true,
	"memory_percent": true, "memory_unit": true, "mem_precision": true, "bar_length": true, "bar_char_elapsed": true, "bar_char_total": true, "memory_display": true,
	"image_backend": true, "image_source": true, "ascii_distro": true, "ascii_colors": true, "ascii_bold": true, "gap": true, "separator": true, "stdout": true,
	"public_ip_host": true, "public_ip_timeout": true, "local_ip_interface": true, "disk_show": true, "disk_subtitle": true, "disk_percent": true,
	"underline_enabled": true, "underline_char": true,
}

func validateNative(fs *flag.FlagSet, c *configs.Config) error {
	defaults := configs.DefaultConfig()
	baseline := configs.Fields(&defaults)
	for key, p := range configs.Fields(c) {
		if !nativeOptions[key] && *p != *baseline[key] {
			return fmt.Errorf("native engine does not yet implement %s; use --engine neofetch", key)
		}
	}
	var unsupported string
	fs.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "gpu_brand", "gpu_type", "refresh_rate", "de_version", "gtk_shorthand", "gtk2", "gtk3", "qt", "song_format", "song_shorthand", "music_player", "bold", "color_blocks", "block_range", "block_width", "block_height", "col_offset", "bar_char", "bar_border", "bar_colors", "battery_display", "disk_display", "loop", "image_size", "catimg_size", "crop_mode", "crop_offset", "xoffset", "yoffset", "background_color", "bg_color", "v", "vv":
			unsupported = f.Name
		}
	})
	if unsupported != "" {
		return fmt.Errorf("native engine does not yet implement --%s; use --engine neofetch", unsupported)
	}
	if c.ImageBackend != "ascii" && c.ImageBackend != "off" {
		return fmt.Errorf("image backend %q requires --engine neofetch", c.ImageBackend)
	}
	enums := map[string][]string{
		"package_managers": {"on", "tiny", "off"}, "package_minimal": {"", "on", "off", "1"},
		"stdout": {"on", "off", "auto"}, "cpu_cores": {"logical", "physical", "off"}, "cpu_temp": {"on", "off", "C", "F"},
		"memory_unit": {"kib", "mib", "gib", "tib"}, "memory_display": {"off", "bar", "infobar", "barinfo"},
		"uptime_shorthand": {"on", "off", "tiny"}, "distro_shorthand": {"on", "off", "tiny"},
	}
	fields := configs.Fields(c)
	for key, values := range enums {
		valid := false
		for _, value := range values {
			if *fields[key] == value {
				valid = true
			}
		}
		if !valid {
			return fmt.Errorf("invalid %s: %q (expected %s)", key, *fields[key], strings.Join(values, ", "))
		}
	}
	for _, key := range []string{"package_separate", "title_fqdn", "os_arch", "kernel_shorthand", "shell_path", "shell_version", "editor_path", "editor_version", "cpu_speed", "speed_shorthand", "cpu_brand", "memory_percent", "disk_percent", "underline_enabled", "ascii_bold"} {
		if *fields[key] != "on" && *fields[key] != "off" {
			return fmt.Errorf("%s must be on or off", key)
		}
	}
	if c.SpeedType == "" || strings.Trim(c.SpeedType, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_") != "" {
		return fmt.Errorf("speed_type must be a cpufreq attribute name")
	}
	if c.ASCIIColors != "(distro)" && c.ASCIIColors != "distro" {
		words, err := configs.Words(strings.Trim(c.ASCIIColors, "()"))
		if err != nil || len(words) == 0 || len(words) > 6 {
			return fmt.Errorf("ascii_colors needs 1..6 palette indexes")
		}
		for _, word := range words {
			n, err := strconv.Atoi(word)
			if err != nil || n < 0 || n > 255 {
				return fmt.Errorf("ascii_colors indexes must be 0..255")
			}
		}
	}
	for _, row := range c.InfoList {
		if !nativeModule(row[len(row)-1]) {
			return fmt.Errorf("module %q requires an upstream engine or is unknown", row[len(row)-1])
		}
	}
	return nil
}

func nativeModule(name string) bool {
	switch name {
	case "title", "underline", "distro", "os", "bios", "model", "kernel", "uptime", "packages", "shell", "editor", "de", "wm", "wm_theme", "cpu", "gpu", "memory", "network", "bluetooth", "resolution", "theme", "icons", "font", "cursor", "term", "term_font", "cols", "disk", "battery", "power_adapter", "local_ip", "public_ip", "users", "locale", "java_ver", "python_ver", "node_ver", "gpu_driver":
		return true
	}
	return false
}

func collect(h *hardware, name string) (string, error) {
	switch name {
	case "cols":
		return "", nil
	case "font":
		h.Get_font()
		return h.term_font, nil
	case "underline":
		if h.UnderlineEnabled == "off" {
			return "", nil
		}
		h.Get_title()
		return strings.Repeat(h.UnderlineChar, min(h.length, 200)), nil
	case "disk":
		return diskInfo(h.Config), nil
	case "battery":
		return batteryInfo(), nil
	case "power_adapter":
		return adapterInfo(), nil
	case "local_ip":
		return localIP(h.LocalIPInterface), nil
	case "public_ip":
		return publicIP(h.PublicIPHost, h.PublicIPTimeout), nil
	case "users":
		return loggedUsers(), nil
	case "locale":
		for _, key := range []string{"LC_ALL", "LC_MESSAGES", "LANG"} {
			if value := os.Getenv(key); value != "" {
				return value, nil
			}
		}
		return "", nil
	case "java_ver":
		return toolVersion("java", "-version"), nil
	case "python_ver":
		return toolVersion("python3", "--version"), nil
	case "node_ver":
		return toolVersion("node", "--version"), nil
	case "gpu_driver":
		return gpuDriver(), nil
	}
	if !nativeModule(name) {
		return "", fmt.Errorf("unknown module %q", name)
	}
	method := reflect.ValueOf(h).MethodByName("Get_" + name)
	if !method.IsValid() {
		return "", fmt.Errorf("module %q is unavailable", name)
	}
	method.Call(nil)
	field := reflect.ValueOf(h).Elem().FieldByName(name)
	if !field.IsValid() || field.Kind() != reflect.String {
		return "", fmt.Errorf("module %q has no value", name)
	}
	return field.String(), nil
}
func atoiDefault(s string, fallback int) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		return fallback
	}
	return n
}

var presets = map[string][]string{
	"rainbow":   {"e40303", "ff8c00", "ffed00", "008026", "24408e", "732982"},
	"trans":     {"5bcefa", "f5a9b8", "ffffff", "f5a9b8", "5bcefa"},
	"bi":        {"d60270", "d60270", "9b4f96", "0038a8", "0038a8"},
	"pan":       {"ff218c", "ffd800", "21b1ff"},
	"nonbinary": {"fff430", "ffffff", "9c59d1", "2d2d2d"},
	"lesbian":   {"d52d00", "ef7627", "ff9a56", "ffffff", "d162a4", "b55690", "a30262"},
}

func tintArtwork(lines []string, palette []string) []string {
	out := make([]string, len(lines))
	height := len(lines)
	for height > 0 && stripANSI(lines[height-1]) == "" {
		height--
	}
	for i, line := range lines {
		index := min(i*len(palette)/max(height, 1), len(palette)-1)
		rgb, _ := strconv.ParseUint(palette[index], 16, 32)
		out[i] = fmt.Sprintf("\x1b[38;2;%d;%d;%dm%s\x1b[0m", rgb>>16, (rgb>>8)&255, rgb&255, stripANSI(line))
	}
	return out
}
