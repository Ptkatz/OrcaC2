package portscanopt

import (
	"github.com/4dogs-cn/TXPortMap/pkg/output"
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

var ResultEvents []*output.ResultEvent

type ScanResultMsg struct {
	Info   Info
	Target string
}

type Info struct {
	Banner  string
	Service string
	Cert    string
}
