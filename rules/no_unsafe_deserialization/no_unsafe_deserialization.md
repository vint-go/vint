---
title: noUnsafeDeserialization
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/security/noUnsafeDeserialization`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/security/noUnsafeDeserialization:
    # rule options here
```

## Details

Detects unsafe deserialization of untrusted data via taint analysis.

This rule uses taint analysis to trace the flow of untrusted data from sources to deserialization sinks. It performs interprocedural data flow analysis to detect cases where user-controlled input is deserialized without proper validation, potentially allowing attackers to craft malicious payloads.

Unsafe deserialization can lead to remote code execution, denial of service, or privilege escalation depending on the deserialization library and the types being deserialized. In Go, while the risk is generally lower than in languages like Java or Python, deserializing untrusted data with `encoding/gob`, `encoding/xml`, or third-party serialization libraries can still lead to resource exhaustion, panic, or logic bugs through crafted input.

Untrusted data should be validated before deserialization. Use strict type constraints, size limits, and input validation. Prefer safe formats like JSON with explicit struct types over flexible formats like `encoding/gob`.

Source: https://github.com/securego/gosec

## Examples

### Invalid

```golang
import "encoding/gob"

func handler(w http.ResponseWriter, r *http.Request) {
    // Deserializing untrusted data from HTTP request body
    decoder := gob.NewDecoder(r.Body)
    var data interface{}
    decoder.Decode(&data) // Unsafe - arbitrary type deserialization
}
```

```golang
import "encoding/xml"

func handler(w http.ResponseWriter, r *http.Request) {
    body, _ := io.ReadAll(r.Body)
    var result interface{}
    // XML deserialization of untrusted input (also vulnerable to XXE)
    xml.Unmarshal(body, &result)
}
```

### Valid

```golang
import "encoding/json"

type UserInput struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

func handler(w http.ResponseWriter, r *http.Request) {
    // Limit body size
    r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

    // Deserialize into strongly-typed struct
    var input UserInput
    decoder := json.NewDecoder(r.Body)
    decoder.DisallowUnknownFields()
    if err := decoder.Decode(&input); err != nil {
        http.Error(w, "invalid input", http.StatusBadRequest)
        return
    }

    // Validate deserialized data
    if input.Name == "" || input.Email == "" {
        http.Error(w, "missing fields", http.StatusBadRequest)
        return
    }
}
```

```golang
import "encoding/json"

func handler(w http.ResponseWriter, r *http.Request) {
    // Limit request body size
    r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

    // Use concrete type instead of interface{}
    var config AppConfig
    if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
        http.Error(w, "invalid JSON", http.StatusBadRequest)
        return
    }
}
```
