// Copyright (c) 2016 Andreas Auernhammer. All rights reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package threefish

func (t *threefish1024) Encrypt(dst, src []byte) {
	var block [16]uint64

	bytesToBlock1024(&block, src)

	Encrypt1024(&block, &(t.keys), &(t.tweak))

	block1024ToBytes(dst, &block)
}

func (t *threefish1024) Decrypt(dst, src []byte) {
	var block [16]uint64

	bytesToBlock1024(&block, src)

	Decrypt1024(&block, &(t.keys), &(t.tweak))

	block1024ToBytes(dst, &block)
}

func newCipher1024(tweak *[TweakSize]byte, key []byte) *threefish1024 {
	c := new(threefish1024)

	c.tweak[0] = uint64(tweak[0]) | uint64(tweak[1])<<8 | uint64(tweak[2])<<16 | uint64(tweak[3])<<24 |
		uint64(tweak[4])<<32 | uint64(tweak[5])<<40 | uint64(tweak[6])<<48 | uint64(tweak[7])<<56

	c.tweak[1] = uint64(tweak[8]) | uint64(tweak[9])<<8 | uint64(tweak[10])<<16 | uint64(tweak[11])<<24 |
		uint64(tweak[12])<<32 | uint64(tweak[13])<<40 | uint64(tweak[14])<<48 | uint64(tweak[15])<<56

	c.tweak[2] = c.tweak[0] ^ c.tweak[1]

	for i := range c.keys[:16] {
		j := i * 8
		c.keys[i] = uint64(key[j]) | uint64(key[j+1])<<8 | uint64(key[j+2])<<16 | uint64(key[j+3])<<24 |
			uint64(key[j+4])<<32 | uint64(key[j+5])<<40 | uint64(key[j+6])<<48 | uint64(key[j+7])<<56
	}
	c.keys[16] = C240 ^ c.keys[0] ^ c.keys[1] ^ c.keys[2] ^ c.keys[3] ^ c.keys[4] ^ c.keys[5] ^ c.keys[6] ^
		c.keys[7] ^ c.keys[8] ^ c.keys[9] ^ c.keys[10] ^ c.keys[11] ^ c.keys[12] ^ c.keys[13] ^ c.keys[14] ^ c.keys[15]

	return c
}

