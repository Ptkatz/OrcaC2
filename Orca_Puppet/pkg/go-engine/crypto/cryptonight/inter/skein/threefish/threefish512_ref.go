// Copyright (c) 2016 Andreas Auernhammer. All rights reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package threefish

func (t *threefish512) Encrypt(dst, src []byte) {
	var block [8]uint64

	bytesToBlock512(&block, src)

	Encrypt512(&block, &(t.keys), &(t.tweak))

	block512ToBytes(dst, &block)
}

func (t *threefish512) Decrypt(dst, src []byte) {
	var block [8]uint64

	bytesToBlock512(&block, src)

	Decrypt512(&block, &(t.keys), &(t.tweak))

	block512ToBytes(dst, &block)
}

func newCipher512(tweak *[TweakSize]byte, key []byte) *threefish512 {
	c := new(threefish512)

	c.tweak[0] = uint64(tweak[0]) | uint64(tweak[1])<<8 | uint64(tweak[2])<<16 | uint64(tweak[3])<<24 |
		uint64(tweak[4])<<32 | uint64(tweak[5])<<40 | uint64(tweak[6])<<48 | uint64(tweak[7])<<56

	c.tweak[1] = uint64(tweak[8]) | uint64(tweak[9])<<8 | uint64(tweak[10])<<16 | uint64(tweak[11])<<24 |
		uint64(tweak[12])<<32 | uint64(tweak[13])<<40 | uint64(tweak[14])<<48 | uint64(tweak[15])<<56

	c.tweak[2] = c.tweak[0] ^ c.tweak[1]

	for i := range c.keys[:8] {
		j := i * 8
		c.keys[i] = uint64(key[j]) | uint64(key[j+1])<<8 | uint64(key[j+2])<<16 | uint64(key[j+3])<<24 |
			uint64(key[j+4])<<32 | uint64(key[j+5])<<40 | uint64(key[j+6])<<48 | uint64(key[j+7])<<56
	}
	c.keys[8] = C240 ^ c.keys[0] ^ c.keys[1] ^ c.keys[2] ^ c.keys[3] ^ c.keys[4] ^ c.keys[5] ^ c.keys[6] ^ c.keys[7]

	return c
}

