---
title: noSecretInSerialization
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/security/noSecretInSerialization`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/security/noSecretInSerialization:
    # rule options here
```

## Details

Detects potential exposure of secrets via JSON, YAML, XML, or TOML marshaling.

This rule identifies struct fields that appear to contain sensitive data (based on field names matching patterns like `password`, `secret`, `token`, `apiKey`, etc.) and are not excluded from serialization output. When structs with secret fields are marshaled to JSON, YAML, XML, or TOML without proper struct tags to omit sensitive fields, secrets may be inadvertently included in API responses, log output, or persisted data.

Sensitive fields should be annotated with appropriate struct tags to exclude them from serialization (e.g., `json:"-"`, `yaml:"-"`, `xml:"-"`) or moved to separate structs that are not serialized.

Source: https://github.com/securego/gosec

## Examples

### Invalid

```golang
type User struct {
    Name     string `json:"name"`
    Email    string `json:"email"`
    Password string `json:"password"` // Secret exposed in JSON output
}

func handler(w http.ResponseWriter, r *http.Request) {
    user := getUser()
    json.NewEncoder(w).Encode(user) // Password included in response
}
```

```golang
type Config struct {
    Host     string `yaml:"host"`
    Port     int    `yaml:"port"`
    APIKey   string `yaml:"api_key"` // Secret exposed in YAML output
}

func saveConfig(cfg Config) error {
    data, _ := yaml.Marshal(cfg)
    return os.WriteFile("config.yaml", data, 0644)
}
```

### Valid

```golang
type User struct {
    Name     string `json:"name"`
    Email    string `json:"email"`
    Password string `json:"-"` // Excluded from JSON output
}

func handler(w http.ResponseWriter, r *http.Request) {
    user := getUser()
    json.NewEncoder(w).Encode(user) // Password not included
}
```

```golang
// Using separate response type without sensitive fields
type UserResponse struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

func handler(w http.ResponseWriter, r *http.Request) {
    user := getUser()
    resp := UserResponse{Name: user.Name, Email: user.Email}
    json.NewEncoder(w).Encode(resp)
}
```
