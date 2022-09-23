// HEAD_PLACEHOLDER
// +build ignore

// Package groestl implements Grøstl-256 algorithm.
//
// This Go implementation is a port of the original C implementation which is
// included in Monero as follows:
//     src/crypto/groestl.c
//     src/crypto/groestl.h
//     src/crypto/groestl_tables.h
//
// Most comments in the original file are copied as well.
//
// In this implementation, we assume all bytes are full.
package groestl

import (
	"hash"
	"unsafe"
)

// This field is for macro definitions.
// We define it in a literal string so that it can trick gofmt(1).
//
// It should be empty after they are expanded by cpp(1).
const _ = `
#undef build
#undef ignore

#define U8_U32(a, begin, end) \
    ( (*[( (end) - (begin) ) / 4]uint32)(unsafe.Pointer(&a[ (begin) ])) )

#define U32_U8(a, begin, end) \
    ( (*[( (end) - (begin) ) * 4]uint8)(unsafe.Pointer(&a[ (begin) ])) )

#define COLUMN(x, y, i, c0, c1, c2, c3, c4, c5, c6, c7, tv1, tv2, tu, tl, t) \
	tu = tab[2*uint32(x[4*c0+0])];				\
	tl = tab[2*uint32(x[4*c0+0])+1];			\
	tv1 = tab[2*uint32(x[4*c1+1])];				\
	tv2 = tab[2*uint32(x[4*c1+1])+1];			\
	ROTATE_COLUMN_DOWN(tv1, tv2, 1, t);			\
	tu ^= tv1;									\
	tl ^= tv2;									\
	tv1 = tab[2*uint32(x[4*c2+2])];				\
	tv2 = tab[2*uint32(x[4*c2+2])+1];			\
	ROTATE_COLUMN_DOWN(tv1, tv2, 2, t);			\
	tu ^= tv1;									\
	tl ^= tv2;									\
	tv1 = tab[2*uint32(x[4*c3+3])];				\
	tv2 = tab[2*uint32(x[4*c3+3])+1];			\
	ROTATE_COLUMN_DOWN(tv1, tv2, 3, t);			\
	tu ^= tv1;									\
	tl ^= tv2;									\
	tl ^= tab[2*uint32(x[4*c4+0])];				\
	tu ^= tab[2*uint32(x[4*c4+0])+1];			\
	tv1 = tab[2*uint32(x[4*c5+1])];				\
	tv2 = tab[2*uint32(x[4*c5+1])+1];			\
	ROTATE_COLUMN_DOWN(tv1, tv2, 1, t);			\
	tl ^= tv1;									\
	tu ^= tv2;									\
	tv1 = tab[2*uint32(x[4*c6+2])];				\
	tv2 = tab[2*uint32(x[4*c6+2])+1];			\
	ROTATE_COLUMN_DOWN(tv1, tv2, 2, t);			\
	tl ^= tv1;									\
	tu ^= tv2;									\
	tv1 = tab[2*uint32(x[4*c7+3])];				\
	tv2 = tab[2*uint32(x[4*c7+3])+1];			\
	ROTATE_COLUMN_DOWN(tv1, tv2, 3, t);			\
	tl ^= tv1;									\
	tu ^= tv2;									\
	y[i] = tu;									\
	y[i+1] = tl;

#define ROTATE_COLUMN_DOWN(v1, v2, amountBytes, tempVar) \
	tempVar = (v1 << (8 * amountBytes)) | (v2 >> (8 * (4 - amountBytes)));	\
	v2 = (v2 << (8 * amountBytes)) | (v1 >> (8 * (4 - amountBytes)));		\
	v1 = tempVar;
`

// To trick goimports(1)
var _ = unsafe.Pointer(nil)

const (
	rows           = 8
	cols512        = 8
	size512        = rows * cols512
	lengthFieldLen = rows

	hashBitLen  = 256
	hashByteLen = 32
)

