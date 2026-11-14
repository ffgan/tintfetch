package cmd

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/ffgan/tintfetch/configs"
	"github.com/ffgan/tintfetch/internal/utils"
)

func diskInfo(c configs.Config) string {
	mounts, err := configs.Words(strings.Trim(c.DiskShow, "()"))
	if err != nil {
		return ""
	}
	var result []string
	for _, mount := range mounts {
		out, err := utils.CommandOutput("df", "-Pk", "--", mount)
		if err != nil {
			continue
		}
		lines := strings.Split(strings.TrimSpace(string(out)), "\n")
		if len(lines) < 2 {
			continue
		}
		fields := strings.Fields(lines[len(lines)-1])
		if len(fields) < 6 {
			continue
		}
		total, e1 := strconv.ParseFloat(fields[1], 64)
		used, e2 := strconv.ParseFloat(fields[2], 64)
		if e1 != nil || e2 != nil {
			continue
		}
		label := mount
		if c.DiskSubtitle == "name" {
			label = fields[0]
		} else if c.DiskSubtitle == "dir" {
			label = filepath.Base(mount)
		}
		value := fmt.Sprintf("%s: %.1f GiB / %.1f GiB", label, used/1048576, total/1048576)
		if c.DiskPercent == "on" {
			value += " (" + fields[4] + ")"
		}
		result = append(result, value)
	}
	return strings.Join(result, ", ")
}
func batteryInfo() string {
	paths, _ := filepath.Glob("/sys/class/power_supply/*")
	var result []string
	for _, path := range paths {
		if utils.ReadFileTrim(filepath.Join(path, "type")) != "Battery" {
			continue
		}
		capacity := utils.ReadFileTrim(filepath.Join(path, "capacity"))
		if capacity != "" {
			result = append(result, filepath.Base(path)+": "+capacity+"% ("+utils.ReadFileTrim(filepath.Join(path, "status"))+")")
		}
	}
	return strings.Join(result, ", ")
}
func adapterInfo() string {
	paths, _ := filepath.Glob("/sys/class/power_supply/*")
	for _, path := range paths {
		if utils.ReadFileTrim(filepath.Join(path, "type")) == "Mains" {
			if utils.ReadFileTrim(filepath.Join(path, "online")) == "1" {
				return "Connected"
			}
			return "Disconnected"
		}
	}
	return ""
}
func localIP(selection string) string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return ""
	}
	names, _ := configs.Words(strings.Trim(selection, "()"))
	auto := len(names) == 0 || (len(names) == 1 && names[0] == "auto")
	var result []string
	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		selected := auto
		for _, name := range names {
			if iface.Name == name {
				selected = true
			}
		}
		if !selected {
			continue
		}
		addresses, _ := iface.Addrs()
		for _, addr := range addresses {
			ip, _, err := net.ParseCIDR(addr.String())
			if err == nil && ip.IsGlobalUnicast() {
				result = append(result, iface.Name+": "+ip.String())
			}
		}
	}
	sort.Strings(result)
	return strings.Join(result, ", ")
}
func publicIP(host, timeout string) string {
	target, err := url.Parse(host)
	if err != nil || (target.Scheme != "https" && target.Scheme != "http") || target.Host == "" {
		return ""
	}
	client := http.Client{Timeout: time.Duration(atoiDefault(timeout, 2)) * time.Second}
	response, err := client.Get(target.String())
	if err != nil {
		return ""
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return ""
	}
	data, err := io.ReadAll(io.LimitReader(response.Body, 256))
	if err != nil {
		return ""
	}
	address := strings.TrimSpace(string(data))
	if net.ParseIP(address) == nil {
		return ""
	}
	return address
}
func loggedUsers() string {
	out, err := utils.CommandOutput("who")
	if err != nil {
		return ""
	}
	users := map[string]bool{}
	for _, line := range strings.Split(string(out), "\n") {
		fields := strings.Fields(line)
		if len(fields) > 0 {
			users[fields[0]] = true
		}
	}
	names := []string{}
	for name := range users {
		names = append(names, name)
	}
	sort.Strings(names)
	return strings.Join(names, ", ")
}
func toolVersion(name string, args ...string) string {
	out, err := utils.CommandCombinedOutput(name, args...)
	if err != nil {
		return ""
	}
	return strings.Split(strings.TrimSpace(string(out)), "\n")[0]
}
func gpuDriver() string {
	paths, _ := filepath.Glob("/sys/class/drm/card[0-9]*/device/driver")
	seen := map[string]bool{}
	for _, path := range paths {
		if target, err := os.Readlink(path); err == nil {
			seen[filepath.Base(target)] = true
		}
	}
	drivers := []string{}
	for name := range seen {
		drivers = append(drivers, name)
	}
	sort.Strings(drivers)
	return strings.Join(drivers, ", ")
}
