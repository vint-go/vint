---
title: noMisspelledWords
description: Detects and corrects commonly misspelled English words in Go source files
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noMisspelledWords`
- This rule is **recommended**, meaning it is enabled by default.
- This rule has a **fix**.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noMisspelledWords:
    locale: "" # Set to "US" or "UK" to enforce locale-specific spelling (e.g., "colour" vs "color")
    mode: "" # Set to "restricted" to check only comments, leave empty to check all text
    extra-words: [] # List of custom misspelling corrections: [{typo: "misspeling", correction: "misspelling"}]
    ignore-rules: [] # List of correction rules to ignore (e.g., ["importas"])
```

## Details

The `noMisspelledWords` rule detects commonly misspelled English words in Go source files and provides automatic corrections. It works by maintaining a large dictionary of known misspellings (thousands of entries) and scanning source code to find and flag them.

The linter supports two scanning modes:
- **Default mode**: Scans the entire file content including identifiers, string literals, comments, and other text using the `Replace` function.
- **Restricted mode** (`mode: "restricted"`): Only scans Go comments using the `ReplaceGo` function, leaving identifiers and string literals unchecked.

### Locale Support

The linter can enforce locale-specific spelling conventions:
- **No locale** (default): Only flags clearly misspelled words without enforcing US or UK spelling conventions.
- **US locale** (`locale: "US"`): Enforces American English spelling. Words like "colour" will be flagged and corrected to "color", "organisation" to "organization", etc.
- **UK locale** (`locale: "UK"` or `locale: "GB"`): Enforces British English spelling. Words like "color" will be flagged and corrected to "colour", "organization" to "organisation", etc.

### Custom Words

You can define additional misspelling corrections via the `extra-words` configuration. Each entry must have a `typo` (the misspelled form) and a `correction` (the correct form). Both values must contain only letters and are normalized to lowercase internally.

### Ignoring Rules

Specific correction rules can be suppressed using the `ignore-rules` option. This is useful when a correction conflicts with domain-specific terminology or identifiers in your codebase.

### Detection Process

The misspell linter uses a multi-stage detection approach:
1. A fast bulk replacement engine identifies potential misspellings.
2. Modified lines are rechecked for accuracy.
3. Individual words are extracted and validated against the known dictionary.
4. CamelCase words are excluded to avoid false positives on identifiers.
5. Only confirmed misspellings that match the predefined dictionary are reported.

The linter reports diagnostics in the format: `filename:line:column: "original" is a misspelling of "correction"`, and provides automatic fix suggestions.

Source: https://github.com/golangci/misspell

## Examples

### Invalid

```golang
// This function calcualtes the total amount
func Calculate(items []Item) int {
	total := 0
	for _, item := range items {
		total += item.Price
	}
	return total
}
```

```golang
// Recieve data from the channel and proccess it
func HandleData(ch chan Data) {
	data := <-ch
	fmt.Println(data)
}
```

```golang
package main

// The seperate function splits a string into its componets
func Separate(input string) []string {
	return strings.Split(input, ",")
}
```

```golang
// This is a convienent utilty function for occurance counting
func CountOccurrences(s string, sub string) int {
	return strings.Count(s, sub)
}
```

```golang
// Definately check if the enviroment variable is set
func CheckEnv(key string) bool {
	val, ok := os.LookupEnv(key)
	return ok && val != ""
}
```

### Valid

```golang
// This function calculates the total amount
func Calculate(items []Item) int {
	total := 0
	for _, item := range items {
		total += item.Price
	}
	return total
}
```

```golang
// Receive data from the channel and process it
func HandleData(ch chan Data) {
	data := <-ch
	fmt.Println(data)
}
```

```golang
package main

// The separate function splits a string into its components
func Separate(input string) []string {
	return strings.Split(input, ",")
}
```

```golang
// This is a convenient utility function for occurrence counting
func CountOccurrences(s string, sub string) int {
	return strings.Count(s, sub)
}
```

```golang
// Definitely check if the environment variable is set
func CheckEnv(key string) bool {
	val, ok := os.LookupEnv(key)
	return ok && val != ""
}
```