// Encrypt512 encrypts the 8 words of block using the expanded 512 bit key and
// the 128 bit tweak. The keys[8] must be keys[0] xor keys[1] xor ... keys[8] xor C240.
// The tweak[2] must be tweak[0] xor tweak[1].
func Encrypt512(block *[8]uint64, keys *[9]uint64, tweak *[3]uint64) {
	b0, b1, b2, b3 := block[0], block[1], block[2], block[3]
	b4, b5, b6, b7 := block[4], block[5], block[6], block[7]

	for r := 0; r < 17; r++ {
		b0 += keys[r%9]
		b1 += keys[(r+1)%9]
		b2 += keys[(r+2)%9]
		b3 += keys[(r+3)%9]
		b4 += keys[(r+4)%9]
		b5 += keys[(r+5)%9] + tweak[r%3]
		b6 += keys[(r+6)%9] + tweak[(r+1)%3]
		b7 += keys[(r+7)%9] + uint64(r)

		b0 += b1
		b1 = (b1<<46 | b1>>(64-46)) ^ b0
		b2 += b3
		b3 = (b3<<36 | b3>>(64-36)) ^ b2
		b4 += b5
		b5 = (b5<<19 | b5>>(64-19)) ^ b4
		b6 += b7
		b7 = (b7<<37 | b7>>(64-37)) ^ b6

		b2 += b1
		b1 = (b1<<33 | b1>>(64-33)) ^ b2
		b4 += b7
		b7 = (b7<<27 | b7>>(64-27)) ^ b4
		b6 += b5
		b5 = (b5<<14 | b5>>(64-14)) ^ b6
		b0 += b3
		b3 = (b3<<42 | b3>>(64-42)) ^ b0

		b4 += b1
		b1 = (b1<<17 | b1>>(64-17)) ^ b4
		b6 += b3
		b3 = (b3<<49 | b3>>(64-49)) ^ b6
		b0 += b5
		b5 = (b5<<36 | b5>>(64-36)) ^ b0
		b2 += b7
		b7 = (b7<<39 | b7>>(64-39)) ^ b2

		b6 += b1
		b1 = (b1<<44 | b1>>(64-44)) ^ b6
		b0 += b7
		b7 = (b7<<9 | b7>>(64-9)) ^ b0
		b2 += b5
		b5 = (b5<<54 | b5>>(64-54)) ^ b2
		b4 += b3
		b3 = (b3<<56 | b3>>(64-56)) ^ b4

		r++

		b0 += keys[r%9]
		b1 += keys[(r+1)%9]
		b2 += keys[(r+2)%9]
		b3 += keys[(r+3)%9]
		b4 += keys[(r+4)%9]
		b5 += keys[(r+5)%9] + tweak[r%3]
		b6 += keys[(r+6)%9] + tweak[(r+1)%3]
		b7 += keys[(r+7)%9] + uint64(r)

		b0 += b1
		b1 = (b1<<39 | b1>>(64-39)) ^ b0
		b2 += b3
		b3 = (b3<<30 | b3>>(64-30)) ^ b2
		b4 += b5
		b5 = (b5<<34 | b5>>(64-34)) ^ b4
		b6 += b7
		b7 = (b7<<24 | b7>>(64-24)) ^ b6

		b2 += b1
		b1 = (b1<<13 | b1>>(64-13)) ^ b2
		b4 += b7
		b7 = (b7<<50 | b7>>(64-50)) ^ b4
		b6 += b5
		b5 = (b5<<10 | b5>>(64-10)) ^ b6
		b0 += b3
		b3 = (b3<<17 | b3>>(64-17)) ^ b0

		b4 += b1
		b1 = (b1<<25 | b1>>(64-25)) ^ b4
		b6 += b3
		b3 = (b3<<29 | b3>>(64-29)) ^ b6
		b0 += b5
		b5 = (b5<<39 | b5>>(64-39)) ^ b0
		b2 += b7
		b7 = (b7<<43 | b7>>(64-43)) ^ b2

		b6 += b1
		b1 = (b1<<8 | b1>>(64-8)) ^ b6
		b0 += b7
		b7 = (b7<<35 | b7>>(64-35)) ^ b0
		b2 += b5
		b5 = (b5<<56 | b5>>(64-56)) ^ b2
		b4 += b3
		b3 = (b3<<22 | b3>>(64-22)) ^ b4
	}

	b0 += keys[0]
	b1 += keys[1]
	b2 += keys[2]
	b3 += keys[3]
	b4 += keys[4]
	b5 += keys[5] + tweak[0]
	b6 += keys[6] + tweak[1]
	b7 += keys[7] + 18

	block[0], block[1], block[2], block[3] = b0, b1, b2, b3
	block[4], block[5], block[6], block[7] = b4, b5, b6, b7
}