// Encrypt1024 encrypts the 16 words of block using the expanded 1024 bit key and
// the 128 bit tweak. The keys[16] must be keys[0] xor keys[1] xor ... keys[15] xor C240.
// The tweak[2] must be tweak[0] xor tweak[1].
func Encrypt1024(block *[16]uint64, keys *[17]uint64, tweak *[3]uint64) {
	b0, b1, b2, b3 := block[0], block[1], block[2], block[3]
	b4, b5, b6, b7 := block[4], block[5], block[6], block[7]
	b8, b9, b10, b11 := block[8], block[9], block[10], block[11]
	b12, b13, b14, b15 := block[12], block[13], block[14], block[15]

	for r := 0; r < 19; r++ {
		b0 += keys[r%17]
		b1 += keys[(r+1)%17]
		b2 += keys[(r+2)%17]
		b3 += keys[(r+3)%17]
		b4 += keys[(r+4)%17]
		b5 += keys[(r+5)%17]
		b6 += keys[(r+6)%17]
		b7 += keys[(r+7)%17]
		b8 += keys[(r+8)%17]
		b9 += keys[(r+9)%17]
		b10 += keys[(r+10)%17]
		b11 += keys[(r+11)%17]
		b12 += keys[(r+12)%17]
		b13 += keys[(r+13)%17] + tweak[r%3]
		b14 += keys[(r+14)%17] + tweak[(r+1)%3]
		b15 += keys[(r+15)%17] + uint64(r)

		b0 += b1
		b1 = ((b1 << 24) | (b1 >> (64 - 24))) ^ b0
		b2 += b3
		b3 = ((b3 << 13) | (b3 >> (64 - 13))) ^ b2
		b4 += b5
		b5 = ((b5 << 8) | (b5 >> (64 - 8))) ^ b4
		b6 += b7
		b7 = ((b7 << 47) | (b7 >> (64 - 47))) ^ b6
		b8 += b9
		b9 = ((b9 << 8) | (b9 >> (64 - 8))) ^ b8
		b10 += b11
		b11 = ((b11 << 17) | (b11 >> (64 - 17))) ^ b10
		b12 += b13
		b13 = ((b13 << 22) | (b13 >> (64 - 22))) ^ b12
		b14 += b15
		b15 = ((b15 << 37) | (b15 >> (64 - 37))) ^ b14

		b0 += b9
		b9 = ((b9 << 38) | (b9 >> (64 - 38))) ^ b0
		b2 += b13
		b13 = ((b13 << 19) | (b13 >> (64 - 19))) ^ b2
		b6 += b11
		b11 = ((b11 << 10) | (b11 >> (64 - 10))) ^ b6
		b4 += b15
		b15 = ((b15 << 55) | (b15 >> (64 - 55))) ^ b4
		b10 += b7
		b7 = ((b7 << 49) | (b7 >> (64 - 49))) ^ b10
		b12 += b3
		b3 = ((b3 << 18) | (b3 >> (64 - 18))) ^ b12
		b14 += b5
		b5 = ((b5 << 23) | (b5 >> (64 - 23))) ^ b14
		b8 += b1
		b1 = ((b1 << 52) | (b1 >> (64 - 52))) ^ b8

		b0 += b7
		b7 = ((b7 << 33) | (b7 >> (64 - 33))) ^ b0
		b2 += b5
		b5 = ((b5 << 4) | (b5 >> (64 - 4))) ^ b2
		b4 += b3
		b3 = ((b3 << 51) | (b3 >> (64 - 51))) ^ b4
		b6 += b1
		b1 = ((b1 << 13) | (b1 >> (64 - 13))) ^ b6
		b12 += b15
		b15 = ((b15 << 34) | (b15 >> (64 - 34))) ^ b12
		b14 += b13
		b13 = ((b13 << 41) | (b13 >> (64 - 41))) ^ b14
		b8 += b11
		b11 = ((b11 << 59) | (b11 >> (64 - 59))) ^ b8
		b10 += b9
		b9 = ((b9 << 17) | (b9 >> (64 - 17))) ^ b10

		b0 += b15
		b15 = ((b15 << 5) | (b15 >> (64 - 5))) ^ b0
		b2 += b11
		b11 = ((b11 << 20) | (b11 >> (64 - 20))) ^ b2
		b6 += b13
		b13 = ((b13 << 48) | (b13 >> (64 - 48))) ^ b6
		b4 += b9
		b9 = ((b9 << 41) | (b9 >> (64 - 41))) ^ b4
		b14 += b1
		b1 = ((b1 << 47) | (b1 >> (64 - 47))) ^ b14
		b8 += b5
		b5 = ((b5 << 28) | (b5 >> (64 - 28))) ^ b8
		b10 += b3
		b3 = ((b3 << 16) | (b3 >> (64 - 16))) ^ b10
		b12 += b7
		b7 = ((b7 << 25) | (b7 >> (64 - 25))) ^ b12

		r++

		b0 += keys[r%17]
		b1 += keys[(r+1)%17]
		b2 += keys[(r+2)%17]
		b3 += keys[(r+3)%17]
		b4 += keys[(r+4)%17]
		b5 += keys[(r+5)%17]
		b6 += keys[(r+6)%17]
		b7 += keys[(r+7)%17]
		b8 += keys[(r+8)%17]
		b9 += keys[(r+9)%17]
		b10 += keys[(r+10)%17]
		b11 += keys[(r+11)%17]
		b12 += keys[(r+12)%17]
		b13 += keys[(r+13)%17] + tweak[r%3]
		b14 += keys[(r+14)%17] + tweak[(r+1)%3]
		b15 += keys[(r+15)%17] + uint64(r)

		b0 += b1
		b1 = ((b1 << 41) | (b1 >> (64 - 41))) ^ b0
		b2 += b3
		b3 = ((b3 << 9) | (b3 >> (64 - 9))) ^ b2
		b4 += b5
		b5 = ((b5 << 37) | (b5 >> (64 - 37))) ^ b4
		b6 += b7
		b7 = ((b7 << 31) | (b7 >> (64 - 31))) ^ b6
		b8 += b9
		b9 = ((b9 << 12) | (b9 >> (64 - 12))) ^ b8
		b10 += b11
		b11 = ((b11 << 47) | (b11 >> (64 - 47))) ^ b10
		b12 += b13
		b13 = ((b13 << 44) | (b13 >> (64 - 44))) ^ b12
		b14 += b15
		b15 = ((b15 << 30) | (b15 >> (64 - 30))) ^ b14

		b0 += b9
		b9 = ((b9 << 16) | (b9 >> (64 - 16))) ^ b0
		b2 += b13
		b13 = ((b13 << 34) | (b13 >> (64 - 34))) ^ b2
		b6 += b11
		b11 = ((b11 << 56) | (b11 >> (64 - 56))) ^ b6
		b4 += b15
		b15 = ((b15 << 51) | (b15 >> (64 - 51))) ^ b4
		b10 += b7
		b7 = ((b7 << 4) | (b7 >> (64 - 4))) ^ b10
		b12 += b3
		b3 = ((b3 << 53) | (b3 >> (64 - 53))) ^ b12
		b14 += b5
		b5 = ((b5 << 42) | (b5 >> (64 - 42))) ^ b14
		b8 += b1
		b1 = ((b1 << 41) | (b1 >> (64 - 41))) ^ b8

		b0 += b7
		b7 = ((b7 << 31) | (b7 >> (64 - 31))) ^ b0
		b2 += b5
		b5 = ((b5 << 44) | (b5 >> (64 - 44))) ^ b2
		b4 += b3
		b3 = ((b3 << 47) | (b3 >> (64 - 47))) ^ b4
		b6 += b1
		b1 = ((b1 << 46) | (b1 >> (64 - 46))) ^ b6
		b12 += b15
		b15 = ((b15 << 19) | (b15 >> (64 - 19))) ^ b12
		b14 += b13
		b13 = ((b13 << 42) | (b13 >> (64 - 42))) ^ b14
		b8 += b11
		b11 = ((b11 << 44) | (b11 >> (64 - 44))) ^ b8
		b10 += b9
		b9 = ((b9 << 25) | (b9 >> (64 - 25))) ^ b10

		b0 += b15
		b15 = ((b15 << 9) | (b15 >> (64 - 9))) ^ b0
		b2 += b11
		b11 = ((b11 << 48) | (b11 >> (64 - 48))) ^ b2
		b6 += b13
		b13 = ((b13 << 35) | (b13 >> (64 - 35))) ^ b6
		b4 += b9
		b9 = ((b9 << 52) | (b9 >> (64 - 52))) ^ b4
		b14 += b1
		b1 = ((b1 << 23) | (b1 >> (64 - 23))) ^ b14
		b8 += b5
		b5 = ((b5 << 31) | (b5 >> (64 - 31))) ^ b8
		b10 += b3
		b3 = ((b3 << 37) | (b3 >> (64 - 37))) ^ b10
		b12 += b7
		b7 = ((b7 << 20) | (b7 >> (64 - 20))) ^ b12
	}

	b0 += keys[3]
	b1 += keys[4]
	b2 += keys[5]
	b3 += keys[6]
	b4 += keys[7]
	b5 += keys[8]
	b6 += keys[9]
	b7 += keys[10]
	b8 += keys[11]
	b9 += keys[12]
	b10 += keys[13]
	b11 += keys[14]
	b12 += keys[15]
	b13 += keys[16] + tweak[2]
	b14 += keys[0] + tweak[0]
	b15 += keys[1] + 20

	block[0], block[1], block[2], block[3] = b0, b1, b2, b3
	block[4], block[5], block[6], block[7] = b4, b5, b6, b7
	block[8], block[9], block[10], block[11] = b8, b9, b10, b11
	block[12], block[13], block[14], block[15] = b12, b13, b14, b15
}

