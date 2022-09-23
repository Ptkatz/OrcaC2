#include "textflag.h"

// func CnRoundsAsm(dst, src *uint64, rkeys *uint32)
TEXT ·CnRoundsAsm(SB), NOSPLIT, $0
	MOVQ    dst+0(FP), AX
	MOVQ    src+8(FP), BX
	MOVQ    rkeys+16(FP), CX
	MOVO    0(BX), X0
	MOVO    0(CX), X1
	AESENC  X1, X0
	MOVO    16(CX), X1
	AESENC  X1, X0
	MOVO    32(CX), X1
	AESENC  X1, X0
	MOVO    48(CX), X1
	AESENC  X1, X0
	MOVO    64(CX), X1
	AESENC  X1, X0
	MOVO    80(CX), X1
	AESENC  X1, X0
	MOVO    96(CX), X1
	AESENC  X1, X0
	MOVO    112(CX), X1
	AESENC  X1, X0
	MOVO    128(CX), X1
	AESENC  X1, X0
	MOVO    144(CX), X1
	AESENC  X1, X0
	MOVO    X0, 0(AX)
	RET

// func CnExpandKeyAsm(key *uint64, rkey *uint32)
// Note that round keys are stored in uint128 format, not uint32
TEXT ·CnExpandKeyAsm(SB), NOSPLIT, $0
	MOVQ    key+0(FP), AX
	MOVQ    rkey+8(FP), BX
	MOVO    (AX), X0
	MOVO    X0, (BX)
	ADDQ    $16, BX
	PXOR    X4, X4 // _expand_key_* expect X4 to be zero

	MOVO    16(AX), X2
	MOVO    X2, (BX)
	ADDQ    $16, BX
	AESKEYGENASSIST $0x01, X2, X1
	CALL    _expand_key_256a<>(SB)
	AESKEYGENASSIST $0x01, X0, X1
	CALL    _expand_key_256b<>(SB)
	AESKEYGENASSIST $0x02, X2, X1
	CALL    _expand_key_256a<>(SB)
	AESKEYGENASSIST $0x02, X0, X1
	CALL    _expand_key_256b<>(SB)
	AESKEYGENASSIST $0x04, X2, X1
	CALL    _expand_key_256a<>(SB)
	AESKEYGENASSIST $0x04, X0, X1
	CALL    _expand_key_256b<>(SB)
	AESKEYGENASSIST $0x08, X2, X1
	CALL    _expand_key_256a<>(SB)
	AESKEYGENASSIST $0x08, X0, X1
	CALL    _expand_key_256b<>(SB)
	AESKEYGENASSIST $0x10, X2, X1
	CALL    _expand_key_256a<>(SB)
	RET

TEXT _expand_key_256a<>(SB), NOSPLIT, $0
	PSHUFD  $0xff, X1, X1
	SHUFPS  $0x10, X0, X4
	PXOR    X4, X0
	SHUFPS  $0x8c, X0, X4
	PXOR    X4, X0
	PXOR    X1, X0
	MOVO    X0, (BX)
	ADDQ    $16, BX
	RET

TEXT _expand_key_256b<>(SB), NOSPLIT, $0
	PSHUFD  $0xaa, X1, X1
	SHUFPS  $0x10, X2, X4
	PXOR    X4, X2
	SHUFPS  $0x8c, X2, X4
	PXOR    X4, X2
	PXOR    X1, X2
	MOVO    X2, (BX)
	ADDQ    $16, BX
	RET

// func CnSingleRoundAsm(dst, src *uint64, rkeys *uint64)
TEXT ·CnSingleRoundAsm(SB), NOSPLIT, $0
	MOVQ    dst+0(FP), AX
	MOVQ    src+8(FP), BX
	MOVQ    rkeys+16(FP), CX
	MOVO    0(BX), X0
	MOVO    0(CX), X1
	AESENC  X1, X0
	MOVO    X0, 0(AX)
	RET
