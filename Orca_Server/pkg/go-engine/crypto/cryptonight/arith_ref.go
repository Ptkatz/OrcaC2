// +build !amd64

package cryptonight

import (
	"math"
)

func mul128(x, y uint64) (lo, hi uint64) {
	xhi, yhi := x>>32, y>>32
	xlo, ylo := x&0xffffffff, y&0xffffffff

	hihi := xhi * yhi
	lolo := xlo * ylo
	lohi := xlo * yhi
	hilo := xhi * ylo

	mid := lolo>>32 + lohi&0xffffffff + hilo&0xffffffff
	lo = mid<<32 | (lolo & 0xffffffff)
	hi = hihi + lohi>>32 + hilo>>32 + mid>>32

	return
}

func v2Sqrt(in uint64) uint64 {
	out := uint64(
		math.Sqrt(
			float64(in)+1<<64,
		)*2 - 1<<33,
	)

	s := out >> 1
	b := out & 1
	r := s*(s+b) + (out << 32)
	if r+b > in {
		out--
	}

	// NOTE: the following branch does not seem to be able to be covered,
	//   i.e. it works without the code below.
	//   In case you find any issue, try de-commenting these.

	// if r+1<<32 < in-s {
	// 	out++
	// }

	return out
}
