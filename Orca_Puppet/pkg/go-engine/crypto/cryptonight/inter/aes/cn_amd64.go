// +build amd64

package aes

import (
	"Orca_Puppet/pkg/go-engine/loggo"
	"golang.org/x/sys/cpu"
	"unsafe"
)

var (
	hasAES = cpu.X86.HasAES
)

func CnExpandKeyGo(key []uint64, rkeys *[40]uint32) {
	if !hasAES {
		CnExpandKeyGoSoft(key, rkeys)
		return
	}
	if uintptr(unsafe.Pointer(&key[0]))%16 != 0 {
		loggo.Error("CnExpandKeyGo %v", &key[0])
	}
	if uintptr(unsafe.Pointer(&rkeys[0]))%16 != 0 {
		loggo.Error("CnExpandKeyGo %v", &rkeys[0])
	}
	CnExpandKeyAsm(&key[0], rkeys)
}

func CnRoundsGo(dst, src []uint64, rkeys *[40]uint32) {
	if !hasAES {
		CnRoundsGoSoft(dst, src, rkeys)
		return
	}
	if uintptr(unsafe.Pointer(&dst[0]))%16 != 0 {
		loggo.Error("CnRoundsGo %v", &dst[0])
	}
	if uintptr(unsafe.Pointer(&src[0]))%16 != 0 {
		loggo.Error("CnRoundsGo %v", &src[0])
	}
	if uintptr(unsafe.Pointer(&rkeys[0]))%16 != 0 {
		loggo.Error("CnRoundsGo %v", &rkeys[0])
	}
	CnRoundsAsm(&dst[0], &src[0], &rkeys[0])
}

func CnSingleRoundGo(dst, src []uint64, rkey *[2]uint64) {
	if !hasAES {
		CnSingleRoundGoSoft(dst, src, rkey)
		return
	}
	if uintptr(unsafe.Pointer(&dst[0]))%16 != 0 {
		loggo.Error("CnSingleRoundGo %v", &dst[0])
	}
	if uintptr(unsafe.Pointer(&src[0]))%16 != 0 {
		loggo.Error("CnSingleRoundGo %v", &src[0])
	}
	if uintptr(unsafe.Pointer(&rkey[0]))%16 != 0 {
		loggo.Error("CnSingleRoundGo %v", &rkey[0])
	}
	CnSingleRoundAsm(&dst[0], &src[0], &rkey[0])
}

//go:noescape
func CnExpandKeyAsm(src *uint64, rkey *[40]uint32)

//go:noescape
func CnRoundsAsm(dst, src *uint64, rkeys *uint32)

//go:noescape
func CnSingleRoundAsm(dst, src *uint64, rkeys *uint64)
