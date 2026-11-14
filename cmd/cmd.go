// Package cmd implements the Tintfetch command line.
package cmd

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/ffgan/tintfetch/configs"
	cli "github.com/ffgan/tintfetch/internal/CLI"
	"github.com/ffgan/tintfetch/internal/logo"
)

// Execute runs an invocation without terminating the host process.
func Execute(ctx context.Context, args []string, input io.Reader, output, diagnostics io.Writer, assets embed.FS) int {
	if engine, rest, err := engineArgs(args); err != nil {
		fmt.Fprintln(diagnostics, "tintfetch:", err)
		return 2
	} else if engine != "native" {
		return runUpstream(ctx, engine, rest, input, output, diagnostics)
	}
	config := configs.DefaultConfig()
	fs := flag.NewFlagSet("tintfetch", flag.ContinueOnError)
	fs.SetOutput(diagnostics)
	var showVersion bool
	setupArgs(&config, fs, &showVersion)
	var configPath, preset string
	var jsonOutput, noConfig, printConfig, plain, shellConfig, listPresets bool
	fs.StringVar(&configPath, "config", "", "Read a configuration file")
	fs.BoolVar(&noConfig, "no-config", false, "Ignore user and system configuration")
	fs.BoolVar(&printConfig, "print-config", false, "Print the default declarative configuration")
	fs.BoolVar(&plain, "plain", false, "Print plain text without artwork or ANSI escapes")
	fs.BoolVar(&jsonOutput, "json", false, "Print structured JSON without ANSI escapes")
	fs.BoolVar(&shellConfig, "shell-config", false, "Execute the selected trusted config with Bash")
	fs.StringVar(&config.ASCIIDistro, "ascii_distro", config.ASCIIDistro, "Select distribution artwork")
	fs.StringVar(&config.ASCIIColors, "ascii_colors", config.ASCIIColors, "Space-separated palette indexes, 0..255")
	fs.StringVar(&config.ASCIIBold, "ascii_bold", config.ASCIIBold, "Bold artwork: on or off")
	fs.StringVar(&preset, "preset", "", "Color artwork with rainbow, trans, bi, pan, nonbinary, or lesbian stripes")
	fs.BoolVar(&listPresets, "list-presets", false, "List native color presets")
	fs.String("engine", "native", "Engine: native, neofetch, or hyfetch")
	fs.Usage = func() {
		fmt.Fprintln(diagnostics, "Tintfetch — system information in color\nUsage: tintfetch [options]\n\nNative options (upstream compatibility: --engine neofetch|hyfetch):")
		fs.PrintDefaults()
	}
	if err := fs.Parse(normalizeArgs(args)); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}
	fail := func(err error) int { fmt.Fprintln(diagnostics, "tintfetch:", err); return 2 }
	if fs.NArg() != 0 {
		return fail(fmt.Errorf("unexpected arguments: %s", strings.Join(fs.Args(), " ")))
	}
	if showVersion {
		fmt.Fprintln(output, "Tintfetch", Version)
		return 0
	}
	if listPresets {
		fmt.Fprintln(output, "bi\nlesbian\nnonbinary\npan\nrainbow\ntrans")
		return 0
	}
	if printConfig || flagSet(fs, "print_config") {
		defaults := configs.DefaultConfig()
		fmt.Fprint(output, configs.Format(&defaults))
		return 0
	}
	if _, ok := presets[preset]; preset != "" && !ok {
		return fail(fmt.Errorf("unknown preset %q; use --list-presets", preset))
	}
	if noConfig && configPath != "" {
		return fail(fmt.Errorf("--config and --no-config are mutually exclusive"))
	}
	if shellConfig && (configPath == "" || noConfig) {
		return fail(fmt.Errorf("--shell-config requires an explicit --config path"))
	}
	if !noConfig && !flagSet(fs, "no_config") {
		if configPath == "" {
			configPath = defaultConfigPath()
		}
		if configPath != "" {
			var loaded *configs.Config
			var err error
			if shellConfig {
				loaded, err = configs.LoadShellConfig(ctx, configPath, diagnostics)
			} else {
				loaded, err = configs.LoadConfig(configPath)
			}
			if err != nil {
				return fail(err)
			}
			// Explicit flags take precedence over configuration, including aliases.
			oldFields, newFields := configs.Fields(&config), configs.Fields(loaded)
			changed := map[string]bool{}
			// Copy only fields addressed by explicitly supplied flag values.
			fs.Visit(func(f *flag.Flag) {
				for key, p := range oldFields {
					if reflect.ValueOf(f.Value).Kind() == reflect.Ptr && reflect.ValueOf(f.Value).Pointer() == reflect.ValueOf(p).Pointer() {
						changed[key] = true
					}
				}
			})
			for key := range changed {
				*newFields[key] = *oldFields[key]
			}
			config = *loaded
		}
	}
	if err := configs.Validate(&config); err != nil {
		return fail(err)
	}
	if err := validateNative(fs, &config); err != nil {
		return fail(err)
	}
	kernel, version, machine, darwin := cache_uname()
	h := hardware{KernelName: kernel, kernel_version: version, KernelMachine: machine, DarwinName: darwin, os: cli.GetOS(kernel, darwin), Config: config}
	h.Get_distro()
	type entry struct {
		Label  string `json:"label"`
		Module string `json:"module"`
		Value  string `json:"value"`
	}
	entries := make([]entry, 0, len(config.InfoList))
	for _, row := range config.InfoList {
		module, label := row[len(row)-1], ""
		if len(row) == 2 {
			label = row[0]
		}
		if module == "cols" {
			continue
		}
		value, err := collect(&h, module)
		if err != nil {
			return fail(err)
		}
		if value != "" {
			entries = append(entries, entry{label, module, stripANSI(value)})
		}
	}
	if jsonOutput {
		err := json.NewEncoder(output).Encode(struct {
			Version string  `json:"version"`
			Entries []entry `json:"entries"`
		}{Version, entries})
		if err != nil {
			fmt.Fprintln(diagnostics, err)
			return 1
		}
		return 0
	}
	terminal := false
	if file, ok := output.(*os.File); ok {
		if stat, err := file.Stat(); err == nil {
			terminal = stat.Mode()&os.ModeCharDevice != 0
		}
	}
	color := terminal && os.Getenv("TERM") != "dumb" && os.Getenv("NO_COLOR") == "" && !plain
	if config.Stdout == "on" {
		plain = true
		color = false
	}
	if config.Stdout == "off" && !plain && os.Getenv("NO_COLOR") == "" {
		color = true
	}
	var art []string
	if color && config.ImageBackend != "off" {
		var override []int
		if config.ASCIIColors != "(distro)" && config.ASCIIColors != "distro" {
			words, _ := configs.Words(strings.Trim(config.ASCIIColors, "()"))
			for _, word := range words {
				n, _ := strconv.Atoi(word)
				override = append(override, n)
			}
		}
		var err error
		art, err = logo.Render(h.ascii_distro, config.ImageSource, config.ASCIIBold == "on", override, assets)
		if err != nil {
			fmt.Fprintln(diagnostics, "tintfetch:", err)
			return 1
		}
		if preset != "" {
			art = tintArtwork(art, presets[preset])
		}
	}
	width := 0
	for _, line := range art {
		if n := len([]rune(stripANSI(line))); n > width {
			width = n
		}
	}
	count := len(entries)
	if len(art) > count {
		count = len(art)
	}
	for i := 0; i < count; i++ {
		var left, right string
		if i < len(art) {
			left = art[i]
		}
		if i < len(entries) {
			e := entries[i]
			right = e.Value
			if e.Label != "" {
				label := e.Label
				if color {
					label = "\x1b[1;34m" + label + "\x1b[0m"
				}
				right = label + config.Separator + " " + right
			}
		}
		if len(art) > 0 {
			fmt.Fprint(output, left, strings.Repeat(" ", width-len([]rune(stripANSI(left)))+atoiDefault(config.Gap, 3)))
		}
		fmt.Fprintln(output, right)
	}
	return 0
}

