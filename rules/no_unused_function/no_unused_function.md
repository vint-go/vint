---
title: noUnusedFunction
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noUnusedFunction`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noUnusedFunction:
    # No specific options for this rule.
```

## Details

Detects functions that are declared but never called or referenced anywhere in the codebase.

The `unused` linter uses a graph-based approach to trace how code elements are referenced. A function is considered "used" if it is reachable from an entry point such as `main()`, `init()`, an exported symbol, a test function, or a cgo-exported function. If a function is not reachable through any of these paths, it is reported as unused.

Key rules governing function usage:
- Packages use their exported functions (but not exported methods; those are handled via their named type).
- `init` functions are always considered used.
- The `main` function is always used when in the `main` package.
- Functions exported to cgo via `//export` are always considered used.
- Functions linked via `//go:linkname` are always considered used.
- Functions in generated files are considered used when the `generated-is-used` option is enabled (default).
- Anonymous functions defined within another function are considered used by that enclosing function.
- Closures and bound methods are considered used by the functions that reference them.

Source: https://github.com/dominikh/go-tools/tree/master/unused

## Examples

### Invalid

```golang
package main

func main() {
    // greet is never called
}

func greet(name string) string { // func greet is unused
    return "Hello, " + name
}
```

```golang
package mypackage

func helper() int { // func helper is unused
    return 42
}

func anotherHelper() int { // func anotherHelper is unused
    return helper()
}
```

```golang
package main

import "fmt"

func main() {
    fmt.Println("running")
}

func unusedComputation(x, y int) int { // func unusedComputation is unused
    return x*x + y*y
}
```

### Valid

```golang
package main

func main() {
    greet("World")
}

func greet(name string) string {
    return "Hello, " + name
}
```

```golang
package mypackage

// Exported functions are considered used by the package
func PublicAPI() string {
    return internalHelper()
}

func internalHelper() string {
    return "data"
}
```

```golang
package main

func init() {
    // init functions are always considered used
    setup()
}

func main() {}

func setup() {
    // used by init
}
```

```golang
package main

/*
#include <stdio.h>
*/
import "C"

//export CGoFunction
func CGoFunction() {
    // Functions exported to cgo are considered used
}

func main() {}
```
