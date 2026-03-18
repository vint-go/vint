---
title: noBenchmarkNAssignment
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noBenchmarkNAssignment`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noBenchmarkNAssignment:
    # rule options here
```

## Details

Assigning to `b.N` in benchmarks distorts the results.

In benchmark functions, `b.N` is controlled by the testing framework to determine how many iterations to run. Assigning to `b.N` directly will distort the benchmark results and break the benchmarking loop.

Source: https://staticcheck.dev/docs/checks/#SA3001

## Examples

### Invalid

```golang
package main

import "testing"

func BenchmarkSomething(b *testing.B) {
    // Wrong: assigning to b.N
    b.N = 1000
    for i := 0; i < b.N; i++ {
        // do work
    }
}
```

### Valid

```golang
package main

import "testing"

func BenchmarkSomething(b *testing.B) {
    // Correct: let the framework control b.N
    for i := 0; i < b.N; i++ {
        // do work
    }
}
```
