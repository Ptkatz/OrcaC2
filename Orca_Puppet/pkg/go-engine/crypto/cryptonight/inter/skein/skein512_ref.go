// Copyright (c) 2016 Andreas Auernhammer. All rights reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

// +build !amd64

package skein

func bytesToBlock(block *[8]uint64, src []byte) {
	for i := range block {
		j := i * 8
		block[i] = uint64(src[j]) | uint64(src[j+1])<<8 | uint64(src[j+2])<<16 |
			uint64(src[j+3])<<24 | uint64(src[j+4])<<32 | uint64(src[j+5])<<40 |
			uint64(src[j+6])<<48 | uint64(src[j+7])<<56
	}
}

func blockToBytes(dst []byte, block *[8]uint64) {
	i := 0
	for _, v := range block {
		dst[i] = byte(v)
		dst[i+1] = byte(v >> 8)
		dst[i+2] = byte(v >> 16)
		dst[i+3] = byte(v >> 24)
		dst[i+4] = byte(v >> 32)
		dst[i+5] = byte(v >> 40)
		dst[i+6] = byte(v >> 48)
		dst[i+7] = byte(v >> 56)
		i += 8
	}
}
