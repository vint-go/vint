---
title: noExcessiveFileLength
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/complexity/noExcessiveFileLength`
- This rule is not recommended, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/complexity/noExcessiveFileLength:
    arguments:
      - max: 100
        skipComments: true
        skipBlankLines: true
```

## Details

This rule enforces a maximum number of lines per file, in order to aid in maintainability and reduce complexity.

Configuration options:

- `max`: (int) a maximum number of lines in a file. Must be non-negative integers. 0 means the rule is disabled (default `0`).
- `skipComments`: (bool) if true, ignore and do not count lines containing just comments (default `false`).
- `skipBlankLines`: (bool) if true, ignore and do not count lines made up purely of whitespace (default `false`).

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
// With max set to 5, a file with more than 5 lines triggers the rule:
package main

import "fmt"

func main() {
    fmt.Println("Hello")
    fmt.Println("World")
}
// file length is 8 lines, which exceeds the limit of 5
```

### Valid

```golang
// With max set to 10, a short file passes the rule:
package main

import "fmt"

func main() {
    fmt.Println("Hello")
}
```
