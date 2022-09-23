// Copyright (c) 2016 Andreas Auernhammer. All rights reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

// +build !amd64

package threefish

func bytesToBlock256(block *[4]uint64, src []byte) {
	for i := range block {
		j := i * 8
		block[i] = uint64(src[j]) | uint64(src[j+1])<<8 | uint64(src[j+2])<<16 | uint64(src[j+3])<<24 |
			uint64(src[j+4])<<32 | uint64(src[j+5])<<40 | uint64(src[j+6])<<48 | uint64(src[j+7])<<56
	}
}

func block256ToBytes(dst []byte, block *[4]uint64) {
	for i, v := range block {
		j := i * 8
		dst[j] = byte(v)
		dst[j+1] = byte(v >> 8)
		dst[j+2] = byte(v >> 16)
		dst[j+3] = byte(v >> 24)
		dst[j+4] = byte(v >> 32)
		dst[j+5] = byte(v >> 40)
		dst[j+6] = byte(v >> 48)
		dst[j+7] = byte(v >> 56)
	}
}

func bytesToBlock512(block *[8]uint64, src []byte) {
	for i := range block {
		j := i * 8
		block[i] = uint64(src[j]) | uint64(src[j+1])<<8 | uint64(src[j+2])<<16 | uint64(src[j+3])<<24 |
			uint64(src[j+4])<<32 | uint64(src[j+5])<<40 | uint64(src[j+6])<<48 | uint64(src[j+7])<<56
	}
}

func block512ToBytes(dst []byte, block *[8]uint64) {
	for i, v := range block {
		j := i * 8
		dst[j] = byte(v)
		dst[j+1] = byte(v >> 8)
		dst[j+2] = byte(v >> 16)
		dst[j+3] = byte(v >> 24)
		dst[j+4] = byte(v >> 32)
		dst[j+5] = byte(v >> 40)
		dst[j+6] = byte(v >> 48)
		dst[j+7] = byte(v >> 56)
	}
}

func bytesToBlock1024(block *[16]uint64, src []byte) {
	for i := range block {
		j := i * 8
		block[i] = uint64(src[j]) | uint64(src[j+1])<<8 | uint64(src[j+2])<<16 | uint64(src[j+3])<<24 |
			uint64(src[j+4])<<32 | uint64(src[j+5])<<40 | uint64(src[j+6])<<48 | uint64(src[j+7])<<56
	}
}

func block1024ToBytes(dst []byte, block *[16]uint64) {
	for i, v := range block {
		j := i * 8
		dst[j] = byte(v)
		dst[j+1] = byte(v >> 8)
		dst[j+2] = byte(v >> 16)
		dst[j+3] = byte(v >> 24)
		dst[j+4] = byte(v >> 32)
		dst[j+5] = byte(v >> 40)
		dst[j+6] = byte(v >> 48)
		dst[j+7] = byte(v >> 56)
	}
}
