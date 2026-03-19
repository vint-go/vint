---
title: noCallToGC
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/performance/noCallToGC`
- This rule is not recommended, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/performance/noCallToGC:
    # rule options here
```

## Details

Explicitly invoking the garbage collector is, except for specific uses in benchmarking, very dubious. The Go runtime is well-optimized for memory management, and manual calls to `runtime.GC()` in production code usually indicate a misunderstanding of Go's garbage collection mechanism. Such calls can degrade performance rather than improve it.

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
package main

import "runtime"

func process() {
	// Explicitly calling the garbage collector is flagged
	runtime.GC()
}
```

### Valid

```golang
package main

import "runtime"

func benchmarkHelper() {
	// In benchmarks, explicit GC calls are acceptable
	// but this rule does not distinguish context
}

func process() {
	// No explicit call to runtime.GC()
	runtime.Goexit()
}
```
