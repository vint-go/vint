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

**Note:** This rule only applies to projects targeting Go < 1.15. Since Go 1.15, the test runner automatically uses the result of `m.Run()` as the exit code, making explicit `os.Exit` calls unnecessary. On Go >= 1.15 projects, this rule is suppressed and the inverse rule `noRedundantTestMainExit` applies instead.

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
