// Assembly stubs for testing noFramePointerClobber rule.

// Bad: clobbers BP without saving it first (frameless function).
TEXT ·badExample(SB), NOSPLIT, $0-8
    MOVQ $0, BP
    MOVQ BP, ret+0(FP)
    RET

// Good: saves and restores BP properly.
TEXT ·goodExample(SB), NOSPLIT, $0-8
    MOVQ BP, saved-8(SP)
    LEAQ saved-8(SP), BP
    MOVQ $42, ret+0(FP)
    MOVQ saved-8(SP), BP
    RET
