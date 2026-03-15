// Assembly stubs for testing noAsmDeclMismatch rule.

// Wrong: should be $0-24 for two int64 args + int64 return (8+8+8=24)
TEXT ·Add(SB), NOSPLIT, $0-16
    MOVQ x+0(FP), AX
    ADDQ y+8(FP), AX
    MOVQ AX, ret+16(FP)
    RET

// Correct: $0-24 matches two int64 args + int64 return (8+8+8=24)
TEXT ·Sub(SB), NOSPLIT, $0-24
    MOVQ x+0(FP), AX
    SUBQ y+8(FP), AX
    MOVQ AX, ret+16(FP)
    RET
