---
title: noInvalidUrlParse
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noInvalidUrlParse`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noInvalidUrlParse:
    # rule options here
```

## Details

Invalid URL in `net/url.Parse`.

This check detects URLs passed to `net/url.Parse` that are obviously invalid, such as URLs with backslashes or other malformed components. While `url.Parse` is very lenient, some inputs are clearly not valid URLs.

Source: https://staticcheck.dev/docs/checks/#SA1007

## Examples

### Invalid

```golang
package main

import "net/url"

func main() {
    // Backslashes are not valid in URLs
    u, _ := url.Parse("http:\\\\example.com")
    _ = u
}
```

### Valid

```golang
package main

import "net/url"

func main() {
    // Correct URL format with forward slashes
    u, _ := url.Parse("http://example.com")
    _ = u
}
```
