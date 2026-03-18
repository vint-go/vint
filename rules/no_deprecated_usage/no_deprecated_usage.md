---
title: noDeprecatedUsage
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noDeprecatedUsage`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noDeprecatedUsage:
    # rule options here
```

## Details

Using a deprecated function, variable, constant, or field.

This check flags the usage of deprecated identifiers from the standard library and third-party packages. Deprecated identifiers are marked with a `// Deprecated:` comment in their documentation. You should use the suggested replacement instead.

Source: https://staticcheck.dev/docs/checks/#SA1019

## Examples

### Invalid

```golang
package main

import "io/ioutil"

func main() {
    // ioutil.ReadAll is deprecated since Go 1.16
    data, _ := ioutil.ReadAll(nil)
    _ = data
}
```

### Valid

```golang
package main

import "io"

func main() {
    // Use io.ReadAll instead
    data, _ := io.ReadAll(nil)
    _ = data
}
```
