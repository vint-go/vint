---
title: noIdenticalIfElseIfBranches
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noIdenticalIfElseIfBranches`
- This rule is not recommended, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noIdenticalIfElseIfBranches:
    # no configuration options
```

## Details

An `if ... else if` chain with identical branches makes maintenance harder and might be a source of bugs. Duplicated branches should be consolidated into one.

When an `if...else if` chain contains two or more branches with the same body, it suggests a copy-paste mistake or unfinished logic. The rule detects structurally identical branch bodies using AST hashing and reports pairs of matching line numbers.

If any condition in the chain contains a function call (considered "complex"), the confidence of the reported failure is reduced to 0.8 to account for possible intentional side effects.

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
// Two branches with identical bodies
if a {
    foo()
} else if b {
    bar()
} else if c {
    foo() // identical to the first branch
}
```

```golang
// Last else branch identical to first if branch
if condition1 {
    print("something")
} else if condition2 {
    print("other")
} else {
    print("something") // identical to the first branch
}
```

### Valid

```golang
// Each branch has a unique body
if a {
    foo()
} else if b {
    bar()
} else {
    baz()
}
```

```golang
// Branches with init statements are not compared (avoids false positives)
if err := something(); err != nil {
    println(err)
} else if err := somethingElse(); err != nil {
    println(err)
}
```