// Decrypt512 decrypts the 8 words of block using the expanded 512 bit key and
// the 128 bit tweak. The keys[8] must be keys[0] xor keys[1] xor ... keys[8] xor C240.
// The tweak[2] must be tweak[0] xor tweak[1].
func Decrypt512(block *[8]uint64, keys *[9]uint64, tweak *[3]uint64) {
	b0, b1, b2, b3 := block[0], block[1], block[2], block[3]
	b4, b5, b6, b7 := block[4], block[5], block[6], block[7]

	var tmp uint64
	for r := 18; r > 1; r-- {
		b0 -= keys[r%9]
		b1 -= keys[(r+1)%9]
		b2 -= keys[(r+2)%9]
		b3 -= keys[(r+3)%9]
		b4 -= keys[(r+4)%9]
		b5 -= keys[(r+5)%9] + tweak[r%3]
		b6 -= keys[(r+6)%9] + tweak[(r+1)%3]
		b7 -= keys[(r+7)%9] + uint64(r)

		tmp = b3 ^ b4
		b3 = tmp>>22 | tmp<<(64-22)
		b4 -= b3
		tmp = b5 ^ b2
		b5 = tmp>>56 | tmp<<(64-56)
		b2 -= b5
		tmp = b7 ^ b0
		b7 = tmp>>35 | tmp<<(64-35)
		b0 -= b7
		tmp = b1 ^ b6
		b1 = tmp>>8 | tmp<<(64-8)
		b6 -= b1

		tmp = b7 ^ b2
		b7 = tmp>>43 | tmp<<(64-43)
		b2 -= b7
		tmp = b5 ^ b0
		b5 = tmp>>39 | tmp<<(64-39)
		b0 -= b5
		tmp = b3 ^ b6
		b3 = tmp>>29 | tmp<<(64-29)
		b6 -= b3
		tmp = b1 ^ b4
		b1 = tmp>>25 | tmp<<(64-25)
		b4 -= b1

		tmp = b3 ^ b0
		b3 = tmp>>17 | tmp<<(64-17)
		b0 -= b3
		tmp = b5 ^ b6
		b5 = tmp>>10 | tmp<<(64-10)
		b6 -= b5
		tmp = b7 ^ b4
		b7 = tmp>>50 | tmp<<(64-50)
		b4 -= b7
		tmp = b1 ^ b2
		b1 = tmp>>13 | tmp<<(64-13)
		b2 -= b1

		tmp = b7 ^ b6
		b7 = tmp>>24 | tmp<<(64-24)
		b6 -= b7
		tmp = b5 ^ b4
		b5 = tmp>>34 | tmp<<(64-34)
		b4 -= b5
		tmp = b3 ^ b2
		b3 = tmp>>30 | tmp<<(64-30)
		b2 -= b3
		tmp = b1 ^ b0
		b1 = tmp>>39 | tmp<<(64-39)
		b0 -= b1

		r--

		b0 -= keys[r%9]
		b1 -= keys[(r+1)%9]
		b2 -= keys[(r+2)%9]
		b3 -= keys[(r+3)%9]
		b4 -= keys[(r+4)%9]
		b5 -= keys[(r+5)%9] + tweak[r%3]
		b6 -= keys[(r+6)%9] + tweak[(r+1)%3]
		b7 -= keys[(r+7)%9] + uint64(r)

		tmp = b3 ^ b4
		b3 = tmp>>56 | tmp<<(64-56)
		b4 -= b3
		tmp = b5 ^ b2
		b5 = tmp>>54 | tmp<<(64-54)
		b2 -= b5
		tmp = b7 ^ b0
		b7 = tmp>>9 | tmp<<(64-9)
		b0 -= b7
		tmp = b1 ^ b6
		b1 = tmp>>44 | tmp<<(64-44)
		b6 -= b1

		tmp = b7 ^ b2
		b7 = tmp>>39 | tmp<<(64-39)
		b2 -= b7
		tmp = b5 ^ b0
		b5 = tmp>>36 | tmp<<(64-36)
		b0 -= b5
		tmp = b3 ^ b6
		b3 = tmp>>49 | tmp<<(64-49)
		b6 -= b3
		tmp = b1 ^ b4
		b1 = tmp>>17 | tmp<<(64-17)
		b4 -= b1

		tmp = b3 ^ b0
		b3 = tmp>>42 | tmp<<(64-42)
		b0 -= b3
		tmp = b5 ^ b6
		b5 = tmp>>14 | tmp<<(64-14)
		b6 -= b5
		tmp = b7 ^ b4
		b7 = tmp>>27 | tmp<<(64-27)
		b4 -= b7
		tmp = b1 ^ b2
		b1 = tmp>>33 | tmp<<(64-33)
		b2 -= b1

		tmp = b7 ^ b6
		b7 = tmp>>37 | tmp<<(64-37)
		b6 -= b7
		tmp = b5 ^ b4
		b5 = tmp>>19 | tmp<<(64-19)
		b4 -= b5
		tmp = b3 ^ b2
		b3 = tmp>>36 | tmp<<(64-36)
		b2 -= b3
		tmp = b1 ^ b0
		b1 = tmp>>46 | tmp<<(64-46)
		b0 -= b1
	}

	b0 -= keys[0]
	b1 -= keys[1]
	b2 -= keys[2]
	b3 -= keys[3]
	b4 -= keys[4]
	b5 -= keys[5] + tweak[0]
	b6 -= keys[6] + tweak[1]
	b7 -= keys[7]

	block[0], block[1], block[2], block[3] = b0, b1, b2, b3
	block[4], block[5], block[6], block[7] = b4, b5, b6, b7
}

// UBI512 does a Threefish512 encryption of the given block using
// the chain values hVal and the tweak.
// The chain values are updated through hVal[i] = block[i] ^ Enc(block)[i]
func UBI512(block *[8]uint64, hVal *[9]uint64, tweak *[3]uint64) {
	b0, b1, b2, b3 := block[0], block[1], block[2], block[3]
	b4, b5, b6, b7 := block[4], block[5], block[6], block[7]

	hVal[8] = C240 ^ hVal[0] ^ hVal[1] ^ hVal[2] ^ hVal[3] ^ hVal[4] ^ hVal[5] ^ hVal[6] ^ hVal[7]
	tweak[2] = tweak[0] ^ tweak[1]

	Encrypt512(block, hVal, tweak)

	hVal[0] = block[0] ^ b0
	hVal[1] = block[1] ^ b1
	hVal[2] = block[2] ^ b2
	hVal[3] = block[3] ^ b3
	hVal[4] = block[4] ^ b4
	hVal[5] = block[5] ^ b5
	hVal[6] = block[6] ^ b6
	hVal[7] = block[7] ^ b7
}