func flagSet(fs *flag.FlagSet, name string) bool {
	found := false
	fs.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func defaultConfigPath() string {
	if dir, err := os.UserConfigDir(); err == nil {
		path := filepath.Join(dir, "tintfetch", "config.conf")
		if _, err := os.Stat(path); err == nil || !errors.Is(err, os.ErrNotExist) {
			return path
		}
	}
	if runtime.GOOS != "windows" {
		path := "/etc/tintfetch/config.conf"
		if _, err := os.Stat(path); err == nil || !errors.Is(err, os.ErrNotExist) {
			return path
		}
	}
	return ""
}

func normalizeArgs(args []string) []string {
	out := make([]string, 0, len(args))
	for _, arg := range args {
		switch arg {
		case "--no_config":
			arg = "--no-config"
		case "--print_config":
			arg = "--print-config"
		}
		out = append(out, arg)
	}
	return out
}

func engineArgs(args []string) (string, []string, error) {
	engine := "native"
	rest := []string{}
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "--" {
			rest = append(rest, args[i:]...)
			break
		}
		if arg == "--engine" || arg == "-engine" {
			if i+1 >= len(args) {
				return "", nil, fmt.Errorf("--engine requires a value")
			}
			i++
			engine = args[i]
		} else if strings.HasPrefix(arg, "--engine=") || strings.HasPrefix(arg, "-engine=") {
			_, engine, _ = strings.Cut(arg, "=")
		} else {
			rest = append(rest, arg)
		}
	}
	if engine != "native" && engine != "neofetch" && engine != "hyfetch" {
		return "", nil, fmt.Errorf("unknown engine %q", engine)
	}
	return engine, rest, nil
}

func runUpstream(ctx context.Context, engine string, args []string, input io.Reader, output, diagnostics io.Writer) int {
	name := engine
	if engine == "neofetch" {
		if _, err := exec.LookPath("neowofetch"); err == nil {
			name = "neowofetch"
		}
	}
	path, err := exec.LookPath(name)
	if err != nil {
		fmt.Fprintf(diagnostics, "tintfetch: install %s to use --engine %s; native is the built-in Go engine\n", engine, engine)
		return 1
	}
	command := exec.CommandContext(ctx, path, args...)
	command.Stdin = input
	command.Stdout = output
	command.Stderr = diagnostics
	if err := command.Run(); err != nil {
		var exit *exec.ExitError
		if errors.As(err, &exit) {
			code := exit.ExitCode()
			if code > 0 {
				return code
			}
			return 1
		}
		fmt.Fprintln(diagnostics, "tintfetch:", err)
		return 1
	}
	return 0
}

var ansiPattern = regexp.MustCompile(`\x1b\[[0-?]*[ -/]*[@-~]`)

func stripANSI(s string) string { return ansiPattern.ReplaceAllString(s, "") }