var (
	zeroBuf64Byte [size512]byte
	zeroBuf64Word [size512 / 4]uint32
)

type state struct {
	chaining [size512 / 4]uint32 // actual state

	blockCounter1,
	blockCounter2 uint32 // message block counter(s)

	buffer [size512]byte // data buffer
	bufPtr int           // data buffer pointer
}

func Sum256(b []byte) []byte {
	h := New256()
	h.Write(b)

	return h.Sum(nil)
}

func New256() hash.Hash {
	s := &state{}
	s.chaining[2*cols512-1] = 65536

	return s
}

func (s *state) Reset() {
	s.chaining = zeroBuf64Word
	s.chaining[2*cols512-1] = 65536
	s.blockCounter1 = 0
	s.blockCounter2 = 0
	s.buffer = zeroBuf64Byte
	s.bufPtr = 0
}

func (s *state) Size() int      { return hashByteLen }
func (s *state) BlockSize() int { return size512 }

// Write updates state with databitlen bits of input
func (s *state) Write(data []byte) (n int, err error) {
	n = len(data)
	index := 0

	// if the buffer contains data that has not yet been digested, first
	// add data to buffer until full
	if s.bufPtr > 0 {
		m := copy(s.buffer[s.bufPtr:], data)
		s.bufPtr += m
		index += m
		if s.bufPtr < size512 {
			// buffer still not full, return
			return
		}

		// digest buffer
		s.bufPtr = 0
		s.transform(s.buffer[:])
	}

	// digest bulk of message
	s.transform(data[index:])
	index += (n - index) / size512 * size512

	// store remaining data in buffer
	m := copy(s.buffer[:], data[index:])
	s.bufPtr += m
	index += m

	return
}

// Sum process remaining data (including padding), perform
// output transformation.
func (s *state) Sum(b []byte) []byte {
	s.buffer[s.bufPtr] = 0x80
	s.bufPtr++

	// pad with '0'-bits
	if s.bufPtr > size512-lengthFieldLen {
		// padding requires two blocks
		n := copy(s.buffer[s.bufPtr:size512], zeroBuf64Byte[:])
		s.bufPtr += n
		// digest first padding block
		s.transform(s.buffer[:size512])
		s.bufPtr = 0
	}
	n := copy(s.buffer[s.bufPtr:size512-lengthFieldLen], zeroBuf64Byte[:])
	s.bufPtr += n

	// length padding
	s.blockCounter1++
	if s.blockCounter1 == 0 {
		s.blockCounter2++
	}
	s.bufPtr = size512

	for s.bufPtr > size512-4 {
		s.bufPtr--
		s.buffer[s.bufPtr] = uint8(s.blockCounter1)
		s.blockCounter1 >>= 8
	}
	for s.bufPtr > size512-lengthFieldLen {
		s.bufPtr--
		s.buffer[s.bufPtr] = uint8(s.blockCounter2)
		s.blockCounter2 >>= 8
	}
	// digest final padding block
	s.transform(s.buffer[:size512])
	// perform output transformation
	s.outputTransformation()

	// store hash result
	return append(b, U32_U8(s.chaining, hashByteLen/4, size512/4)[:]...)
}

// digest up to msglen bytes of input (full blocks only)
func (s *state) transform(b []byte) {
	n := len(b)
	offset := 0

	// digest message, one block at a time
	for n >= size512 {
		input := b[offset:]
		// length of input is known and constant
		f512(&s.chaining, U8_U32(input, 0, size512))

		// increment block counter
		s.blockCounter1++
		if s.blockCounter1 == 0 {
			s.blockCounter2++
		}

		n -= size512
		offset += size512
	}
}

