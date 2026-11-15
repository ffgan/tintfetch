# Tintfetch

一个用于快速展示当前系统信息的命令行工具，对标 neofetch / hyfetch。纯 Go 实现，单个二进制文件，默认在左侧绘制发行版 ASCII logo，右侧列出主机信息。

hyfetch 的旗帜配色、预设、`hyfetch.json` 和命令行参数可以直接用。没有 hyfetch 配置、也没传相关参数时，输出仍是原来的 neofetch 风格。

## 安装

[Releases](https://github.com/ffgan/tintfetch/releases) 里按平台放了压缩包。Windows 是 `.zip`，其他系统是 `.tar.gz`。新包里是可执行文件、`LICENSE`、`COPYING`、`THIRD_PARTY_NOTICES`、`SOURCE.txt` 和这份说明。同一个 release 里的 `tintfetch_<version>_checksums.txt` 是 SHA-256。

这些包用当前还在维护的最老 Go 版本构建，版本号写在 [`.go-version`](.go-version)（现在是 Go 1.26）。二进制沿用这个工具链的系统要求，例如 macOS 12，以及 Linux ppc64 的 ELFv1。

目标来自该版本的 `go tool dist list`。32 位 ARM 会另外打出 ARMv5、ARMv6、ARMv7；MIPS 会同时打出 hardfloat 和 softfloat。文件名里的系统、架构和变体能直接对上：

- `linux_armv6` 给树莓派 Zero 这种 ARMv6，`linux_armv5` 给更老的设备，一般的 32 位 ARM 用 `linux_armv7`。FreeBSD、NetBSD、OpenBSD、Plan 9 也是同样三份。
- 没有硬件浮点的 MIPS 机器用 `*_softfloat`，有 FPU 的用 `*_hardfloat`。
- `android_arm64` 不依赖 NDK，可以放进 Termux。`android_386`、`android_amd64`、`android_armv7` 链接 Android 系统 libc（API 21）。
- `ios_arm64` 面向真机（最低 iOS 12），`ios_amd64` 面向 x86_64 模拟器。
- `js_wasm` 和 `wasip1_wasm` 里的文件名是 `tintfetch.wasm`。

Go 1.26 已经去掉 `windows/arm`，也不再提供 `freebsd/riscv64`，所以发布里没有这两项。

## 构建

需要 Go 1.26 或更新的版本。

```bash
go build -o tintfetch .
```

打出和 release 相同的一套包：

```bash
# android/386、android/amd64、android/arm 需要 ANDROID_NDK_HOME。
# ios/arm64 只能在装了 Xcode 的 macOS 上编。
VERSION=0.3.6 ANDROID_NDK_HOME=/opt/android-ndk bash scripts/build-release.sh
```

## 用法

```bash
tintfetch
tintfetch --stdout
tintfetch --logo Arch
tintfetch --disable disk network
tintfetch --print_config > ~/.config/tintfetch/config.conf

tintfetch -p transgender          # 用跨性别旗给 logo 上色
tintfetch -p '#ff0000,#00ff00,#0000ff' -m rgb
tintfetch --distro Arch           # 指定发行版 logo，并按 hyfetch 配置上色
tintfetch -c                      # 写入 ~/.config/hyfetch.json
tintfetch -b fastfetch            # 把上好色的 logo 交给 fastfetch
tintfetch --print-font-logo
tintfetch --june
```

已安装 hyfetch 时，`~/.config/hyfetch.json`（或 `-C` 指定的文件）会被直接读取：预设、`8bit`/`rgb`、水平/垂直/随机上色、明度、`--args` 里的 neofetch 参数都会生效。预设表与 hyfetch 附带的 `presets.json` 相同。`-b` 的取值如果是 `ascii` 这类图片后端，仍然交给原来的 logo 逻辑；只有 `neofetch`、`fastfetch`、`macchina` 会切到 hyfetch 的后端。

不带参数时，`tintfetch` 会依次查找这些配置：

- `$XDG_CONFIG_HOME/tintfetch/config.conf`
- `~/.config/tintfetch/config.conf`
- `~/.config/neofetch/config.conf`
- `/etc/tintfetch/config.conf`
- `/etc/neofetch/config.conf`

配置格式和 neofetch 相同：`print_info()` 决定显示哪些行，其余是 `key="value"`。用 `--print_config` 可以打印一份当前默认配置。

信息行里暂时检测不到的项目会跳过，不会占一行空输出。`--v` 会把跳过的项目写到标准错误。

## fastfetch 模块

fastfetch 的模块都可以写进 `print_info()`，名字和 fastfetch 一致。已经有的 neofetch 行会接着用原来的检测，例如 `os` 就是发行版，`host` 就是机型，`display` 就是分辨率。主板、机箱、启动项、亮度、Btrfs、摄像头、蓝牙适配器、CPU 缓存和占用、磁盘 IO、物理磁盘、内存条、DNS、Wi-Fi、流量、键鼠手柄、登录管理器、负载、进程、声音、TPM、壁纸、天气、OpenGL、OpenCL、Vulkan、VA-API、ZFS、终端主题，以及 `break`、`separator`、`custom`、`command`、`datetime`、`version` 都会检测。

Linux 读 sysfs 和 `/proc`。Windows 读 SMBIOS、注册表和系统 API，读不到磁盘时再问 CIM。macOS 用 `sysctl`、`ioreg`、`diskutil` 和 `system_profiler`。FreeBSD、OpenBSD、NetBSD、DragonFly 用 `sysctl`、`kenv`、`geom` 和 `iostat`。磁盘 IO 和网卡流量默认是一小段采样后的速率，和 fastfetch 一样。

```bash
tintfetch --list_modules
tintfetch -s os:host:kernel:cpu:cpucache:memory:disk
tintfetch -s all
tintfetch --json
tintfetch --custom "hello from tintfetch"
```

`command="uname -r"` 配 `info "Uname" command` 会跑这条命令。`weather_location` 和 `weather_format` 交给 wttr.in，格式和 fastfetch 一样。`prin "标签" "内容"` 直接印一行，不检测。

`--structure all` 会把每个模块都跑一遍。天气、CPU 占用和最忙进程会各等一小会儿；平时不写进默认配置，避免拖慢原来的输出。

## 默认会显示的内容

用户名和主机名、操作系统、机型、内核、开机时间、软件包、Shell、分辨率、桌面环境和窗口管理器、主题、CPU、内存、磁盘、本地 IP、网卡，以及两行色块。桌面、电池、蓝牙这些只在检测得到时出现。

图片 logo 会在安装了 `chafa`、`catimg` 或 `viu` 时使用，例如 `tintfetch --source ./logo.png`。没有这些工具时回退到发行版 ASCII。

## 开源与发行打包

当前开发版为 `0.3.7-dev`。本次更名保留了 master 中的原生 HyFetch、多平台和
fastfetch 模块实现。命令、配置目录、源码模块和发布归档统一使用 tintfetch。
源码仓库为 `ffgan/tintfetch`。

```sh
make check  # 另需 Python 3，用于离线来源校验
make build
make install PREFIX=/usr DESTDIR="$PWD/build/stage"
./build/tintfetch --plain --no-config
./build/tintfetch --shell-config --config /path/to/trusted-config.conf
```

`--plain`、`--no-config`、`--print-config` 是对应旧参数的别名。
Bash 配置必须显式启用且提供文件路径；普通配置不会作为 shell 脚本执行。

DEB/RPM 配方、手册、贡献规范和 CI 均已加入。详见
[兼容性说明](docs/compatibility.md)、[发行说明](docs/releasing.md)、
[贡献指南](CONTRIBUTING.md) 和 [来源审计](docs/provenance.md)。
401 个 ASCII 素材和两份 HyFetch 数据均已匹配到固定上游版本，
并记录校验和；`make check-provenance` 在离线构建时检查审计快照。
OKHSL 转换保留 Björn Ottosson 的 MIT 署名。正式入库仍需发行版维护者审核。

本项目大量使用 AI 辅助开发，包括历史上的 Cursor Agent，以及近期由 OpenAI Codex
辅助的来源审计和发行打包工作。AI 使用披露、贡献者责任和人工复核要求见
[贡献指南](CONTRIBUTING.md#ai-assistance-and-human-review)。自动化检查通过不代表
发行版批准，也不代表所有代码来源已完成人工核查。

## 许可与对应源码

ARM ID 表明确采用 util-linux 的 GPL-2.0-or-later 数据，并保留 Fastfetch 的 MIT
贡献声明。包含该表的 Tintfetch 二进制按 GPL-2.0-or-later 分发；其余 MIT 源文件
继续保留 MIT 许可，Go 运行时保留 BSD 声明。完整 GPL 文本见 `COPYING`，逐项说明
见 `THIRD_PARTY_NOTICES` 和 [来源说明](docs/provenance.md)。

正式 release 同时提供 `tintfetch_<version>_source.tar.gz`，内含构建二进制所需的
全部项目源码、素材、ARM 原始输入、生成脚本和构建安装脚本。重新分发二进制时，
应按 GPL 提供完整对应源码；仅附仓库链接或许可证文件不足以替代该义务。
