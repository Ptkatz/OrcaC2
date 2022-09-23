package aes

import (
	"math/bits"
	"unsafe"
)

func CnExpandKeyGoSoft(key []uint64, rkeys *[40]uint32) {
	for i := 0; i < 4; i++ {
		rkeys[2*i] = bits.ReverseBytes32(uint32(key[i]))
		rkeys[2*i+1] = bits.ReverseBytes32(uint32(key[i] >> 32))
	}

	for i := 8; i < 40; i++ {
		t := rkeys[i-1]
		if i%8 == 0 {
			t = subw(rotw(t)) ^ (uint32(powx[i/8-1]) << 24)
		} else if 8 > 6 && i%8 == 4 {
			t = subw(t)
		}
		rkeys[i] = rkeys[i-8] ^ t
	}
}

func CnRoundsGoSoft(dst, src []uint64, rkeys *[40]uint32) {
	src8 := (*[16]byte)(unsafe.Pointer(&src[0]))
	dst8 := (*[16]byte)(unsafe.Pointer(&dst[0]))

	var s0, s1, s2, s3, t0, t1, t2, t3 uint32

	s0 = uint32(src8[0])<<24 | uint32(src8[1])<<16 | uint32(src8[2])<<8 | uint32(src8[3])
	s1 = uint32(src8[4])<<24 | uint32(src8[5])<<16 | uint32(src8[6])<<8 | uint32(src8[7])
	s2 = uint32(src8[8])<<24 | uint32(src8[9])<<16 | uint32(src8[10])<<8 | uint32(src8[11])
	s3 = uint32(src8[12])<<24 | uint32(src8[13])<<16 | uint32(src8[14])<<8 | uint32(src8[15])

	for r := 0; r < 10; r++ {
		t0 = rkeys[4*r+0] ^ te0[uint8(s0>>24)] ^ te1[uint8(s1>>16)] ^ te2[uint8(s2>>8)] ^ te3[uint8(s3)]
		t1 = rkeys[4*r+1] ^ te0[uint8(s1>>24)] ^ te1[uint8(s2>>16)] ^ te2[uint8(s3>>8)] ^ te3[uint8(s0)]
		t2 = rkeys[4*r+2] ^ te0[uint8(s2>>24)] ^ te1[uint8(s3>>16)] ^ te2[uint8(s0>>8)] ^ te3[uint8(s1)]
		t3 = rkeys[4*r+3] ^ te0[uint8(s3>>24)] ^ te1[uint8(s0>>16)] ^ te2[uint8(s1>>8)] ^ te3[uint8(s2)]
		s0, s1, s2, s3 = t0, t1, t2, t3
	}

	dst8[0], dst8[1], dst8[2], dst8[3] = byte(s0>>24), byte(s0>>16), byte(s0>>8), byte(s0)
	dst8[4], dst8[5], dst8[6], dst8[7] = byte(s1>>24), byte(s1>>16), byte(s1>>8), byte(s1)
	dst8[8], dst8[9], dst8[10], dst8[11] = byte(s2>>24), byte(s2>>16), byte(s2>>8), byte(s2)
	dst8[12], dst8[13], dst8[14], dst8[15] = byte(s3>>24), byte(s3>>16), byte(s3>>8), byte(s3)
}

func CnSingleRoundGoSoft(dst, src []uint64, rkey *[2]uint64) {
	src8 := (*[16]byte)(unsafe.Pointer(&src[0]))
	dst8 := (*[16]byte)(unsafe.Pointer(&dst[0]))
	rkey32 := (*[4]uint32)(unsafe.Pointer(&rkey[0]))

	var t0, t1, t2, t3 uint32

	t0 = rkey32[0] ^ ter0[src8[0]] ^ ter1[src8[5]] ^ ter2[src8[10]] ^ ter3[src8[15]]
	t1 = rkey32[1] ^ ter0[src8[4]] ^ ter1[src8[9]] ^ ter2[src8[14]] ^ ter3[src8[3]]
	t2 = rkey32[2] ^ ter0[src8[8]] ^ ter1[src8[13]] ^ ter2[src8[2]] ^ ter3[src8[7]]
	t3 = rkey32[3] ^ ter0[src8[12]] ^ ter1[src8[1]] ^ ter2[src8[6]] ^ ter3[src8[11]]

	dst8[0], dst8[1], dst8[2], dst8[3] = byte(t0), byte(t0>>8), byte(t0>>16), byte(t0>>24)
	dst8[4], dst8[5], dst8[6], dst8[7] = byte(t1), byte(t1>>8), byte(t1>>16), byte(t1>>24)
	dst8[8], dst8[9], dst8[10], dst8[11] = byte(t2), byte(t2>>8), byte(t2>>16), byte(t2>>24)
	dst8[12], dst8[13], dst8[14], dst8[15] = byte(t3), byte(t3>>8), byte(t3>>16), byte(t3>>24)
}

func CnSingleRoundHeavyGo(dst, src []uint64, rkey *[2]uint64) {
	dst[0] = src[0]
	dst[1] = src[1]

	var x [2]uint64
	var k [2]uint64
	k[0] = rkey[0]
	k[1] = rkey[1]

	x[0] = dst[0] ^ 0xffffffffffffffff
	x[1] = dst[1] ^ 0xffffffffffffffff

	kk := (*[4]uint32)(unsafe.Pointer(&k[0]))
	xx := (*[4]uint32)(unsafe.Pointer(&x[0]))
	xxx := (*[16]byte)(unsafe.Pointer(&x[0]))

	kk[0] ^= ter0[xxx[0*4+0]] ^ ter1[xxx[1*4+1]] ^ ter2[xxx[2*4+2]] ^ ter3[xxx[3*4+3]]
	xx[0] ^= kk[0]
	kk[1] ^= ter0[xxx[1*4+0]] ^ ter1[xxx[2*4+1]] ^ ter2[xxx[3*4+2]] ^ ter3[xxx[0*4+3]]
	xx[1] ^= kk[1]
	kk[2] ^= ter0[xxx[2*4+0]] ^ ter1[xxx[3*4+1]] ^ ter2[xxx[0*4+2]] ^ ter3[xxx[1*4+3]]
	xx[2] ^= kk[2]
	kk[3] ^= ter0[xxx[3*4+0]] ^ ter1[xxx[0*4+1]] ^ ter2[xxx[1*4+2]] ^ ter3[xxx[2*4+3]]

	dst[0] = k[0]
	dst[1] = k[1]
}

// Apply sbox0 to each byte in w.
func subw(w uint32) uint32 {
	return uint32(sbox0[w>>24])<<24 |
		uint32(sbox0[w>>16&0xff])<<16 |
		uint32(sbox0[w>>8&0xff])<<8 |
		uint32(sbox0[w&0xff])
}

// Rotate
func rotw(w uint32) uint32 { return w<<8 | w>>24 }
