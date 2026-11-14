package cli

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ffgan/tintfetch/configs"
	"github.com/ffgan/tintfetch/internal/utils"
)

type PkgInfo struct {
	Manager string
	Count   int
}

func DetectPackages(osname string) string {
	return DetectPackagesWithOptions(osname, configs.DefaultConfig().Packages)
}

func DetectPackagesWithOptions(osname string, options configs.Packages) string {
	// TODO: 修复与hyfetch不一致的地方
	// Packages: 227 (pip), 2699 (rpm), 25 (flatpak)
	// Packages: 227 (pip), 2699 (rpm), 19 (flatpak-system), 6 (flatpak-user)
	managers := []PkgInfo{}
	add := func(manager string, count int) {
		if count > 0 {
			managers = append(managers, PkgInfo{manager, count})
		}
	}

	// --- Language package managers ---
	if options.PackageMinimal != "on" && options.PackageMinimal != "1" {
		if utils.CommandExists("pipx") {
			add("pipx", utils.GetPkgCount("pipx", "list", "--short"))
		}
		if utils.CommandExists("pip") {
			add("pip", utils.GetPkgCount("pip", "freeze"))
		}
		if utils.CommandExists("cargo") {
			out, _ := utils.CommandOutput("cargo", "install", "--list")
			lines := 0
			sc := bufio.NewScanner(bytes.NewReader(out))
			for sc.Scan() {
				if !strings.HasPrefix(sc.Text(), " ") {
					lines++
				}
			}
			add("cargo", lines)
		}
		if utils.CommandExists("npm") {
			if st, err := os.Stat("/usr/lib/node_modules"); err == nil && st.IsDir() {
				add("npm", utils.FilePathGlobLens("/usr/lib/node_modules/*/"))
			} else if st, err := os.Stat("/usr/local/lib/node_modules"); err == nil && st.IsDir() {
				add("npm", utils.FilePathGlobLens("/usr/local/lib/node_modules/*/"))
			} else {
				out, err := utils.CommandOutput("npm", "root", "-g")
				if err == nil {
					root := strings.TrimSpace(string(out))
					add("npm", utils.FilePathGlobLens(filepath.Join(root, "*/")))
				}
			}
		}
		if utils.CommandExists("pnpm") {
			p := filepath.Join(os.Getenv("HOME"), ".local/share/pnpm/global/5/node_modules/*/")
			add("pnpm", utils.FilePathGlobLens(p))
		}

	}
	// --- System package managers by OS ---
	switch osname {
	case Linux:
		if utils.CommandExists("dpkg") {
			add("dpkg", utils.GetPkgCount("dpkg-query", "-f", ".\\n", "-W"))
		}
		if utils.CommandExists("pacman") {
			add("pacman", utils.GetPkgCount("pacman", "-Qq", "--color", "never"))
		}
		if utils.CommandExists("apk") {
			add("apk", utils.GetPkgCount("apk", "info"))
		}
		if utils.CommandExists("rpm") {
			add("rpm", utils.GetPkgCount("rpm", "-qa", "--nodigest", "--nosignature"))
		}
		if utils.CommandExists("flatpak") {
			if options.PackageSeparate == "on" {
				add("flatpak-system", utils.GetPkgCount("flatpak", "list", "--system", "--columns=application"))
				add("flatpak-user", utils.GetPkgCount("flatpak", "list", "--user", "--columns=application"))
			} else {
				add("flatpak", utils.GetPkgCount("flatpak", "list", "--columns=application"))
			}
		}
		if utils.CommandExists("snap") {
			add("snap", utils.GetPkgCount("snap", "list"))
		}
		if utils.CommandExists("brew") {
			root, err := utils.CommandOutput("brew", "--cellar")
			if err == nil {
				p1 := strings.TrimSpace(string(root)) + "/*"
				add("brew", utils.FilePathGlobLens(p1))
			}
		}
	case Darwin:
		if utils.CommandExists("brew") {
			add("brew", utils.FilePathGlobLens("/usr/local/Cellar/*/"))
		}
		if utils.CommandExists("port") {
			add("macports", utils.GetPkgCount("port", "installed"))
		}
	case Windows:
		if utils.CommandExists("choco") {
			pdata := os.Getenv("ProgramData")
			if pdata == "" {
				pdata = `C:\ProgramData`
			}
			add("choco", utils.FilePathGlobLens(filepath.Join(pdata, "chocolatey/lib/*")))
		}
		if utils.CommandExists("winget") {
			add("winget", utils.GetPkgCount("winget", "list", "--accept-source-agreements"))
		}
	}

	return formatPackages(managers, options.PackageManagers)
}

func formatPackages(managers []PkgInfo, mode string) string {
	total := 0
	counts, names := []string{}, []string{}
	for _, m := range managers {
		total += m.Count
		counts = append(counts, fmt.Sprintf("%d (%s)", m.Count, m.Manager))
		names = append(names, m.Manager)
	}
	if len(managers) == 0 {
		return ""
	}
	switch mode {
	case "off":
		return fmt.Sprint(total)
	case "tiny":
		return fmt.Sprintf("%d (%s)", total, strings.Join(names, ", "))
	default:
		return strings.Join(counts, ", ")
	}
}
