// Package aes implements AES utilities for CryptoNight usage.
//
// Most files are ported from Go's crypto/aes package.
//
// Since CryptoNight's use of AES is quite non-standard and not intended
// for encryption, you must use this package this package with care for
// project that's not CryptoNight associated.
package aes

// CnExpandKey expands exactly 10 round keys.
//
// key must have at least 2 elements.
//
// The result may vary from different architecture, but the output parameter
// rkeys is guranteed to give correct result when used as input in CnRounds.
//
// Note that this is CryptoNight specific.
// This is non-standard AES!
func CnExpandKey(key []uint64, rkeys *[40]uint32) {
	CnExpandKeyGo(key, rkeys)
}

// CnRounds = (SubBytes, ShiftRows, MixColumns, AddRoundKey) * 10,
//
// dst and src must have at least 2 elements.
//
// Note that this is CryptoNight specific.
// This is non-standard AES!
func CnRounds(dst, src []uint64, rkeys *[40]uint32) {
	CnRoundsGo(dst, src, rkeys)
}

// CnSingleRound performs exactly one AES round, i.e.
// one (SubBytes, ShiftRows, MixColumns, AddRoundKey).
//
// dst and src must have at least 2 elements.
//
// Note that this is CryptoNight specific.
// CnSingleRound * 10 might not be equivalent to one CnRounds.
func CnSingleRound(dst, src []uint64, rkey *[2]uint64) {
	CnSingleRoundGo(dst, src, rkey)
}
