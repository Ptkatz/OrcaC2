// Copyright (c) 2016 Andreas Auernhammer. All rights reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package threefish

func (t *threefish256) Encrypt(dst, src []byte) {
	var block [4]uint64

	bytesToBlock256(&block, src)

	Encrypt256(&block, &(t.keys), &(t.tweak))

	block256ToBytes(dst, &block)
}

func (t *threefish256) Decrypt(dst, src []byte) {
	var block [4]uint64

	bytesToBlock256(&block, src)

	Decrypt256(&block, &(t.keys), &(t.tweak))

	block256ToBytes(dst, &block)
}

func newCipher256(tweak *[TweakSize]byte, key []byte) *threefish256 {
	c := new(threefish256)

	c.tweak[0] = uint64(tweak[0]) | uint64(tweak[1])<<8 | uint64(tweak[2])<<16 | uint64(tweak[3])<<24 |
		uint64(tweak[4])<<32 | uint64(tweak[5])<<40 | uint64(tweak[6])<<48 | uint64(tweak[7])<<56

	c.tweak[1] = uint64(tweak[8]) | uint64(tweak[9])<<8 | uint64(tweak[10])<<16 | uint64(tweak[11])<<24 |
		uint64(tweak[12])<<32 | uint64(tweak[13])<<40 | uint64(tweak[14])<<48 | uint64(tweak[15])<<56

	c.tweak[2] = c.tweak[0] ^ c.tweak[1]

	for i := range c.keys[:4] {
		j := i * 8
		c.keys[i] = uint64(key[j]) | uint64(key[j+1])<<8 | uint64(key[j+2])<<16 | uint64(key[j+3])<<24 |
			uint64(key[j+4])<<32 | uint64(key[j+5])<<40 | uint64(key[j+6])<<48 | uint64(key[j+7])<<56
	}
	c.keys[4] = C240 ^ c.keys[0] ^ c.keys[1] ^ c.keys[2] ^ c.keys[3]

	return c
}

// Encrypt256 encrypts the 4 words of block using the expanded 256 bit key and
// the 128 bit tweak. The keys[4] must be keys[0] xor keys[1] xor ... keys[3] xor C240.
// The tweak[2] must be tweak[0] xor tweak[1].
func Encrypt256(block *[4]uint64, keys *[5]uint64, tweak *[3]uint64) {
	b0, b1, b2, b3 := block[0], block[1], block[2], block[3]

	for r := 0; r < 17; r++ {
		b0 += keys[r%5]
		b1 += keys[(r+1)%5] + tweak[r%3]
		b2 += keys[(r+2)%5] + tweak[(r+1)%3]
		b3 += keys[(r+3)%5] + uint64(r)

		b0 += b1
		b1 = ((b1 << 14) | (b1 >> (64 - 14))) ^ b0
		b2 += b3
		b3 = ((b3 << 16) | (b3 >> (64 - 16))) ^ b2

		b0 += b3
		b3 = ((b3 << 52) | (b3 >> (64 - 52))) ^ b0
		b2 += b1
		b1 = ((b1 << 57) | (b1 >> (64 - 57))) ^ b2

		b0 += b1
		b1 = ((b1 << 23) | (b1 >> (64 - 23))) ^ b0
		b2 += b3
		b3 = ((b3 << 40) | (b3 >> (64 - 40))) ^ b2

		b0 += b3
		b3 = ((b3 << 5) | (b3 >> (64 - 5))) ^ b0
		b2 += b1
		b1 = ((b1 << 37) | (b1 >> (64 - 37))) ^ b2

		r++

		b0 += keys[r%5]
		b1 += keys[(r+1)%5] + tweak[r%3]
		b2 += keys[(r+2)%5] + tweak[(r+1)%3]
		b3 += keys[(r+3)%5] + uint64(r)

		b0 += b1
		b1 = ((b1 << 25) | (b1 >> (64 - 25))) ^ b0
		b2 += b3
		b3 = ((b3 << 33) | (b3 >> (64 - 33))) ^ b2

		b0 += b3
		b3 = ((b3 << 46) | (b3 >> (64 - 46))) ^ b0
		b2 += b1
		b1 = ((b1 << 12) | (b1 >> (64 - 12))) ^ b2

		b0 += b1
		b1 = ((b1 << 58) | (b1 >> (64 - 58))) ^ b0
		b2 += b3
		b3 = ((b3 << 22) | (b3 >> (64 - 22))) ^ b2

		b0 += b3
		b3 = ((b3 << 32) | (b3 >> (64 - 32))) ^ b0
		b2 += b1
		b1 = ((b1 << 32) | (b1 >> (64 - 32))) ^ b2
	}

	b0 += keys[3]
	b1 += keys[4] + tweak[0]
	b2 += keys[0] + tweak[1]
	b3 += keys[1] + uint64(18)

	block[0], block[1], block[2], block[3] = b0, b1, b2, b3
}

