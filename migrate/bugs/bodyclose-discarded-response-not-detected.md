# bodyclose: discarded response not detected

## Affected rule
`lint/correctness/noUnclosedBodies` (maps to golangci-lint `bodyclose`)

## Behavior in golangci-lint
golangci-lint's `bodyclose` linter detects both patterns:
1. Response assigned to a variable but `Body.Close()` never called
2. Response entirely discarded (return values not captured at all)

For example, both of these trigger `bodyclose`:
```go
// Pattern 1: assigned but not closed
resp, err := http.Get("https://example.com")
if err != nil { return }
// missing resp.Body.Close()

// Pattern 2: entirely discarded
http.Get("https://example.com") //nolint:bodyclose
```

## Behavior in vint
Vint's `noUnclosedBodies` rule only detects **Pattern 1** — it walks AST assignment statements looking for variables typed `*http.Response`, then checks whether `Body.Close()` is called on them.

It does **not** detect **Pattern 2** — when the return value of `http.Get()` (or similar functions returning `*http.Response`) is entirely discarded without assignment.

## Gap
Functions returning `*http.Response` whose return values are completely ignored are not flagged by `noUnclosedBodies`. The rule's AST walker only inspects assignment targets, so a bare expression statement like `http.Get(url)` is invisible to it.

## Example
```go
package example

import "net/http"

func leakyRequest() {
	// This is caught by golangci-lint's bodyclose but NOT by vint's noUnclosedBodies
	http.Get("https://example.com")
}

func leakyButAssigned() {
	// This IS caught by both golangci-lint's bodyclose and vint's noUnclosedBodies
	resp, _ := http.Get("https://example.com")
	_ = resp
}
```

## Impact on migration
When migrating `//nolint:bodyclose` directives, the smart nolint conversion (which runs vint rules on stripped source to determine which rules fire) will not generate `// vint-ignore` lines for Pattern 2, because `noUnclosedBodies` doesn't fire there. This means:

- The migration is still **safe** — no false suppressions are generated
- But the original `//nolint:bodyclose` directive is removed without a replacement, which is correct since vint wouldn't flag that line anyway
- If `noUnclosedBodies` is later enhanced to catch discarded responses, those lines would start producing warnings without suppressions
