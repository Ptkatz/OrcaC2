// Copyright (c) 2016 Andreas Auernhammer. All rights reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

// +build amd64

package threefish

import "unsafe"

func bytesToBlock256(block *[4]uint64, src []byte) {
	srcPtr := (*[4]uint64)(unsafe.Pointer(&src[0]))

	block[0] = srcPtr[0]
	block[1] = srcPtr[1]
	block[2] = srcPtr[2]
	block[3] = srcPtr[3]
}

func block256ToBytes(dst []byte, block *[4]uint64) {
	dstPtr := (*[4]uint64)(unsafe.Pointer(&dst[0]))

	dstPtr[0] = block[0]
	dstPtr[1] = block[1]
	dstPtr[2] = block[2]
	dstPtr[3] = block[3]
}

func bytesToBlock512(block *[8]uint64, src []byte) {
	srcPtr := (*[8]uint64)(unsafe.Pointer(&src[0]))

	block[0] = srcPtr[0]
	block[1] = srcPtr[1]
	block[2] = srcPtr[2]
	block[3] = srcPtr[3]
	block[4] = srcPtr[4]
	block[5] = srcPtr[5]
	block[6] = srcPtr[6]
	block[7] = srcPtr[7]
}

func block512ToBytes(dst []byte, block *[8]uint64) {
	dstPtr := (*[8]uint64)(unsafe.Pointer(&dst[0]))

	dstPtr[0] = block[0]
	dstPtr[1] = block[1]
	dstPtr[2] = block[2]
	dstPtr[3] = block[3]
	dstPtr[4] = block[4]
	dstPtr[5] = block[5]
	dstPtr[6] = block[6]
	dstPtr[7] = block[7]
}

func bytesToBlock1024(block *[16]uint64, src []byte) {
	srcPtr := (*[16]uint64)(unsafe.Pointer(&src[0]))

	block[0] = srcPtr[0]
	block[1] = srcPtr[1]
	block[2] = srcPtr[2]
	block[3] = srcPtr[3]
	block[4] = srcPtr[4]
	block[5] = srcPtr[5]
	block[6] = srcPtr[6]
	block[7] = srcPtr[7]
	block[8] = srcPtr[8]
	block[9] = srcPtr[9]
	block[10] = srcPtr[10]
	block[11] = srcPtr[11]
	block[12] = srcPtr[12]
	block[13] = srcPtr[13]
	block[14] = srcPtr[14]
	block[15] = srcPtr[15]
}

func block1024ToBytes(dst []byte, block *[16]uint64) {
	dstPtr := (*[16]uint64)(unsafe.Pointer(&dst[0]))

	dstPtr[0] = block[0]
	dstPtr[1] = block[1]
	dstPtr[2] = block[2]
	dstPtr[3] = block[3]
	dstPtr[4] = block[4]
	dstPtr[5] = block[5]
	dstPtr[6] = block[6]
	dstPtr[7] = block[7]
	dstPtr[8] = block[8]
	dstPtr[9] = block[9]
	dstPtr[10] = block[10]
	dstPtr[11] = block[11]
	dstPtr[12] = block[12]
	dstPtr[13] = block[13]
	dstPtr[14] = block[14]
	dstPtr[15] = block[15]
}
