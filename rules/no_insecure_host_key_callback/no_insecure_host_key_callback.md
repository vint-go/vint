---
title: noInsecureHostKeyCallback
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/security/noInsecureHostKeyCallback`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/security/noInsecureHostKeyCallback:
    # rule options here
```

## Details

Audits the use of `ssh.InsecureIgnoreHostKey` function from the `golang.org/x/crypto/ssh` package.

This rule flags calls to `ssh.InsecureIgnoreHostKey()`, which returns a function that accepts any host key. Using this function disables SSH host key verification, making the connection vulnerable to man-in-the-middle (MITM) attacks. An attacker can intercept the SSH connection and impersonate the remote server without detection.

SSH host key verification is a critical security mechanism that ensures the client is connecting to the intended server. In production environments, applications should use proper host key verification by maintaining a known hosts list or implementing a custom `HostKeyCallback` that validates the server's identity.

Source: https://github.com/securego/gosec

## Examples

### Invalid

```golang
import "golang.org/x/crypto/ssh"

// Using InsecureIgnoreHostKey disables host key verification
config := &ssh.ClientConfig{
    User: "admin",
    Auth: []ssh.AuthMethod{
        ssh.Password("password"),
    },
    HostKeyCallback: ssh.InsecureIgnoreHostKey(),
}
client, err := ssh.Dial("tcp", "example.com:22", config)
```

### Valid

```golang
import "golang.org/x/crypto/ssh"
import "golang.org/x/crypto/ssh/knownhosts"

// Using known hosts file for host key verification
hostKeyCallback, err := knownhosts.New("/home/user/.ssh/known_hosts")
if err != nil {
    log.Fatal(err)
}

config := &ssh.ClientConfig{
    User: "admin",
    Auth: []ssh.AuthMethod{
        ssh.Password("password"),
    },
    HostKeyCallback: hostKeyCallback,
}
client, err := ssh.Dial("tcp", "example.com:22", config)
```

```golang
import "golang.org/x/crypto/ssh"

// Using a fixed host key for verification
config := &ssh.ClientConfig{
    User: "admin",
    Auth: []ssh.AuthMethod{
        ssh.Password("password"),
    },
    HostKeyCallback: ssh.FixedHostKey(expectedKey),
}
```
