---
title: useFmtPrint
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useFmtPrint`
- This rule is **not recommended**, meaning it must be explicitly enabled.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useFmtPrint:
    # no configuration options
```

## Details

This rule proposes to replace calls to built-in `print` and `println` with their equivalents from the `fmt` standard package.

`print` and `println` built-in functions are not recommended for use-cases other than
[language bootstrapping and are not guaranteed to stay in the language](https://go.dev/ref/spec#Bootstrapping).

The rule is smart enough to skip calls to user-defined `print` or `println` functions that shadow the built-ins, as well as method calls on types that define `print` or `println` methods.

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
func main() {
    println("just testing", something)
    print("just testing", some, thing+1)
}
```

### Valid

```golang
func main() {
    fmt.Fprintln(os.Stderr, "just testing", something)
    fmt.Fprint(os.Stderr, "just testing", some, thing+1)
}
```

```golang
// User-defined print/println are allowed
func print()   {}
func println() {}

func main() {
    println("just testing")
    print("just testing")
}
```

```golang
// Method calls are allowed
type T struct{}
func (T) print(s string)   {}
func (T) println(s string) {}

func main() {
    t := T{}
    t.print("just testing")
    t.println("just testing")
}
```
