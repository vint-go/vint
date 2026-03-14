---
title: noSmtpInjectionTaint
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/security/noSmtpInjectionTaint`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/security/noSmtpInjectionTaint:
    # rule options here
```

## Details

Detects SMTP command and header injection vulnerabilities via taint analysis.

This rule uses taint analysis to trace the flow of untrusted data from sources to SMTP-related sinks such as `net/smtp.SendMail` or SMTP client methods. It performs interprocedural data flow analysis to detect cases where user-controlled input is used in email headers or SMTP commands without proper sanitization.

SMTP injection occurs when an attacker injects SMTP protocol commands or additional email headers through unsanitized input. By injecting carriage return and line feed (CRLF) sequences, an attacker can add arbitrary headers (such as BCC recipients to exfiltrate emails) or inject additional SMTP commands. This can lead to spam relay abuse, phishing, or information disclosure.

Email headers and SMTP command parameters derived from user input must be sanitized to remove CRLF sequences and validated against expected formats.

Source: https://github.com/securego/gosec

## Examples

### Invalid

```golang
import "net/smtp"

func sendEmail(w http.ResponseWriter, r *http.Request) {
    to := r.FormValue("to")
    subject := r.FormValue("subject")
    body := r.FormValue("body")

    // User input directly used in email headers
    msg := []byte("To: " + to + "\r\n" +
        "Subject: " + subject + "\r\n" +
        "\r\n" + body)

    smtp.SendMail("smtp.example.com:587", auth, "from@example.com", []string{to}, msg)
}
```

### Valid

```golang
import (
    "net/smtp"
    "mime"
    "strings"
    "regexp"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func sendEmail(w http.ResponseWriter, r *http.Request) {
    to := r.FormValue("to")
    subject := r.FormValue("subject")
    body := r.FormValue("body")

    // Validate email address format
    if !emailRegex.MatchString(to) {
        http.Error(w, "invalid email", http.StatusBadRequest)
        return
    }

    // Sanitize subject - remove CRLF
    cleanSubject := strings.NewReplacer("\r", "", "\n", "").Replace(subject)
    encodedSubject := mime.QEncoding.Encode("utf-8", cleanSubject)

    msg := []byte("To: " + to + "\r\n" +
        "Subject: " + encodedSubject + "\r\n" +
        "MIME-Version: 1.0\r\n" +
        "Content-Type: text/plain; charset=utf-8\r\n" +
        "\r\n" + body)

    smtp.SendMail("smtp.example.com:587", auth, "from@example.com", []string{to}, msg)
}
```
