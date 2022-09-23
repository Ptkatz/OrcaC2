// Copyright (c) 2016 Andreas Auernhammer. All rights reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

// +build amd64

package skein

import "unsafe"

func bytesToBlock(block *[8]uint64, src []byte) {
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

func blockToBytes(dst []byte, block *[8]uint64) {
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
