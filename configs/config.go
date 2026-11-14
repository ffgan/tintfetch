package configs

type Config struct {
	Title
	Kernel
	Distro
	Uptime
	Memory
	Packages
	Shell
	Editor
	CPU
	GPU
	Resolution
	GTK
	IPAddress
	DE
	Disk
	Song
	TextColor
	TextOptions
	ColorBlocks
	ProgressBar
	BackendSettings
	ASCIIOptions
	ImageOptions
	MiscOptions
	PrintInfo
}

type Title struct {
	TitleFqdn string `config:"title_fqdn"`
}

type Kernel struct {
	KernelShorthand string `config:"kernel_shorthand"`
}

type Distro struct {
	DistroShorthand string `config:"distro_shorthand"`
	OSArch          string `config:"os_arch"`
}

type Uptime struct {
	UptimeShorthand string `config:"uptime_shorthand"`
}

type Memory struct {
	MemoryPercent string `config:"memory_percent"`
	MemoryUnit    string `config:"memory_unit"`
	MemPrecision  string `config:"mem_precision"`
}

type Packages struct {
	PackageManagers string `config:"package_managers"`
	PackageSeparate string `config:"package_separate"`
	PackageMinimal  string `config:"package_minimal"`
}

type Shell struct {
	ShellPath    string `config:"shell_path"`
	ShellVersion string `config:"shell_version"`
}

type Editor struct {
	EditorPath    string `config:"editor_path"`
	EditorVersion string `config:"editor_version"`
}

type CPU struct {
	SpeedType      string `config:"speed_type"`
	SpeedShorthand string `config:"speed_shorthand"`
	CPUBrand       string `config:"cpu_brand"`
	CPUSpeed       string `config:"cpu_speed"`
	CPUCores       string `config:"cpu_cores"`
	CPUTemp        string `config:"cpu_temp"`
}

type GPU struct {
	GPUBrand string `config:"gpu_brand"`
	GPUType  string `config:"gpu_type"`
}

type Resolution struct {
	RefreshRate string `config:"refresh_rate"`
}

type GTK struct {
	GTKShorthand string `config:"gtk_shorthand"`
	GTK2         string `config:"gtk2"`
	GTK3         string `config:"gtk3"`
	QT           string `config:"qt"`
}

type IPAddress struct {
	PublicIPHost     string `config:"public_ip_host"`
	PublicIPTimeout  string `config:"public_ip_timeout"`
	LocalIPInterface string `config:"local_ip_interface"`
}

type DE struct {
	DeVersion string `config:"de_version"`
}

type Disk struct {
	DiskShow     string `config:"disk_show"`
	DiskSubtitle string `config:"disk_subtitle"`
	DiskPercent  string `config:"disk_percent"`
}

type Song struct {
	MusicPlayer   string `config:"music_player"`
	SongFormat    string `config:"song_format"`
	SongShorthand string `config:"song_shorthand"`
	MPCArgs       string `config:"mpc_args"`
}

type TextColor struct {
	Colors string `config:"colors"`
}

type TextOptions struct {
	Bold             string `config:"bold"`
	UnderlineEnabled string `config:"underline_enabled"`
	UnderlineChar    string `config:"underline_char"`
	Separator        string `config:"separator"`
}

type ColorBlocks struct {
	BlockRange  string `config:"block_range"`
	ColorBlocks string `config:"color_blocks"`
	BlockWidth  string `config:"block_width"`
	BlockHeight string `config:"block_height"`
	ColOffset   string `config:"col_offset"`
}

type ProgressBar struct {
	BarCharElapsed  string `config:"bar_char_elapsed"`
	BarCharTotal    string `config:"bar_char_total"`
	BarBorder       string `config:"bar_border"`
	BarLength       string `config:"bar_length"`
	BarColorElapsed string `config:"bar_color_elapsed"`
	BarColorTotal   string `config:"bar_color_total"`
	MemoryDisplay   string `config:"memory_display"`
	BatteryDisplay  string `config:"battery_display"`
	DiskDisplay     string `config:"disk_display"`
}

type BackendSettings struct {
	ImageBackend string `config:"image_backend"`
	ImageSource  string `config:"image_source"`
}

type ASCIIOptions struct {
	ASCIIDistro string `config:"ascii_distro"`
	ASCIIColors string `config:"ascii_colors"`
	ASCIIBold   string `config:"ascii_bold"`
}

type ImageOptions struct {
	ImageLoop       string `config:"image_loop"`
	ThumbnailDir    string `config:"thumbnail_dir"`
	CropMode        string `config:"crop_mode"`
	CropOffset      string `config:"crop_offset"`
	ImageSize       string `config:"image_size"`
	CatimgSize      string `config:"catimg_size"`
	Gap             string `config:"gap"`
	Yoffset         string `config:"yoffset"`
	Xoffset         string `config:"xoffset"`
	BackgroundColor string `config:"background_color"`
}

type MiscOptions struct {
	Stdout string `config:"stdout"`
}

type PrintInfo struct {
	InfoList [][]string `config:"info_list"`
}
