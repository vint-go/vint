---
title: noRangeVariableAlias
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noRangeVariableAlias`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noRangeVariableAlias:
    # rule options here
```

## Details

Detects implicit memory aliasing in `for...range` statements (applicable to Go 1.21 or lower).

In Go versions 1.21 and earlier, the loop variable in a `for...range` statement is reused across iterations. Taking the address of the loop variable (`&v`) inside the loop body creates a pointer that refers to the same memory location in every iteration, causing all pointers to reference the last element after the loop completes. This is a common source of bugs in Go programs.

Starting from Go 1.22, each iteration of a `for...range` loop creates a new variable, eliminating this issue. However, for code that must support Go 1.21 or earlier, explicitly copying the loop variable before taking its address is necessary.

Source: https://github.com/securego/gosec

## Examples

### Invalid

```golang
// Go 1.21 or lower
type Item struct {
    Name string
}

func process(items []Item) []*Item {
    var result []*Item
    for _, item := range items {
        // All pointers will reference the last item
        result = append(result, &item)
    }
    return result
}
```

```golang
// Go 1.21 or lower
func processMap(m map[string]int) []*int {
    var ptrs []*int
    for _, v := range m {
        ptrs = append(ptrs, &v) // All point to same memory
    }
    return ptrs
}
```

### Valid

```golang
// Go 1.21 or lower - explicit copy
type Item struct {
    Name string
}

func process(items []Item) []*Item {
    var result []*Item
    for _, item := range items {
        item := item // Create a new variable
        result = append(result, &item)
    }
    return result
}
```

```golang
// Go 1.21 or lower - using index
func process(items []Item) []*Item {
    var result []*Item
    for i := range items {
        result = append(result, &items[i])
    }
    return result
}
```

```golang
// Go 1.22+ - this is safe, each iteration has its own variable
func process(items []Item) []*Item {
    var result []*Item
    for _, item := range items {
        result = append(result, &item) // Safe in Go 1.22+
    }
    return result
}
```
