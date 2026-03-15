---
title: noUnusedVariable
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noUnusedVariable`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noUnusedVariable:
    # Controls whether writing to a variable counts as a use.
    # post-statements-are-reads: true
    # Controls whether local variables are always considered used.
    # local-variables-are-used: true
    # Controls whether variables in generated files are considered used.
    # generated-is-used: true
```

## Details

Detects package-level variables that are declared but never read or referenced in the codebase.

The `unused` linter traces variable usage through a dependency graph. A variable is considered "used" if it is read somewhere in the code that is itself reachable from an entry point. Key considerations include:

- Exported package-level variables are always considered used (since external packages may reference them).
- Writing to a variable does not count as using it by default, unless the `post-statements-are-reads` option is enabled.
- In test files, assigning to a package-level variable does count as a use (rule 4.9).
- Variables use their types, so even if a variable is unused, its type may still be marked as used through other paths.
- The blank identifier (`_`) is always considered used.
- Variables in generated files are considered used when the `generated-is-used` option is enabled (default).
- When the `local-variables-are-used` option is enabled (default), local variables are not reported. This check primarily targets package-level variables unless that option is disabled.

Note: Go's compiler already catches most unused local variables. This rule is particularly useful for detecting unused package-level variables that the compiler does not flag.

Source: https://github.com/dominikh/go-tools/tree/master/unused

## Examples

### Invalid

```golang
package mypackage

var unusedConfig = map[string]string{ // var unusedConfig is unused
    "key": "value",
}

func DoWork() string {
    return "done"
}
```

```golang
package mypackage

var counter int // var counter is unused

func Increment() {
    // Writing to counter does not count as using it
    // (when field-writes-are-uses is false)
    counter++
}
```

```golang
package mypackage

var (
    Active   = true
    logLevel = "debug" // var logLevel is unused
)

func Status() bool {
    return Active
}
```

### Valid

```golang
package mypackage

var config = map[string]string{
    "key": "value",
}

func GetConfig() map[string]string {
    return config // config is read here, so it is used
}
```

```golang
package mypackage

// Exported variables are always considered used
var Version = "1.0.0"
```

```golang
package mypackage

var _ Interface = (*Impl)(nil) // blank identifier is always used

type Interface interface {
    Method()
}

type Impl struct{}

func (i *Impl) Method() {}
```

```golang
package mypackage

var cache map[string]string

func init() {
    cache = make(map[string]string)
}

func Get(key string) string {
    return cache[key]
}
```
