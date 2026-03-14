---
title: noInsecureCookie
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/security/noInsecureCookie`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/security/noInsecureCookie:
    # rule options here
```

## Details

Detects insecure HTTP cookie configurations missing `Secure`, `HttpOnly`, or `SameSite` attributes.

This rule uses SSA analysis to identify `http.Cookie` instances that are missing important security attributes. Cookies without the `Secure` flag can be transmitted over unencrypted HTTP connections, exposing their contents to network attackers. Cookies without the `HttpOnly` flag are accessible to JavaScript, making them vulnerable to cross-site scripting (XSS) attacks. Cookies without a proper `SameSite` attribute may be sent with cross-site requests, enabling cross-site request forgery (CSRF) attacks.

All session cookies and cookies containing sensitive data should set `Secure: true`, `HttpOnly: true`, and `SameSite: http.SameSiteStrictMode` or `http.SameSiteLaxMode`.

Source: https://github.com/securego/gosec

## Examples

### Invalid

```golang
import "net/http"

func setCookie(w http.ResponseWriter) {
    http.SetCookie(w, &http.Cookie{
        Name:  "session",
        Value: sessionToken,
        // Missing Secure, HttpOnly, and SameSite attributes
    })
}
```

```golang
import "net/http"

func setCookie(w http.ResponseWriter) {
    http.SetCookie(w, &http.Cookie{
        Name:     "session",
        Value:    sessionToken,
        Secure:   false, // Should be true
        HttpOnly: false, // Should be true
    })
}
```

### Valid

```golang
import "net/http"

func setCookie(w http.ResponseWriter) {
    http.SetCookie(w, &http.Cookie{
        Name:     "session",
        Value:    sessionToken,
        Secure:   true,
        HttpOnly: true,
        SameSite: http.SameSiteStrictMode,
        Path:     "/",
    })
}
```
