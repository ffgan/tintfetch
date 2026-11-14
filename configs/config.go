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
	Fastfetch
	PrintInfo
}

// Fastfetch holds options for the modules that came from fastfetch.
type Fastfetch struct {
	CustomText      string `tintfetch:"custom"`
	CommandText     string `tintfetch:"command"`
	WeatherLocation string `tintfetch:"weather_location"`
	WeatherFormat   string `tintfetch:"weather_format"`
	SeparatorString string `tintfetch:"separator_string"`
}

type Title struct {
	TitleFqdn string `tintfetch:"title_fqdn"`
}

type Kernel struct {
	KernelShorthand string `tintfetch:"kernel_shorthand"`
}

type Distro struct {
	DistroShorthand string `tintfetch:"distro_shorthand"`
	OSArch          string `tintfetch:"os_arch"`
}

type Uptime struct {
	UptimeShorthand string `tintfetch:"uptime_shorthand"`
}

type Memory struct {
	MemoryPercent string `tintfetch:"memory_percent"`
	MemoryUnit    string `tintfetch:"memory_unit"`
	MemPrecision  string `tintfetch:"memory_precision"`
}

type Packages struct {
	PackageManagers string `tintfetch:"package_managers"`
	PackageSeparate string `tintfetch:"package_separate"`
	PackageMinimal  string `tintfetch:"package_minimal"`
}

type Shell struct {
	ShellPath    string `tintfetch:"shell_path"`
	ShellVersion string `tintfetch:"shell_version"`
}

type Editor struct {
	EditorPath    string `tintfetch:"editor_path"`
	EditorVersion string `tintfetch:"editor_version"`
}

type CPU struct {
	SpeedType      string `tintfetch:"speed_type"`
	SpeedShorthand string `tintfetch:"speed_shorthand"`
	CPUBrand       string `tintfetch:"cpu_brand"`
	CPUSpeed       string `tintfetch:"cpu_speed"`
	CPUCores       string `tintfetch:"cpu_cores"`
	CPUTemp        string `tintfetch:"cpu_temp"`
}

type GPU struct {
	GPUBrand string `tintfetch:"gpu_brand"`
	GPUType  string `tintfetch:"gpu_type"`
}

type Resolution struct {
	RefreshRate string `tintfetch:"refresh_rate"`
}

type GTK struct {
	GTKShorthand string `tintfetch:"gtk_shorthand"`
	GTK2         string `tintfetch:"gtk2"`
	GTK3         string `tintfetch:"gtk3"`
	QT           string `tintfetch:"qt"`
}

type IPAddress struct {
	PublicIPHost     string `tintfetch:"public_ip_host"`
	PublicIPTimeout  string `tintfetch:"public_ip_timeout"`
	LocalIPInterface string `tintfetch:"local_ip_interface"`
}

type DE struct {
	DeVersion string `tintfetch:"de_version"`
}

type Disk struct {
	DiskShow     string `tintfetch:"disk_show"`
	DiskSubtitle string `tintfetch:"disk_subtitle"`
	DiskPercent  string `tintfetch:"disk_percent"`
}

type Song struct {
	MusicPlayer   string `tintfetch:"music_player"`
	SongFormat    string `tintfetch:"song_format"`
	SongShorthand string `tintfetch:"song_shorthand"`
	MPCArgs       string `tintfetch:"mpc_args"`
}

type TextColor struct {
	Colors string `tintfetch:"colors"`
}

type TextOptions struct {
	Bold             string `tintfetch:"bold"`
	UnderlineEnabled string `tintfetch:"underline_enabled"`
	UnderlineChar    string `tintfetch:"underline_char"`
	Separator        string `tintfetch:"separator"`
}

type ColorBlocks struct {
	BlockRange  string `tintfetch:"block_range"`
	Enabled     string `tintfetch:"color_blocks"`
	BlockWidth  string `tintfetch:"block_width"`
	BlockHeight string `tintfetch:"block_height"`
	ColOffset   string `tintfetch:"col_offset"`
}

type ProgressBar struct {
	BarCharElapsed  string `tintfetch:"bar_char_elapsed"`
	BarCharTotal    string `tintfetch:"bar_char_total"`
	BarBorder       string `tintfetch:"bar_border"`
	BarLength       string `tintfetch:"bar_length"`
	BarColorElapsed string `tintfetch:"bar_color_elapsed"`
	BarColorTotal   string `tintfetch:"bar_color_total"`
	MemoryDisplay   string `tintfetch:"memory_display"`
	BatteryDisplay  string `tintfetch:"battery_display"`
	DiskDisplay     string `tintfetch:"disk_display"`
}

type BackendSettings struct {
	ImageBackend string `tintfetch:"image_backend"`
	ImageSource  string `tintfetch:"image_source"`
}

type ASCIIOptions struct {
	ASCIIDistro string `tintfetch:"ascii_distro"`
	ASCIIColors string `tintfetch:"ascii_colors"`
	ASCIIBold   string `tintfetch:"ascii_bold"`
}

type ImageOptions struct {
	ImageLoop       string `tintfetch:"image_loop"`
	ThumbnailDir    string `tintfetch:"thumbnail_dir"`
	CropMode        string `tintfetch:"crop_mode"`
	CropOffset      string `tintfetch:"crop_offset"`
	ImageSize       string `tintfetch:"image_size"`
	CatimgSize      string `tintfetch:"catimg_size"`
	Gap             string `tintfetch:"gap"`
	Yoffset         string `tintfetch:"yoffset"`
	Xoffset         string `tintfetch:"xoffset"`
	BackgroundColor string `tintfetch:"background_color"`
}

type MiscOptions struct {
	Stdout string `tintfetch:"stdout"`
}

type PrintInfo struct {
	InfoList [][]string `tintfetch:"info_list"`
}