// Decrypt1024 decrypts the 16 words of block using the expanded 1024 bit key and
// the 128 bit tweak. The keys[16] must be keys[0] xor keys[1] xor ... keys[15] xor C240.
// The tweak[2] must be tweak[0] xor tweak[1].
func Decrypt1024(block *[16]uint64, keys *[17]uint64, tweak *[3]uint64) {
	b0, b1, b2, b3 := block[0], block[1], block[2], block[3]
	b4, b5, b6, b7 := block[4], block[5], block[6], block[7]
	b8, b9, b10, b11 := block[8], block[9], block[10], block[11]
	b12, b13, b14, b15 := block[12], block[13], block[14], block[15]

	var tmp uint64
	for r := 20; r > 1; r-- {
		b0 -= keys[r%17]
		b1 -= keys[(r+1)%17]
		b2 -= keys[(r+2)%17]
		b3 -= keys[(r+3)%17]
		b4 -= keys[(r+4)%17]
		b5 -= keys[(r+5)%17]
		b6 -= keys[(r+6)%17]
		b7 -= keys[(r+7)%17]
		b8 -= keys[(r+8)%17]
		b9 -= keys[(r+9)%17]
		b10 -= keys[(r+10)%17]
		b11 -= keys[(r+11)%17]
		b12 -= keys[(r+12)%17]
		b13 -= keys[(r+13)%17] + tweak[r%3]
		b14 -= keys[(r+14)%17] + tweak[(r+1)%3]
		b15 -= keys[(r+15)%17] + uint64(r)

		tmp = b7 ^ b12
		b7 = (tmp >> 20) | (tmp << (64 - 20))
		b12 -= b7
		tmp = b3 ^ b10
		b3 = (tmp >> 37) | (tmp << (64 - 37))
		b10 -= b3
		tmp = b5 ^ b8
		b5 = (tmp >> 31) | (tmp << (64 - 31))
		b8 -= b5
		tmp = b1 ^ b14
		b1 = (tmp >> 23) | (tmp << (64 - 23))
		b14 -= b1
		tmp = b9 ^ b4
		b9 = (tmp >> 52) | (tmp << (64 - 52))
		b4 -= b9
		tmp = b13 ^ b6
		b13 = (tmp >> 35) | (tmp << (64 - 35))
		b6 -= b13
		tmp = b11 ^ b2
		b11 = (tmp >> 48) | (tmp << (64 - 48))
		b2 -= b11
		tmp = b15 ^ b0
		b15 = (tmp >> 9) | (tmp << (64 - 9))
		b0 -= b15

		tmp = b9 ^ b10
		b9 = (tmp >> 25) | (tmp << (64 - 25))
		b10 -= b9
		tmp = b11 ^ b8
		b11 = (tmp >> 44) | (tmp << (64 - 44))
		b8 -= b11
		tmp = b13 ^ b14
		b13 = (tmp >> 42) | (tmp << (64 - 42))
		b14 -= b13
		tmp = b15 ^ b12
		b15 = (tmp >> 19) | (tmp << (64 - 19))
		b12 -= b15
		tmp = b1 ^ b6
		b1 = (tmp >> 46) | (tmp << (64 - 46))
		b6 -= b1
		tmp = b3 ^ b4
		b3 = (tmp >> 47) | (tmp << (64 - 47))
		b4 -= b3
		tmp = b5 ^ b2
		b5 = (tmp >> 44) | (tmp << (64 - 44))
		b2 -= b5
		tmp = b7 ^ b0
		b7 = (tmp >> 31) | (tmp << (64 - 31))
		b0 -= b7

		tmp = b1 ^ b8
		b1 = (tmp >> 41) | (tmp << (64 - 41))
		b8 -= b1
		tmp = b5 ^ b14
		b5 = (tmp >> 42) | (tmp << (64 - 42))
		b14 -= b5
		tmp = b3 ^ b12
		b3 = (tmp >> 53) | (tmp << (64 - 53))
		b12 -= b3
		tmp = b7 ^ b10
		b7 = (tmp >> 4) | (tmp << (64 - 4))
		b10 -= b7
		tmp = b15 ^ b4
		b15 = (tmp >> 51) | (tmp << (64 - 51))
		b4 -= b15
		tmp = b11 ^ b6
		b11 = (tmp >> 56) | (tmp << (64 - 56))
		b6 -= b11
		tmp = b13 ^ b2
		b13 = (tmp >> 34) | (tmp << (64 - 34))
		b2 -= b13
		tmp = b9 ^ b0
		b9 = (tmp >> 16) | (tmp << (64 - 16))
		b0 -= b9

		tmp = b15 ^ b14
		b15 = (tmp >> 30) | (tmp << (64 - 30))
		b14 -= b15
		tmp = b13 ^ b12
		b13 = (tmp >> 44) | (tmp << (64 - 44))
		b12 -= b13
		tmp = b11 ^ b10
		b11 = (tmp >> 47) | (tmp << (64 - 47))
		b10 -= b11
		tmp = b9 ^ b8
		b9 = (tmp >> 12) | (tmp << (64 - 12))
		b8 -= b9
		tmp = b7 ^ b6
		b7 = (tmp >> 31) | (tmp << (64 - 31))
		b6 -= b7
		tmp = b5 ^ b4
		b5 = (tmp >> 37) | (tmp << (64 - 37))
		b4 -= b5
		tmp = b3 ^ b2
		b3 = (tmp >> 9) | (tmp << (64 - 9))
		b2 -= b3
		tmp = b1 ^ b0
		b1 = (tmp >> 41) | (tmp << (64 - 41))
		b0 -= b1

		r--

		b0 -= keys[r%17]
		b1 -= keys[(r+1)%17]
		b2 -= keys[(r+2)%17]
		b3 -= keys[(r+3)%17]
		b4 -= keys[(r+4)%17]
		b5 -= keys[(r+5)%17]
		b6 -= keys[(r+6)%17]
		b7 -= keys[(r+7)%17]
		b8 -= keys[(r+8)%17]
		b9 -= keys[(r+9)%17]
		b10 -= keys[(r+10)%17]
		b11 -= keys[(r+11)%17]
		b12 -= keys[(r+12)%17]
		b13 -= keys[(r+13)%17] + tweak[r%3]
		b14 -= keys[(r+14)%17] + tweak[(r+1)%3]
		b15 -= keys[(r+15)%17] + uint64(r)

		tmp = b7 ^ b12
		b7 = (tmp >> 25) | (tmp << (64 - 25))
		b12 -= b7
		tmp = b3 ^ b10
		b3 = (tmp >> 16) | (tmp << (64 - 16))
		b10 -= b3
		tmp = b5 ^ b8
		b5 = (tmp >> 28) | (tmp << (64 - 28))
		b8 -= b5
		tmp = b1 ^ b14
		b1 = (tmp >> 47) | (tmp << (64 - 47))
		b14 -= b1
		tmp = b9 ^ b4
		b9 = (tmp >> 41) | (tmp << (64 - 41))
		b4 -= b9
		tmp = b13 ^ b6
		b13 = (tmp >> 48) | (tmp << (64 - 48))
		b6 -= b13
		tmp = b11 ^ b2
		b11 = (tmp >> 20) | (tmp << (64 - 20))
		b2 -= b11
		tmp = b15 ^ b0
		b15 = (tmp >> 5) | (tmp << (64 - 5))
		b0 -= b15

		tmp = b9 ^ b10
		b9 = (tmp >> 17) | (tmp << (64 - 17))
		b10 -= b9
		tmp = b11 ^ b8
		b11 = (tmp >> 59) | (tmp << (64 - 59))
		b8 -= b11
		tmp = b13 ^ b14
		b13 = (tmp >> 41) | (tmp << (64 - 41))
		b14 -= b13
		tmp = b15 ^ b12
		b15 = (tmp >> 34) | (tmp << (64 - 34))
		b12 -= b15
		tmp = b1 ^ b6
		b1 = (tmp >> 13) | (tmp << (64 - 13))
		b6 -= b1
		tmp = b3 ^ b4
		b3 = (tmp >> 51) | (tmp << (64 - 51))
		b4 -= b3
		tmp = b5 ^ b2
		b5 = (tmp >> 4) | (tmp << (64 - 4))
		b2 -= b5
		tmp = b7 ^ b0
		b7 = (tmp >> 33) | (tmp << (64 - 33))
		b0 -= b7

		tmp = b1 ^ b8
		b1 = (tmp >> 52) | (tmp << (64 - 52))
		b8 -= b1
		tmp = b5 ^ b14
		b5 = (tmp >> 23) | (tmp << (64 - 23))
		b14 -= b5
		tmp = b3 ^ b12
		b3 = (tmp >> 18) | (tmp << (64 - 18))
		b12 -= b3
		tmp = b7 ^ b10
		b7 = (tmp >> 49) | (tmp << (64 - 49))
		b10 -= b7
		tmp = b15 ^ b4
		b15 = (tmp >> 55) | (tmp << (64 - 55))
		b4 -= b15
		tmp = b11 ^ b6
		b11 = (tmp >> 10) | (tmp << (64 - 10))
		b6 -= b11
		tmp = b13 ^ b2
		b13 = (tmp >> 19) | (tmp << (64 - 19))
		b2 -= b13
		tmp = b9 ^ b0
		b9 = (tmp >> 38) | (tmp << (64 - 38))
		b0 -= b9

		tmp = b15 ^ b14
		b15 = (tmp >> 37) | (tmp << (64 - 37))
		b14 -= b15
		tmp = b13 ^ b12
		b13 = (tmp >> 22) | (tmp << (64 - 22))
		b12 -= b13
		tmp = b11 ^ b10
		b11 = (tmp >> 17) | (tmp << (64 - 17))
		b10 -= b11
		tmp = b9 ^ b8
		b9 = (tmp >> 8) | (tmp << (64 - 8))
		b8 -= b9
		tmp = b7 ^ b6
		b7 = (tmp >> 47) | (tmp << (64 - 47))
		b6 -= b7
		tmp = b5 ^ b4
		b5 = (tmp >> 8) | (tmp << (64 - 8))
		b4 -= b5
		tmp = b3 ^ b2
		b3 = (tmp >> 13) | (tmp << (64 - 13))
		b2 -= b3
		tmp = b1 ^ b0
		b1 = (tmp >> 24) | (tmp << (64 - 24))
		b0 -= b1
	}

	b0 -= keys[0]
	b1 -= keys[1]
	b2 -= keys[2]
	b3 -= keys[3]
	b4 -= keys[4]
	b5 -= keys[5]
	b6 -= keys[6]
	b7 -= keys[7]
	b8 -= keys[8]
	b9 -= keys[9]
	b10 -= keys[10]
	b11 -= keys[11]
	b12 -= keys[12]
	b13 -= keys[13] + tweak[0]
	b14 -= keys[14] + tweak[1]
	b15 -= keys[15]

	block[0], block[1], block[2], block[3] = b0, b1, b2, b3
	block[4], block[5], block[6], block[7] = b4, b5, b6, b7
	block[8], block[9], block[10], block[11] = b8, b9, b10, b11
	block[12], block[13], block[14], block[15] = b12, b13, b14, b15
}

