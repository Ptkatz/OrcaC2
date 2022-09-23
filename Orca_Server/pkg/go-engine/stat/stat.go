package stat

import (
	"errors"
	"io/ioutil"
	"math"
	"os/exec"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type SysInfo struct {
	CPU           float64
	Memory        float64
	VirtualMemory float64
}

type Stat struct {
	utime  float64
	stime  float64
	cutime float64
	cstime float64
	start  float64
	rss    float64
	virt   float64
	uptime float64
}

func parseFloat(val string) float64 {
	floatVal, _ := strconv.ParseFloat(val, 64)
	return floatVal
}

func formatStdOut(stdout []byte, userfulIndex int) []string {
	infoArr := strings.Split(string(stdout), "\n")[userfulIndex]
	ret := strings.Fields(infoArr)
	return ret
}

func GetStat(pid int) (*SysInfo, error) {

	if runtime.GOOS != "linux" {
		return nil, errors.New("GOOS not support")
	}

	_history := Stat{}
	_, newhis, err := getStat(_history, pid)
	if err != nil {
		return nil, err
	}
	time.Sleep(time.Second)
	info, _, err := getStat(newhis, pid)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func getStat(_history Stat, pid int) (SysInfo, Stat, error) {
	sysInfo := SysInfo{}

	var clkTck float64 = 100
	var pageSize float64 = 4096

	uptimeFileBytes, err := ioutil.ReadFile(path.Join("/proc", "uptime"))
	uptime := parseFloat(strings.Split(string(uptimeFileBytes), " ")[0])

	clkTckStdout, err := exec.Command("getconf", "CLK_TCK").Output()
	if err == nil {
		clkTck = parseFloat(formatStdOut(clkTckStdout, 0)[0])
	}

	pageSizeStdout, err := exec.Command("getconf", "PAGESIZE").Output()
	if err == nil {
		pageSize = parseFloat(formatStdOut(pageSizeStdout, 0)[0])
	}

	procStatFileBytes, err := ioutil.ReadFile(path.Join("/proc", strconv.Itoa(pid), "stat"))
	splitAfter := strings.SplitAfter(string(procStatFileBytes), ")")

	if len(splitAfter) == 0 || len(splitAfter) == 1 {
		return sysInfo, _history, errors.New("Can't find process with this PID: " + strconv.Itoa(pid))
	}
	infos := strings.Split(splitAfter[1], " ")
	stat := &Stat{
		utime:  parseFloat(infos[12]),
		stime:  parseFloat(infos[13]),
		cutime: parseFloat(infos[14]),
		cstime: parseFloat(infos[15]),
		start:  parseFloat(infos[20]) / clkTck,
		uptime: uptime,
	}

	procStatFileBytes, err = ioutil.ReadFile(path.Join("/proc", strconv.Itoa(pid), "statm"))
	splitAfter = strings.Split(string(procStatFileBytes), " ")
	stat.virt = parseFloat(splitAfter[0])
	stat.rss = parseFloat(splitAfter[1])

	_stime := 0.0
	_utime := 0.0
	if _history.stime != 0 {
		_stime = _history.stime
	}

	if _history.utime != 0 {
		_utime = _history.utime
	}
	total := stat.stime - _stime + stat.utime - _utime
	total = total / clkTck

	seconds := stat.start - uptime
	if _history.uptime != 0 {
		seconds = uptime - _history.uptime
	}

	seconds = math.Abs(seconds)
	if seconds == 0 {
		seconds = 1
	}

	_history = *stat
	sysInfo.CPU = (total / seconds) * 100
	sysInfo.Memory = stat.rss * pageSize
	sysInfo.VirtualMemory = stat.virt * pageSize

	return sysInfo, _history, nil
}
