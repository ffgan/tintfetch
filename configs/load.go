package configs

import (
	"fmt"
	"os"

	"github.com/ffgan/tintfetch/internal/utils"
)

func LoadConfig(fileName string) (*Config, error) {
	config := DefaultConfig()
	if fileName == "" {
		return &config, nil
	}

	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exist, %w", err)
	}

	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	err = ParseConfig(string(data), &config)
	if err != nil {
		return nil, fmt.Errorf("parse config file: %w", err)
	}

	return &config, nil
}

func DefaultConfig() Config {
	return Config{
		Title: Title{
			TitleFqdn: utils.OFF,
		},
		Kernel: Kernel{
			KernelShorthand: utils.ON,
		},
		Distro: Distro{
			DistroShorthand: utils.OFF,
			OSArch:          utils.ON,
		},
		Uptime: Uptime{
			UptimeShorthand: utils.ON,
		},
		Memory: Memory{
			MemoryPercent: utils.ON,
			MemoryUnit:    "gib",
			MemPrecision:  "2",
		},
		Packages: Packages{
			PackageManagers: utils.ON,
			PackageSeparate: utils.ON,
			PackageMinimal:  "",
		},
		Shell: Shell{
			ShellPath:    utils.OFF,
			ShellVersion: utils.ON,
		},
		Editor: Editor{
			EditorPath:    utils.OFF,
			EditorVersion: utils.ON,
		},
		CPU: CPU{
			SpeedType:      "bios_limit",
			SpeedShorthand: utils.ON,
			CPUBrand:       utils.ON,
			CPUSpeed:       utils.ON,
			CPUCores:       "logical",
			CPUTemp:        utils.OFF,
		},
		GPU: GPU{
			GPUBrand: utils.ON,
			GPUType:  "all",
		},
		Resolution: Resolution{
			RefreshRate: utils.ON,
		},
		GTK: GTK{
			GTKShorthand: utils.OFF,
			GTK2:         utils.ON,
			GTK3:         utils.ON,
			QT:           utils.ON,
		},
		IPAddress: IPAddress{
			PublicIPHost:     "https://ident.me",
			PublicIPTimeout:  "2",
			LocalIPInterface: "('auto')",
		},
		DE: DE{
			DeVersion: utils.ON,
		},
		Disk: Disk{
			DiskShow:     "('/')",
			DiskSubtitle: "mount",
			DiskPercent:  utils.ON,
		},
		Song: Song{
			MusicPlayer:   utils.AUTO,
			SongFormat:    "%artist% - %album% - %title%",
			SongShorthand: utils.OFF,
			MPCArgs:       "()",
		},
		TextColor: TextColor{
			Colors: "(distro)",
		},
		TextOptions: TextOptions{
			Bold:             utils.ON,
			UnderlineEnabled: utils.ON,
			UnderlineChar:    "-",
			Separator:        ":",
		},
		ColorBlocks: ColorBlocks{
			BlockRange:  "(0 15)",
			ColorBlocks: utils.ON,
			BlockWidth:  "3",
			BlockHeight: "1",
			ColOffset:   utils.AUTO,
		},
		ProgressBar: ProgressBar{
			BarCharElapsed:  "-",
			BarCharTotal:    "=",
			BarBorder:       utils.ON,
			BarLength:       "15",
			BarColorElapsed: "distro",
			BarColorTotal:   "distro",
			MemoryDisplay:   utils.OFF,
			BatteryDisplay:  utils.OFF,
			DiskDisplay:     utils.OFF,
		},
		BackendSettings: BackendSettings{
			ImageBackend: "ascii",
			ImageSource:  utils.AUTO,
		},
		ASCIIOptions: ASCIIOptions{
			ASCIIDistro: utils.AUTO,
			ASCIIColors: "(distro)",
			ASCIIBold:   utils.ON,
		},
		ImageOptions: ImageOptions{
			ImageLoop:       utils.OFF,
			ThumbnailDir:    "${XDG_CACHE_HOME:-${HOME}/.cache}/thumbnails/tintfetch",
			CropMode:        "normal",
			CropOffset:      "center",
			ImageSize:       utils.AUTO,
			CatimgSize:      "2",
			Gap:             "3",
			Yoffset:         "0",
			Xoffset:         "0",
			BackgroundColor: "",
		},
		MiscOptions: MiscOptions{
			Stdout: utils.AUTO,
		},
		PrintInfo: PrintInfo{
			InfoList: [][]string{
				{"title"},
				{"underline"},
				{"OS", "distro"},
				{"Host", "model"},
				{"Kernel", "kernel"},
				{"Uptime", "uptime"},
				{"Packages", "packages"},
				{"Shell", "shell"},
				{"Editor", "editor"},
				{"Resolution", "resolution"},
				{"DE", "de"},
				{"WM", "wm"},
				{"WM Theme", "wm_theme"},
				{"Theme", "theme"},
				{"Icons", "icons"},
				{"Cursor", "cursor"},
				{"Terminal", "term"},
				{"Terminal Font", "term_font"},
				{"CPU", "cpu"},
				{"GPU", "gpu"},
				{"Memory", "memory"},
				{"Network", "network"},
				{"Bluetooth", "bluetooth"},
				{"BIOS", "bios"},
				{"cols"},
			},
		},
	}
}
