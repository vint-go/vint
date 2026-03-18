---
title: noTestMainWithoutExit
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noTestMainWithoutExit`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noTestMainWithoutExit:
    # rule options here
```

## Details

`TestMain` doesn't call `os.Exit`, hiding test failures.

When using a custom `TestMain` function, you must call `os.Exit` with the result of `m.Run()`. Otherwise, the test binary will always exit with code 0, even if tests fail, effectively hiding failures.

Source: https://staticcheck.dev/docs/checks/#SA3000

## Examples

### Invalid

```golang
package main

import "testing"

func TestMain(m *testing.M) {
    // setup
    m.Run()
    // teardown
    // Missing os.Exit - test failures will be hidden
}
```

### Valid

```golang
package main

import (
    "os"
    "testing"
)

func TestMain(m *testing.M) {
    // setup
    code := m.Run()
    // teardown
    os.Exit(code)
}
```
