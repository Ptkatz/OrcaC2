// spec
#define PAD_SIZE (2 << 20) // 2 MiB
#define ITER     (1 << 19) // 524288

// common
#define STATE R8
#define I     R9
#define CHUNK R10
#define A     X0
#define B     X1
#define C     X2
#define D     X3
#define TMP0  R15
#define TMP1  R14
#define TMP2  R13
#define TMPX0 X15
#define TMPX1 X14
#define TMPX2 X13
#define TMPX3 X12

// variant 1
#define TWEAK X4

// variant 2
#define E     X4
#define DIV_RESULT  R11
#define SQRT_RESULT R12
