---
title: noDeepExit
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noDeepExit`
- This rule is not recommended, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noDeepExit:
    # no configuration options
```

## Details

Packages exposing functions that can stop program execution by exiting are hard to reuse.
This rule looks for program exits in functions other than `main()` or `init()`.

The rule detects calls to the following exit functions outside of `main()` and `init()`:
- `os.Exit`
- `syscall.Exit`
- `log.Fatal`, `log.Fatalf`, `log.Fatalln`
- `log.Panic`, `log.Panicf`, `log.Panicln`
- `flag.Parse`
- `flag.NewFlagSet` with `flag.ExitOnError`

In test files, `TestMain` and testable example functions (e.g. `Example`, `ExampleFoo`) are also excluded from checking.

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
package mypackage

import "os"

func foo() {
	os.Exit(1) // calls to os.Exit only in main() or init() functions
}
```

```golang
package mypackage

import "log"

func bar() {
	log.Fatal("something went wrong") // calls to log.Fatal only in main() or init() functions
}
```

```golang
package mypackage

import "flag"

func setup() {
	flag.Parse() // calls to flag.Parse only in main() or init() functions
}
```

### Valid

```golang
package main

import "os"

func main() {
	os.Exit(0) // ok, in main()
}

func init() {
	os.Exit(1) // ok, in init()
}
```

```golang
package mypackage

import "log"

func bar() {
	log.Println("this is fine") // not an exit function
}
```
