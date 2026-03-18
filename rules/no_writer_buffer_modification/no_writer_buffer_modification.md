---
title: noWriterBufferModification
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noWriterBufferModification`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noWriterBufferModification:
    # rule options here
```

## Details

Modifying the buffer in an `io.Writer` implementation.

An `io.Writer` implementation must not modify the slice data, even temporarily. The caller retains ownership of the slice, and modifying it violates the `io.Writer` contract and can lead to data corruption.

Source: https://staticcheck.dev/docs/checks/#SA1023

## Examples

### Invalid

```golang
package main

type myWriter struct{}

func (w myWriter) Write(p []byte) (int, error) {
    // Wrong: modifying the input slice
    for i := range p {
        p[i] = p[i] ^ 0xFF
    }
    return len(p), nil
}
```

### Valid

```golang
package main

type myWriter struct{}

func (w myWriter) Write(p []byte) (int, error) {
    // Correct: create a copy if modification is needed
    buf := make([]byte, len(p))
    copy(buf, p)
    for i := range buf {
        buf[i] = buf[i] ^ 0xFF
    }
    return len(p), nil
}
```
