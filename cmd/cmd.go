package cmd

import (
	"context"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"bufio"
	"time"

	"github.com/ffgan/tintfetch/configs"
	cli "github.com/ffgan/tintfetch/internal/CLI"
	"github.com/ffgan/tintfetch/internal/hyfetch"
	"github.com/ffgan/tintfetch/internal/logo"
	"github.com/ffgan/tintfetch/internal/utils"
)

const neofetchConf = "/etc/neofetch/config.conf"

func Run(ASCIIFiles embed.FS) {
	for _, arg := range os.Args[1:] {
		if arg == "--list-presets" {
			for _, name := range hyfetch.PresetNames(nil) {
				fmt.Println(name)
			}
			return
		}
	}
	for i, arg := range os.Args[1:] {
		switch arg {
		case "--plain":
			os.Args[i+1] = "--stdout"
		case "--no-config":
			os.Args[i+1] = "--no_config"
		case "--print-config":
			os.Args[i+1] = "--print_config"
		}
	}
	shellConfig := false
	filtered := os.Args[:1]
	for _, arg := range os.Args[1:] {
		if arg == "--shell-config" {
			shellConfig = true
		} else {
			filtered = append(filtered, arg)
		}
	}
	os.Args = filtered
	hyOpts, rest, err := hyfetch.SplitArgs(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "tintfetch: %v\n", err)
		os.Exit(2)
	}
	if hyOpts.Version {
		fmt.Printf("tintfetch %s\n", Version)
		os.Exit(0)
	}
	os.Args = append([]string{os.Args[0]}, rest...)

	kernelName, kernelVersion, kernelMachine, darwinName := cacheUname()
	osName := cli.GetOS(kernelName, darwinName)

	hyPath := hyOpts.ConfigFile
	if hyPath == "" {
		hyPath = hyfetch.ConfigPath()
	}
	var hyCfg hyfetch.Config
	hyExists := false
	if hyOpts.Configure {
		hyCfg, err = hyfetch.Configure(hyPath, os.Stdin, os.Stdout)
		if err != nil {
			fmt.Fprintf(os.Stderr, "tintfetch: %v\n", err)
			os.Exit(2)
		}
		hyExists = true
	} else {
		hyCfg, hyExists, err = hyfetch.Load(hyPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "tintfetch: %v\n", err)
			os.Exit(1)
		}
	}
	hyCfg, err = hyOpts.Apply(hyCfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "tintfetch: %v\n", err)
		os.Exit(2)
	}
	colorOn := hyExists || hyOpts.WantsColor()
	if colorOn && (hyCfg.Backend == "" || hyCfg.Backend == "neofetch") && len(hyCfg.Args) > 0 {
		os.Args = append(os.Args, hyCfg.Args...)
	}

	explicit, skipFile := configFromArgs(os.Args[1:])
	var filePath string
	switch {
	case explicit != "":
		filePath = explicit
	case !skipFile:
		filePath = findConfig()
	}

	if shellConfig && (explicit == "" || skipFile) {
		fmt.Fprintln(os.Stderr, "tintfetch: --shell-config requires an explicit --config path")
		os.Exit(2)
	}
	config, err := configs.LoadConfig(filePath)
	if shellConfig {
		config, err = configs.LoadShellConfig(context.Background(), filePath, os.Stderr)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "tintfetch: %v\n", err)
		os.Exit(2)
	}

	disabled := parseArgs(config)

	h := &hardware{
		cache_dir:      getCacheDir(),
		KernelName:     kernelName,
		kernel_version: kernelVersion,
		KernelMachine:  kernelMachine,
		DarwinName:     darwinName,
		os:             osName,
		Config:         *config,
		verbose:        verbose,
		disabled:       disabled,
	}
	if h.Stdout == utils.ON {
		h.plain = true
	}

	logo.SetAppearance(h.ASCIIBold != utils.OFF, h.ASCIIColors, h.Colors)
	if colorOn && hyCfg.Distro != "" {
		h.ASCIIDistro = hyCfg.Distro
	}
	h.Get_distro()

	if hyOpts.PrintFontLogo {
		glyph, err := hyfetch.FontLogo(h.distro)
		if err != nil {
			fmt.Fprintf(os.Stderr, "tintfetch: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(glyph)
		return
	}
	if hyOpts.TestPrint {
		fmt.Println(logo.RawASCII(h.ascii_distro, h.ImageSource, ASCIIFiles))
		return
	}

	theme := hyfetch.ResolveTheme(hyCfg, hyOpts)
	now := time.Now()
	if hyfetch.ShouldPride(hyCfg, hyOpts.June, now) {
		mode := hyCfg.Mode
		if mode == "" {
			mode = hyfetch.Mode8bit
		}
		hyfetch.Pride(mode)
		if !hyOpts.June {
			hyfetch.MarkPrideShown(now)
		}
	}

	var art []string
	if colorOn && !h.plainMode() && !strings.EqualFold(h.ImageBackend, "off") {
		art = hyfetchArt(h, hyCfg, hyOpts, theme, ASCIIFiles)
	} else if !h.plainMode() && !strings.EqualFold(h.ImageBackend, "off") {
		if lines := imageArt(h.ImageSource); len(lines) > 0 {
			art = lines
			logo.ApplyDistro(h.ascii_distro)
		} else {
			art = logo.PrintASCII(h.ascii_distro, h.ImageSource, ASCIIFiles)
		}
	} else {
		logo.ApplyDistro(h.ascii_distro)
	}

	if jsonMode {
		h.writeJSON(os.Stdout)
		return
	}
	show(art, h.Rows(), h.gap())
	if hyOpts.AskExit {
		fmt.Fprint(os.Stderr, "Press enter to exit...")
		_, _ = bufio.NewReader(os.Stdin).ReadString('\n')
	}
}

func hyfetchArt(h *hardware, cfg hyfetch.Config, opt hyfetch.Options, theme hyfetch.Theme, files embed.FS) []string {
	source := h.ImageSource
	if cfg.CustomASCIIPath != "" {
		source = cfg.CustomASCIIPath
	}
	raw := logo.RawASCII(h.ascii_distro, source, files)
	lines, profile, name, err := hyfetch.Paint(raw, cfg, opt, theme)
	if err != nil {
		fmt.Fprintf(os.Stderr, "tintfetch: %v\n", err)
		return logo.PrintASCII(h.ascii_distro, h.ImageSource, files)
	}
	if opt.Debug {
		fmt.Fprintf(os.Stderr, "tintfetch: preset %s, mode %s, align %s, colors %d\n", name, cfg.Mode, cfg.ColorAlign.Mode, len(profile.Colors))
	}
	art := strings.Join(lines, "\n")
	switch cfg.Backend {
	case "fastfetch":
		if err := hyfetch.RunFastfetch(art, []string(cfg.Args)); err != nil {
			fmt.Fprintf(os.Stderr, "tintfetch: %v\n", err)
			break
		}
		os.Exit(0)
	case "macchina":
		if err := hyfetch.RunMacchina(art, []string(cfg.Args)); err != nil {
			fmt.Fprintf(os.Stderr, "tintfetch: %v\n", err)
			break
		}
		os.Exit(0)
	}
	return lines
}

func findConfig() string {
	var candidates []string
	if dir, err := os.UserConfigDir(); err == nil {
		candidates = append(candidates, filepath.Join(dir, "tintfetch", "config.conf"))
	}
	if home := os.Getenv("HOME"); home != "" {
		candidates = append(candidates,
			filepath.Join(home, ".config", "tintfetch", "config.conf"),
			filepath.Join(home, ".config", "neofetch", "config.conf"),
		)
	}
	candidates = append(candidates, "/etc/tintfetch/config.conf", neofetchConf, "/etc/neofetch/default.conf")
	for _, path := range candidates {
		if utils.FileExists(path) {
			return path
		}
	}
	return ""
}

func configFromArgs(args []string) (path string, skip bool) {
	for i := 0; i < len(args); i++ {
		switch {
		case args[i] == "--no_config" || args[i] == "--noconfig":
			skip = true
		case args[i] == "--config" && i+1 < len(args):
			i++
			path = args[i]
		case strings.HasPrefix(args[i], "--config="):
			path = strings.TrimPrefix(args[i], "--config=")
		}
	}
	return path, skip
}

func getCacheDir() string {
	cacheDir := os.Getenv("TMPDIR")
	if cacheDir == "" {
		cacheDir = "/tmp"
	}
	return cacheDir
}

func cacheUname() (string, string, string, string) {
	parts := strings.Fields(cli.UName("-srm"))
	for len(parts) < 3 {
		parts = append(parts, "")
	}
	return parts[0], parts[1], parts[2], ""
}

func imageArt(source string) []string {
	fi, err := os.Stat(source)
	if err != nil || !fi.Mode().IsRegular() {
		return nil
	}
	switch strings.ToLower(filepath.Ext(source)) {
	case ".png", ".jpg", ".jpeg", ".jpe", ".gif", ".webp":
	default:
		return nil
	}
	tools := []struct {
		bin  string
		args []string
	}{
		{"chafa", []string{"--symbols", "block", "--size", "40x40", source}},
		{"catimg", []string{"-w", "40", source}},
		{"viu", []string{"-w", "40", source}},
	}
	for _, tool := range tools {
		if !utils.CommandExists(tool.bin) {
			continue
		}
		out := utils.RunCommand(tool.bin, tool.args...)
		if strings.TrimSpace(out) == "" {
			continue
		}
		return strings.Split(strings.TrimRight(out, "\n"), "\n")
	}
	return nil
}
