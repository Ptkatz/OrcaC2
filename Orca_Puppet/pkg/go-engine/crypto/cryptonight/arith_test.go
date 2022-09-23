package cryptonight

import (
	"testing"
)

// taken from monero: tests/hash/main.cpp:test_variant2_int_sqrt
//
// comments are reserved as well.
func TestV2Sqrt(t *testing.T) {
	if o := v2Sqrt(0); o != 0 {
		t.Fatalf("expected 0, got %v\n", o)
	}
	if o := v2Sqrt(1 << 63); o != 1930543745 {
		t.Fatalf("expected 1930543745, got %v\n", o)
	}
	if o := v2Sqrt(^uint64(0)); o != 3558067407 {
		t.Fatalf("expected 3558067407, got %v\n", o)
	}

	for i := uint64(1); i <= 3558067407; i++ {
		// "i" is integer part of "sqrt(2^64 + n) * 2 - 2^33"
		// n = (i/2 + 2^32)^2 - 2^64
		i0 := i >> 1
		var n1 uint64
		if (i & 1) == 0 {
			// n = (i/2 + 2^32)^2 - 2^64
			// n = i^2/4 + 2*2^32*i/2 + 2^64 - 2^64
			// n = i^2/4 + 2^32*i
			// i is even, so i^2 is divisible by 4:
			// n = (i^2 >> 2) + (i << 32)
			// int_sqrt_v2(i^2/4 + 2^32*i - 1) must be equal to i - 1
			// int_sqrt_v2(i^2/4 + 2^32*i) must be equal to i
			n1 = i0*i0 + (i << 32) - 1
		} else {
			// n = (i/2 + 2^32)^2 - 2^64
			// n = i^2/4 + 2*2^32*i/2 + 2^64 - 2^64
			// n = i^2/4 + 2^32*i
			// i is odd, so i = i0*2+1 (i0 = i >> 1)
			// n = (i0*2+1)^2/4 + 2^32*i
			// n = (i0^2*4+i0*4+1)/4 + 2^32*i
			// n = i0^2+i0+1/4 + 2^32*i
			// i0^2+i0 + 2^32*i < n < i0^2+i0+1 + 2^32*i
			// int_sqrt_v2(i0^2+i0 + 2^32*i) must be equal to i - 1
			// int_sqrt_v2(i0^2+i0+1 + 2^32*i) must be equal to i
			n1 = i0*i0 + i0 + (i << 32)
		}
		if o := v2Sqrt(n1); o != i-1 {
			t.Fatalf("expected %v, got %v\n", i-1, o)
		}
		if o := v2Sqrt(n1 + 1); o != i {
			t.Fatalf("expected %v, got %v\n", i, o)
		}
	}
}
