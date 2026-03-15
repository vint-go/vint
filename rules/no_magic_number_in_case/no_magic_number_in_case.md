---
title: noMagicNumberInCase
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noMagicNumberInCase`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noMagicNumberInCase:
    ignored-numbers: "0,0.0,1,1.0"
    ignored-files: "_test.go"
```

## Details

Detects magic numbers used in `switch` statement `case` clauses. A magic number is a numeric literal that is not defined as a constant, but which may change, and therefore can be hard to update. Using magic numbers directly in case clauses makes code less readable and harder to maintain, as the meaning of the case value is not clear without additional context.

This check examines both direct numeric literals in case values and numeric literals within binary expressions used in case conditions (e.g., comparison expressions in expressionless switch statements).

By default, the numbers `0`, `0.0`, `1`, and `1.0` are excluded from detection. Test files (`_test.go`) are also excluded by default.

Source: https://github.com/tommy-muehle/go-mnd

## Examples

### Invalid

```golang
package example

func categorize(x interface{}) {
	switch x {
	case "test":
	case 3: // Magic number: 3, in <case> detected
	}
}
```

```golang
package example

import (
	"fmt"
	"time"
)

func greet() {
	t := time.Now()
	switch {
	case t.Hour() < 12: // Magic number: 12, in <case> detected
		fmt.Println("Good morning!")
	case 17 > t.Hour(): // Magic number: 17, in <case> detected
		fmt.Println("Good afternoon.")
	default:
		fmt.Println("Good evening.")
	}
}
```

```golang
package example

func categorizeFloat(x interface{}) {
	switch x {
	case 1.0:
	case 0.0:
	case 3.0: // Magic number: 3.0, in <case> detected
	}
}
```

### Valid

```golang
package example

const categoryThreshold = 3

func categorize(x interface{}) {
	switch x {
	case "test":
	case categoryThreshold:
	}
}
```

```golang
package example

import (
	"fmt"
	"time"
)

const (
	morningEnd   = 12
	afternoonEnd = 17
)

func greet() {
	t := time.Now()
	switch {
	case t.Hour() < morningEnd:
		fmt.Println("Good morning!")
	case afternoonEnd > t.Hour():
		fmt.Println("Good afternoon.")
	default:
		fmt.Println("Good evening.")
	}
}
```

```golang
package example

func categorizeDefault(x interface{}) {
	switch x {
	case 1.0: // 1.0 is excluded by default
	case 0.0: // 0.0 is excluded by default
	}
}
```
