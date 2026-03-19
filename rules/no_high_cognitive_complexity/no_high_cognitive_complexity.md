---
title: noHighCognitiveComplexity
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/complexity/noHighCognitiveComplexity`
- This rule is not recommended, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/complexity/noHighCognitiveComplexity:
    arguments: [7]
```

## Details

[Cognitive complexity](https://www.sonarsource.com/docs/CognitiveComplexity.pdf) is a measure of how hard code is to understand.
While cyclomatic complexity is good to measure "testability" of the code,
cognitive complexity aims to provide a more precise measure of the difficulty of understanding the code.
Enforcing a maximum complexity per function helps to keep code readable and maintainable.

The rule accepts a single integer argument specifying the maximum allowed cognitive complexity per function.
The default threshold is 7.

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
func sumOfPrimes(max int) int {
	total := 0
OUT:
	for i := 1; i <= max; i++ {           // +1
		for j := 2; j < i; j++ {          // +1 +1(nesting)
			if i%j == 0 {                  // +1 +2(nesting)
				continue OUT               // +1
			}
		}
		total += i
	}
	return total
}
// cognitive complexity = 7
```

```golang
func example(a, b, c bool) {
	if a {                    // +1
		if b {                // +1 +1(nesting)
			if c {            // +1 +2(nesting)
				doSomething()
			}
		}
	}
}
// cognitive complexity = 6
```

### Valid

```golang
func simple(x int) bool {
	if x > 0 {
		return true
	}
	return false
}
// cognitive complexity = 1
```

```golang
func linearFlow(items []string) {
	for _, item := range items {
		process(item)
	}
}
// cognitive complexity = 1
```