// given state h, do h <- P(h)+h
func (s *state) outputTransformation() {
	var j int
	var temp, y, z [2 * cols512]uint32

	for j = 0; j < 2*cols512; j++ {
		temp[j] = s.chaining[j]
	}
	rnd512p(U32_U8(temp, 0, 2*cols512), &y, 0x00000000)
	rnd512p(U32_U8(y, 0, 2*cols512), &z, 0x00000001)
	rnd512p(U32_U8(z, 0, 2*cols512), &y, 0x00000002)
	rnd512p(U32_U8(y, 0, 2*cols512), &z, 0x00000003)
	rnd512p(U32_U8(z, 0, 2*cols512), &y, 0x00000004)
	rnd512p(U32_U8(y, 0, 2*cols512), &z, 0x00000005)
	rnd512p(U32_U8(z, 0, 2*cols512), &y, 0x00000006)
	rnd512p(U32_U8(y, 0, 2*cols512), &z, 0x00000007)
	rnd512p(U32_U8(z, 0, 2*cols512), &y, 0x00000008)
	rnd512p(U32_U8(y, 0, 2*cols512), &temp, 0x00000009)
	for j = 0; j < 2*cols512; j++ {
		s.chaining[j] ^= temp[j]
	}
}

// compute compression function (short variants)
func f512(h *[16]uint32, m *[size512 / 4]uint32) {
	var i int
	var Ptmp, Qtmp, y, z [2 * cols512]uint32

	for i = 0; i < 2*cols512; i++ {
		z[i] = m[i]
		Ptmp[i] = h[i] ^ m[i]
	}

	// compute Q(m)
	rnd512q(U32_U8(z, 0, 2*cols512), &y, 0x00000000)
	rnd512q(U32_U8(y, 0, 2*cols512), &z, 0x01000000)
	rnd512q(U32_U8(z, 0, 2*cols512), &y, 0x02000000)
	rnd512q(U32_U8(y, 0, 2*cols512), &z, 0x03000000)
	rnd512q(U32_U8(z, 0, 2*cols512), &y, 0x04000000)
	rnd512q(U32_U8(y, 0, 2*cols512), &z, 0x05000000)
	rnd512q(U32_U8(z, 0, 2*cols512), &y, 0x06000000)
	rnd512q(U32_U8(y, 0, 2*cols512), &z, 0x07000000)
	rnd512q(U32_U8(z, 0, 2*cols512), &y, 0x08000000)
	rnd512q(U32_U8(y, 0, 2*cols512), &Qtmp, 0x09000000)

	// compute P(h+m)
	rnd512p(U32_U8(Ptmp, 0, 2*cols512), &y, 0x00000000)
	rnd512p(U32_U8(y, 0, 2*cols512), &z, 0x00000001)
	rnd512p(U32_U8(z, 0, 2*cols512), &y, 0x00000002)
	rnd512p(U32_U8(y, 0, 2*cols512), &z, 0x00000003)
	rnd512p(U32_U8(z, 0, 2*cols512), &y, 0x00000004)
	rnd512p(U32_U8(y, 0, 2*cols512), &z, 0x00000005)
	rnd512p(U32_U8(z, 0, 2*cols512), &y, 0x00000006)
	rnd512p(U32_U8(y, 0, 2*cols512), &z, 0x00000007)
	rnd512p(U32_U8(z, 0, 2*cols512), &y, 0x00000008)
	rnd512p(U32_U8(y, 0, 2*cols512), &Ptmp, 0x00000009)

	// compute P(h+m) + Q(m) + h
	for i = 0; i < 2*cols512; i++ {
		h[i] ^= Ptmp[i] ^ Qtmp[i]
	}
}