// UBI1024 does a Threefish1024 encryption of the given block using
// the chain values hVal and the tweak.
// The chain values are updated through hVal[i] = block[i] ^ Enc(block)[i]
func UBI1024(block *[16]uint64, hVal *[17]uint64, tweak *[3]uint64) {
	b0, b1, b2, b3 := block[0], block[1], block[2], block[3]
	b4, b5, b6, b7 := block[4], block[5], block[6], block[7]
	b8, b9, b10, b11 := block[8], block[9], block[10], block[11]
	b12, b13, b14, b15 := block[12], block[13], block[14], block[15]

	hVal[16] = C240 ^ hVal[0] ^ hVal[1] ^ hVal[2] ^ hVal[3] ^ hVal[4] ^ hVal[5] ^ hVal[6] ^ hVal[7] ^
		hVal[8] ^ hVal[9] ^ hVal[10] ^ hVal[11] ^ hVal[12] ^ hVal[13] ^ hVal[14] ^ hVal[15]
	tweak[2] = tweak[0] ^ tweak[1]

	Encrypt1024(block, hVal, tweak)

	hVal[0] = block[0] ^ b0
	hVal[1] = block[1] ^ b1
	hVal[2] = block[2] ^ b2
	hVal[3] = block[3] ^ b3
	hVal[4] = block[4] ^ b4
	hVal[5] = block[5] ^ b5
	hVal[6] = block[6] ^ b6
	hVal[7] = block[7] ^ b7
	hVal[8] = block[8] ^ b8
	hVal[9] = block[9] ^ b9
	hVal[10] = block[10] ^ b10
	hVal[11] = block[11] ^ b11
	hVal[12] = block[12] ^ b12
	hVal[13] = block[13] ^ b13
	hVal[14] = block[14] ^ b14
	hVal[15] = block[15] ^ b15
}
