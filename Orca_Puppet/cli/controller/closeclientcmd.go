package controller

import (
	"Orca_Puppet/define/debug"
	"os"
)

func closeClientCmd() {
	debug.DebugPrint("close client")
	os.Exit(10)
}
