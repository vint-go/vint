---
title: noUnboundedDecompression
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/security/noUnboundedDecompression`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/security/noUnboundedDecompression:
    # rule options here
```

## Details

Detects the use of `io.Copy` instead of `io.CopyN` when decompressing data, which can lead to decompression bomb (zip bomb) attacks.

A decompression bomb is a maliciously crafted compressed file that expands to an enormous size when decompressed. When `io.Copy` is used to read from a decompression reader (such as `gzip.Reader`, `zlib.Reader`, `flate.Reader`, or `bzip2.Reader`), it will read until EOF without any size limit, potentially consuming all available memory or disk space.

Using `io.CopyN` instead allows specifying a maximum number of bytes to copy, preventing resource exhaustion attacks. Alternatively, wrapping the reader with `io.LimitReader` provides similar protection.

Source: https://github.com/securego/gosec

## Examples

### Invalid

```golang
import (
    "compress/gzip"
    "io"
    "os"
)

func decompress(src io.Reader) error {
    gz, err := gzip.NewReader(src)
    if err != nil {
        return err
    }
    defer gz.Close()

    // Unbounded copy from decompression reader
    _, err = io.Copy(os.Stdout, gz)
    return err
}
```

```golang
import (
    "compress/zlib"
    "io"
)

func decompressZlib(data []byte) ([]byte, error) {
    r, err := zlib.NewReader(bytes.NewReader(data))
    if err != nil {
        return nil, err
    }
    defer r.Close()

    var buf bytes.Buffer
    _, err = io.Copy(&buf, r) // Potential decompression bomb
    return buf.Bytes(), err
}
```

### Valid

```golang
import (
    "compress/gzip"
    "io"
    "os"
)

func decompress(src io.Reader) error {
    gz, err := gzip.NewReader(src)
    if err != nil {
        return err
    }
    defer gz.Close()

    // Bounded copy with maximum size
    _, err = io.CopyN(os.Stdout, gz, 1024*1024*100) // Max 100MB
    return err
}
```

```golang
import (
    "compress/gzip"
    "io"
)

func decompress(src io.Reader) ([]byte, error) {
    gz, err := gzip.NewReader(src)
    if err != nil {
        return nil, err
    }
    defer gz.Close()

    // Using LimitReader to bound the decompression
    limited := io.LimitReader(gz, 1024*1024*100)
    return io.ReadAll(limited)
}
```
