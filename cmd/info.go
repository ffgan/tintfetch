package cmd

import (
	"runtime"
	"strings"

	"github.com/ffgan/tintfetch/configs"
	cli "github.com/ffgan/tintfetch/internal/CLI"
	gui "github.com/ffgan/tintfetch/internal/GUI"
	dev "github.com/ffgan/tintfetch/internal/devices"
	"github.com/ffgan/tintfetch/internal/utils"
)

type hardware struct {
	KernelName     string
	kernel_version string
	KernelMachine  string
	DarwinName     string
	os             string
	distro         string
	ascii_distro   string

	bios string

	model string

	user     string
	hostname string
	title    string
	length   int

	kernel string

	uptime string

	packages string

	shell string

	editor string

	de string

	wm string

	wm_theme string

	cpu string
	gpu string

	memory string

	network string

	bluetooth string

	resolution string

	term      string
	term_font string

	underline string
	theme     string
	icons     string
	cursor    string
	cols      string
	configs.Config
}

func (h *hardware) Get_os() {
	h.os = cli.GetOS(h.KernelName, h.DarwinName)
}

func (h *hardware) Get_distro() {
	h.distro, h.ascii_distro = cli.GetDistro(h.os, h.OSArch, h.KernelMachine, h.DistroShorthand, h.ASCIIDistro)
}

func (h *hardware) Get_bios() {
	h.bios = dev.GetBIOS(h.os)
}

func (h *hardware) Get_model() {
	h.model = cli.GetModel(h.os, h.KernelName, h.KernelMachine)
}
func (h *hardware) Get_title() {
	title, err := cli.GetTitleFQDN(h.TitleFqdn == utils.ON)
	if err != nil {
		return
	}
	h.user = title.User
	h.hostname = title.Hostname
	h.title = title.Title
	h.length = title.Length
}

func (h *hardware) Get_kernel() {
	h.kernel = cli.GetKernel(h.OSArch, h.DistroShorthand, h.KernelShorthand, h.ascii_distro, h.os, h.kernel_version, h.KernelName, h.KernelMachine)
}

func (h *hardware) Get_uptime() {
	h.uptime = cli.GetUptime(h.os, h.UptimeShorthand)
}

func (h *hardware) Get_packages() {
	h.packages = cli.DetectPackagesWithOptions(h.os, h.Packages)
}

func (h *hardware) Get_shell() {
	h.shell = cli.GetShell(h.ShellPath, h.ShellVersion)
}

func (h *hardware) Get_editor() {
	h.editor = cli.GetEditor(h.EditorPath, h.EditorVersion)
}

func (h *hardware) Get_de() {
	h.de = gui.GetDE(h.os, h.distro, h.KernelName)
}

func (h *hardware) Get_wm() {
	h.wm = gui.GetWM(h.os, h.KernelName)
}

func (h *hardware) Get_wm_theme() {
	h.wm_theme = gui.GetTheme()
}

func (h *hardware) Get_cpu() {
	unit := h.CPUTemp
	if unit == utils.ON {
		unit = "C"
	}
	h.cpu = cli.GetCPU(h.os, h.KernelMachine, h.CPUCores, h.CPUSpeed, h.CPUTemp, h.CPUBrand, h.CPUCores, h.SpeedShorthand, unit, h.SpeedType)
}

func (h *hardware) Get_gpu() {
	h.gpu = dev.GetGPU(h.os)
}

func (h *hardware) Get_memory() {
	h.memory = cli.PrintMem(h.os, h.Memory, h.ProgressBar)
}

func (h *hardware) Get_network() {
	h.network = dev.GetNetwork(h.os)
}

func (h *hardware) Get_bluetooth() {
	h.bluetooth = dev.GetBluetooth()
}

func (h *hardware) Get_resolution() {
	h.resolution = dev.GetResolution(h.os, h.KernelMachine)
}

func (h *hardware) Get_theme() {
	h.theme = gui.GetTheme()
}

func (h *hardware) Get_icons() {
	h.icons = gui.GetIcons()
}

func (h *hardware) Get_font() {
	h.term_font = gui.GetFont()
}

func (h *hardware) Get_cursor() {
	h.cursor = gui.GetCursor()
}

func (h *hardware) Get_term() {
	h.term = cli.GetTerm(h.os)
}

func (h *hardware) Get_term_font() {
	h.term_font = gui.GetFont()
}

func cache_uname() (string, string, string, string) {
	fields := strings.Fields(cli.UName("-srm"))
	if len(fields) >= 3 {
		return fields[0], fields[1], fields[2], ""
	}
	name := map[string]string{"linux": "Linux", "darwin": "Darwin", "windows": "Windows_NT", "freebsd": "FreeBSD", "openbsd": "OpenBSD", "netbsd": "NetBSD", "dragonfly": "DragonFly", "solaris": "SunOS"}[runtime.GOOS]
	if name == "" {
		name = runtime.GOOS
	}
	return name, "", runtime.GOARCH, ""
}
