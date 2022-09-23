// Copyright (c) 2016 Andreas Auernhammer. All rights reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package threefish

import "testing"

// The UBI256, UBI512 and UBI1024 functions are tested within
// the skein packages (skein, skein256 and skein1024)

func testBlockSize(t *testing.T, blocksize int) {
	var tweak [TweakSize]byte
	c, err := NewCipher(&tweak, make([]byte, blocksize))
	if err != nil {
		t.Fatalf("Failed to create Threefish-%d instance: %s", blocksize*8, err)
	}

	if bs := c.BlockSize(); bs != blocksize {
		t.Fatalf("BlockSize() returned unexpected value: %d - expected %d", bs, blocksize)
	}
}

func TestBlockSize(t *testing.T) {
	testBlockSize(t, BlockSize256)
	testBlockSize(t, BlockSize512)
	testBlockSize(t, BlockSize1024)
}

func TestNew(t *testing.T) {
	badKeyLengths := []int{
		0, 31, 33, 63, 65, 127, 129,
	}
	var tweak [TweakSize]byte
	for i, v := range badKeyLengths {
		_, err := NewCipher(&tweak, make([]byte, v))
		if err == nil {
			t.Fatalf("BadKey %d:  NewCipher accepted inavlid key length %d", i, v)
		}
	}
}

func TestIncrementTweak(t *testing.T) {
	var tweak [3]uint64

	IncrementTweak(&tweak, 1)
	if tweak[0] != 1 {
		t.Fatalf("IncrementTweak failed by increment of %d", 1)
	}

	tweak[0] = ^uint64(0)
	IncrementTweak(&tweak, 2)
	if tweak[0] != 1 && tweak[1] != 1 {
		t.Fatalf("IncrementTweak failed by increment of %d", 2)
	}

	tweak[0] = ^uint64(0)
	tweak[1] = uint64(0xFFFFFFFF)
	IncrementTweak(&tweak, 1)
	if tweak[0] != 0 && tweak[1] != 0 {
		t.Fatalf("IncrementTweak failed by increment of %d", 1)
	}
}

// Benchmarks

func benchmarkEncrypt(b *testing.B, blocksize, size int) {
	key := make([]byte, blocksize)
	var tweak [TweakSize]byte

	c, err := NewCipher(&tweak, key)
	if err != nil {
		b.Fatalf("Failed to create Threefish-%d instance: %s", blocksize*8, err)
	}
	n := size / blocksize
	buf := make([]byte, blocksize)
	b.SetBytes(int64(blocksize * n))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < n; j++ {
			c.Encrypt(buf, buf)
		}
	}
}

func benchmarkDecrypt(b *testing.B, blocksize, size int) {
	key := make([]byte, blocksize)
	var tweak [TweakSize]byte

	c, err := NewCipher(&tweak, key)
	if err != nil {
		b.Fatalf("Failed to create Threefish-%d instance: %s", blocksize*8, err)
	}

	n := size / blocksize
	buf := make([]byte, blocksize)
	b.SetBytes(int64(blocksize * n))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < n; j++ {
			c.Decrypt(buf, buf)
		}
	}
}

func BenchmarkEncrypt256_32(b *testing.B)    { benchmarkEncrypt(b, BlockSize256, 32) }
func BenchmarkEncrypt256_1024(b *testing.B)  { benchmarkEncrypt(b, BlockSize256, 1024) }
func BenchmarkEncrypt512_64(b *testing.B)    { benchmarkEncrypt(b, BlockSize512, 64) }
func BenchmarkEncrypt512_1024(b *testing.B)  { benchmarkEncrypt(b, BlockSize512, 1024) }
func BenchmarkEncrypt1024_128(b *testing.B)  { benchmarkEncrypt(b, BlockSize1024, 128) }
func BenchmarkEncrypt1024_1024(b *testing.B) { benchmarkEncrypt(b, BlockSize1024, 1024) }

func BenchmarkDecrypt256_32(b *testing.B)    { benchmarkDecrypt(b, BlockSize256, 32) }
func BenchmarkDecrypt256_1024(b *testing.B)  { benchmarkDecrypt(b, BlockSize256, 1024) }
func BenchmarkDecrypt512_64(b *testing.B)    { benchmarkDecrypt(b, BlockSize512, 64) }
func BenchmarkDecrypt512_1024(b *testing.B)  { benchmarkDecrypt(b, BlockSize512, 1024) }
func BenchmarkDecrypt1024_128(b *testing.B)  { benchmarkDecrypt(b, BlockSize1024, 128) }
func BenchmarkDecrypt1024_1024(b *testing.B) { benchmarkDecrypt(b, BlockSize1024, 1024) }
