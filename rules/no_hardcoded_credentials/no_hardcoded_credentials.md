---
title: noHardcodedCredentials
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/security/noHardcodedCredentials`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/security/noHardcodedCredentials:
    pattern: "(?i)(passwd|password|pwd|secret|token|api[_-]?key|apikey|access[_-]?key|auth[_-]?token|credentials?)"
    entropyThreshold: 3.5
```

## Details

Detects hardcoded credentials such as passwords, API keys, and tokens embedded directly in Go source code.

Configuration options:

- `pattern`: Custom regex pattern for matching credential variable names. Defaults to matching common credential identifiers.
- `entropyThreshold`: Minimum Shannon entropy for a string to be flagged during entropy analysis. Default is `3.5`.

This rule scans variable assignments, declarations, equality comparisons, and composite literals for credential-like patterns. It matches variable names against common credential identifiers (e.g., `passwd`, `password`, `secret`, `token`, `apiKey`) and checks string values against known secret formats including AWS access keys, Slack tokens, GitHub tokens, and Google API keys.

The rule also performs entropy analysis using the zxcvbn password strength library to identify high-entropy strings that are likely secrets. The entropy threshold, per-character threshold, and minimum entropy length are configurable.

Hardcoded credentials are a significant security risk because they can be extracted from source code, version control history, or compiled binaries. Secrets should be stored in environment variables, secret management systems, or configuration files that are excluded from version control.

Source: https://github.com/securego/gosec

## Examples

### Invalid

```golang
// Hardcoded password in variable assignment
password := "supersecretpassword123"
```

```golang
// Hardcoded API key in constant declaration
const apiKey = "AKIAIOSFODNN7EXAMPLE"
```

```golang
// Hardcoded token in struct literal
config := Config{
    Token: "ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
}
```

```golang
// Password comparison with hardcoded value
if userPassword == "admin123" {
    grantAccess()
}
```

### Valid

```golang
// Reading password from environment variable
password := os.Getenv("DB_PASSWORD")
```

```golang
// Using a secret management library
apiKey, err := secrets.Get("api-key")
if err != nil {
    log.Fatal(err)
}
```

```golang
// Reading token from configuration file
token := config.GetString("auth.token")
```
