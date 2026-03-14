---
title: noCommandInjectionTaint
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/security/noCommandInjectionTaint`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/security/noCommandInjectionTaint:
    # rule options here
```

## Details

Detects command injection vulnerabilities via taint analysis.

This rule uses taint analysis to trace the flow of untrusted data from sources (such as HTTP request parameters or user input) to command execution sinks (`os/exec.Command`, `syscall.Exec`, etc.). It performs interprocedural analysis to detect cases where user-controlled data flows through multiple functions before reaching a command execution call.

Command injection allows attackers to execute arbitrary system commands on the server, potentially gaining full control of the host. Even partial control over command arguments can be exploited to read files, exfiltrate data, or escalate privileges.

All command arguments derived from user input should be strictly validated against an allowlist of expected values, and shell execution with user-controlled strings should be avoided entirely.

Source: https://github.com/securego/gosec

## Examples

### Invalid

```golang
func handler(w http.ResponseWriter, r *http.Request) {
    filename := r.URL.Query().Get("file")
    // Taint flows from HTTP request to command execution
    output := runCommand(filename)
    w.Write(output)
}

func runCommand(file string) []byte {
    out, _ := exec.Command("cat", file).Output()
    return out
}
```

```golang
func processInput(input string) error {
    // User input passed to shell
    cmd := exec.Command("sh", "-c", input)
    return cmd.Run()
}
```

### Valid

```golang
func handler(w http.ResponseWriter, r *http.Request) {
    filename := r.URL.Query().Get("file")
    // Validate against allowlist
    if !isAllowedFile(filename) {
        http.Error(w, "not allowed", http.StatusForbidden)
        return
    }
    // Read file directly instead of using command
    data, err := os.ReadFile(filepath.Join("/safe/dir", filepath.Base(filename)))
    if err != nil {
        http.Error(w, "error", http.StatusInternalServerError)
        return
    }
    w.Write(data)
}
```

```golang
func processFile(filename string) error {
    // Use constant command with validated argument
    clean := filepath.Base(filename)
    if !regexp.MustCompile(`^[a-zA-Z0-9_.-]+$`).MatchString(clean) {
        return fmt.Errorf("invalid filename")
    }
    return exec.Command("/usr/bin/process", clean).Run()
}
```
