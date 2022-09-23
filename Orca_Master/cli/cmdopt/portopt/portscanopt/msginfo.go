package portscanopt

import (
	"github.com/4dogs-cn/TXPortMap/pkg/output"
	"sync"
)

type ScanCmdMsg struct {
	CmdIps     []string
	CmdPorts   []string
	CmdT1000   bool
	CmdRandom  bool
	NumThreads int
	Limit      int
	ExcIps     []string
	ExcPorts   []string
	Tout       float64
	Nbtscan    bool
}

type NBTScanIPMap struct {
	sync.Mutex
	IPS map[string]struct{}
}

var (
	Writer     output.Writer
	NBTScanIPs = NBTScanIPMap{IPS: make(map[string]struct{})}
)
