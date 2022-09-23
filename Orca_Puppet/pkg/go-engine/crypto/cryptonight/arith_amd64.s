#include "textflag.h"
#include "sum_defs_amd64.h"

// func mul128(x, y uint64) (lo, hi uint64)
TEXT ·mul128(SB), NOSPLIT, $0
	MOVQ    x+0(FP), AX
	MULQ    y+8(FP)
	MOVQ    AX, lo+16(FP)
	MOVQ    DX, hi+24(FP)
	RET

// func v2Sqrt(in uint64) (out uint64)
TEXT ·v2Sqrt(SB), NOSPLIT, $0
	MOVQ    in+0(FP), TMP1

	// <BEGIN> VARIANT2_INTEGER_MATH_SQRT_STEP
	MOVQ    TMP1, AX
	SHRQ    $12, AX
	MOVQ    $(1023 << 52), BX
	ADDQ    BX, AX
	MOVQ    AX, TMPX0
	SQRTSD  TMPX0, TMPX0
	MOVQ    TMPX0, TMP2
	SUBQ    BX, TMP2
	SHRQ    $19, TMP2        // not yet sanitized sqrt result
	// <END> VARIANT2_INTEGER_MATH_SQRT_STEP
	// <BEGIN> VARIANT2_INTEGER_MATH_SQRT_FIXUP
	MOVQ    TMP2, AX
	SHRQ    $1, AX            // s = sqrtResult >> 1
	MOVQ    TMP2, BX
	ANDQ    $1, BX            // b = sqrtResult & 1
	MOVQ    TMP2, CX
	SHLQ    $32, CX
	LEAQ    0(AX)(BX*1), DX
	IMULQ   AX, DX
	ADDQ    DX, CX            // r2 = s * (s + b) + (sqrtResult << 32)

	ADDQ    CX, BX
	XORQ    DX, DX
	CMPQ    BX, TMP1
	SETHI   DL
	SUBQ    DX, TMP2         // sqrtResult += ((r2 + b > sqrtInput) ? -1 : 0)
	// NOTE: the following branch does not seem to be able to be covered,
	//   i.e. it works without the code below.
	//   In case you find any issue, try de-commenting these.

	// MOVQ    $0x100000000, DX
	// LEAQ    0(CX)(DX*1), BX
	// SUBQ    AX, TMP1
	// XORQ    DX, DX
	// CMPQ    BX, TMP1
	// SETCS   DL
	// ADDQ    DX, TMP2      // sqrtResult += ((r2 + (1 << 32) < sqrtInput - s) ? 1 : 0)

	// <END> VARIANT2_INTEGER_MATH_SQRT_FIXUP

	MOVQ    TMP2, out+8(FP)
	RET
