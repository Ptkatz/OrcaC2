package cryptonight

//go:noescape
func mul128(x, y uint64) (lo, hi uint64)

//go:noescape
func v2Sqrt(in uint64) (out uint64)