// Decrypt256 decrypts the 4 words of block using the expanded 256 bit key and
// the 128 bit tweak. The keys[4] must be keys[0] xor keys[1] xor ... keys[3] xor C240.
// The tweak[2] must be tweak[0] xor tweak[1].
func Decrypt256(block *[4]uint64, keys *[5]uint64, tweak *[3]uint64) {
	b0, b1, b2, b3 := block[0], block[1], block[2], block[3]

	var tmp uint64
	for r := 18; r > 1; r-- {
		b0 -= keys[r%5]
		b1 -= keys[(r+1)%5] + tweak[r%3]
		b2 -= keys[(r+2)%5] + tweak[(r+1)%3]
		b3 -= keys[(r+3)%5] + uint64(r)

		tmp = b1 ^ b2
		b1 = (tmp >> 32) | (tmp << (64 - 32))
		b2 -= b1
		tmp = b3 ^ b0
		b3 = (tmp >> 32) | (tmp << (64 - 32))
		b0 -= b3

		tmp = b3 ^ b2
		b3 = (tmp >> 22) | (tmp << (64 - 22))
		b2 -= b3
		tmp = b1 ^ b0
		b1 = (tmp >> 58) | (tmp << (64 - 58))
		b0 -= b1

		tmp = b1 ^ b2
		b1 = (tmp >> 12) | (tmp << (64 - 12))
		b2 -= b1
		tmp = b3 ^ b0
		b3 = (tmp >> 46) | (tmp << (64 - 46))
		b0 -= b3

		tmp = b3 ^ b2
		b3 = (tmp >> 33) | (tmp << (64 - 33))
		b2 -= b3
		tmp = b1 ^ b0
		b1 = (tmp >> 25) | (tmp << (64 - 25))
		b0 -= b1

		r--

		b0 -= keys[r%5]
		b1 -= keys[(r+1)%5] + tweak[r%3]
		b2 -= keys[(r+2)%5] + tweak[(r+1)%3]
		b3 -= keys[(r+3)%5] + uint64(r)

		tmp = b1 ^ b2
		b1 = (tmp >> 37) | (tmp << (64 - 37))
		b2 -= b1
		tmp = b3 ^ b0
		b3 = (tmp >> 5) | (tmp << (64 - 5))
		b0 -= b3

		tmp = b3 ^ b2
		b3 = (tmp >> 40) | (tmp << (64 - 40))
		b2 -= b3
		tmp = b1 ^ b0
		b1 = (tmp >> 23) | (tmp << (64 - 23))
		b0 -= b1

		tmp = b1 ^ b2
		b1 = (tmp >> 57) | (tmp << (64 - 57))
		b2 -= b1
		tmp = b3 ^ b0
		b3 = (tmp >> 52) | (tmp << (64 - 52))
		b0 -= b3

		tmp = b3 ^ b2
		b3 = (tmp >> 16) | (tmp << (64 - 16))
		b2 -= b3
		tmp = b1 ^ b0
		b1 = (tmp >> 14) | (tmp << (64 - 14))
		b0 -= b1
	}

	b0 -= keys[0]
	b1 -= keys[1] + tweak[0]
	b2 -= keys[2] + tweak[1]
	b3 -= keys[3]

	block[0], block[1], block[2], block[3] = b0, b1, b2, b3
}

// UBI256 does a Threefish256 encryption of the given block using
// the chain values hVal and the tweak.
// The chain values are updated through hVal[i] = block[i] ^ Enc(block)[i]
func UBI256(block *[4]uint64, hVal *[5]uint64, tweak *[3]uint64) {
	b0, b1, b2, b3 := block[0], block[1], block[2], block[3]

	hVal[4] = C240 ^ hVal[0] ^ hVal[1] ^ hVal[2] ^ hVal[3]
	tweak[2] = tweak[0] ^ tweak[1]

	Encrypt256(block, hVal, tweak)

	hVal[0] = block[0] ^ b0
	hVal[1] = block[1] ^ b1
	hVal[2] = block[2] ^ b2
	hVal[3] = block[3] ^ b3
}
