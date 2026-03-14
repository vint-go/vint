---
title: noTrojanSourceBidi
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/security/noTrojanSourceBidi`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/security/noTrojanSourceBidi:
    # rule options here
```

## Details

Detects Trojan Source attacks using bidirectional Unicode characters in Go source files.

Trojan Source is an attack that exploits Unicode bidirectional (BiDi) control characters to make source code appear different from what compilers and interpreters execute. Characters like Right-to-Left Override (U+202E), Left-to-Right Override (U+202D), Right-to-Left Embedding (U+202B), and similar BiDi control characters can reorder the display of source code in text editors, allowing malicious logic to be hidden in plain sight.

This attack was documented in CVE-2021-42574 and can be used to introduce vulnerabilities into code reviews where reviewers see different logic than what actually executes. The rule scans source files for the presence of these dangerous Unicode control characters.

Source: https://github.com/securego/gosec

## Examples

### Invalid

```golang
// Code containing bidirectional Unicode control characters
// that could reorder displayed text
// (Note: actual BiDi characters are invisible in plain text)
func checkAccess(isAdmin bool) bool {
    // The following line contains a hidden RLO character
    // that could make "if !isAdmin" appear as "if isAdmin"
    if isAdmin {
        return true
    }
    return false
}
```

### Valid

```golang
// Clean source code without bidirectional Unicode characters
func checkAccess(isAdmin bool) bool {
    if isAdmin {
        return true
    }
    return false
}
```

```golang
// Using only standard ASCII or safe Unicode characters
func greet(name string) string {
    return fmt.Sprintf("Hello, %s!", name)
}
```
