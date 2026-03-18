---
title: noNonOctalFileMode
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noNonOctalFileMode`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noNonOctalFileMode:
    # rule options here
```

## Details

Using a non-octal `os.FileMode` that looks like it was meant to be in octal.

File permissions in Go should be specified in octal (e.g., `0644`). Writing `644` (without the leading zero) is decimal 644, which is octal 1204 -- likely not what you intended.

Source: https://staticcheck.dev/docs/checks/#SA9002

## Examples

### Invalid

```golang
package main

import "os"

func main() {
    // 644 decimal = 01204 octal, not 0644
    os.WriteFile("file.txt", []byte("data"), 644)
}
```

### Valid

```golang
package main

import "os"

func main() {
    // Octal notation for file permissions
    os.WriteFile("file.txt", []byte("data"), 0644)
}
```