// compute one round of Q (short variants)
func rnd512q(x *[64]byte, y *[16]uint32, r uint32) {
	var temp1, temp2, tempUpperValue, tempLowerValue, temp uint32
	x32 := U8_U32(x, 0, 64)
	x32[0] = ^x32[0]
	x32[1] ^= 0xffffffff ^ r
	x32[2] = ^x32[2]
	x32[3] ^= 0xefffffff ^ r
	x32[4] = ^x32[4]
	x32[5] ^= 0xdfffffff ^ r
	x32[6] = ^x32[6]
	x32[7] ^= 0xcfffffff ^ r
	x32[8] = ^x32[8]
	x32[9] ^= 0xbfffffff ^ r
	x32[10] = ^x32[10]
	x32[11] ^= 0xafffffff ^ r
	x32[12] = ^x32[12]
	x32[13] ^= 0x9fffffff ^ r
	x32[14] = ^x32[14]
	x32[15] ^= 0x8fffffff ^ r
	COLUMN(x, y, 0, 2, 6, 10, 14, 1, 5, 9, 13, temp1, temp2, tempUpperValue, tempLowerValue, temp)
	COLUMN(x, y, 2, 4, 8, 12, 0, 3, 7, 11, 15, temp1, temp2, tempUpperValue, tempLowerValue, temp)
	COLUMN(x, y, 4, 6, 10, 14, 2, 5, 9, 13, 1, temp1, temp2, tempUpperValue, tempLowerValue, temp)
	COLUMN(x, y, 6, 8, 12, 0, 4, 7, 11, 15, 3, temp1, temp2, tempUpperValue, tempLowerValue, temp)
	COLUMN(x, y, 8, 10, 14, 2, 6, 9, 13, 1, 5, temp1, temp2, tempUpperValue, tempLowerValue, temp)
	COLUMN(x, y, 10, 12, 0, 4, 8, 11, 15, 3, 7, temp1, temp2, tempUpperValue, tempLowerValue, temp)
	COLUMN(x, y, 12, 14, 2, 6, 10, 13, 1, 5, 9, temp1, temp2, tempUpperValue, tempLowerValue, temp)
	COLUMN(x, y, 14, 0, 4, 8, 12, 15, 3, 7, 11, temp1, temp2, tempUpperValue, tempLowerValue, temp)
}

// compute one round of P (short variants)
func rnd512p(x *[64]byte, y *[16]uint32, r uint32) {
	var temp1, temp2, tempUpperValue, tempLowerValue, temp uint32
	x32 := U8_U32(x, 0, 64)
	x32[0] ^= 0x00000000 ^ r
	x32[2] ^= 0x00000010 ^ r
	x32[4] ^= 0x00000020 ^ r
	x32[6] ^= 0x00000030 ^ r
	x32[8] ^= 0x00000040 ^ r
	x32[10] ^= 0x00000050 ^ r
	x32[12] ^= 0x00000060 ^ r
	x32[14] ^= 0x00000070 ^ r
	COLUMN(x, y, 0, 0, 2, 4, 6, 9, 11, 13, 15, temp1, temp2, tempUpperValue, tempLowerValue, temp)
	COLUMN(x, y, 2, 2, 4, 6, 8, 11, 13, 15, 1, temp1, temp2, tempUpperValue, tempLowerValue, temp)
	COLUMN(x, y, 4, 4, 6, 8, 10, 13, 15, 1, 3, temp1, temp2, tempUpperValue, tempLowerValue, temp)
	COLUMN(x, y, 6, 6, 8, 10, 12, 15, 1, 3, 5, temp1, temp2, tempUpperValue, tempLowerValue, temp)
	COLUMN(x, y, 8, 8, 10, 12, 14, 1, 3, 5, 7, temp1, temp2, tempUpperValue, tempLowerValue, temp)
	COLUMN(x, y, 10, 10, 12, 14, 0, 3, 5, 7, 9, temp1, temp2, tempUpperValue, tempLowerValue, temp)
	COLUMN(x, y, 12, 12, 14, 0, 2, 5, 7, 9, 11, temp1, temp2, tempUpperValue, tempLowerValue, temp)
	COLUMN(x, y, 14, 14, 0, 2, 4, 7, 9, 11, 13, temp1, temp2, tempUpperValue, tempLowerValue, temp)
}
