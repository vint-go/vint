package fixtures

import "ssh"

// Invalid: Using InsecureIgnoreHostKey disables host key verification
func badConfig() {
	config := &ssh.ClientConfig{
		User: "admin",
		Auth: []ssh.AuthMethod{
			ssh.Password("password"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // MATCH /use of ssh.InsecureIgnoreHostKey disables host key verification and is vulnerable to MITM attacks/
	}
	_ = config
}

// Invalid: Using InsecureIgnoreHostKey in a variable assignment
func badAssignment() {
	cb := ssh.InsecureIgnoreHostKey() // MATCH /use of ssh.InsecureIgnoreHostKey disables host key verification and is vulnerable to MITM attacks/
	_ = cb
}

// Valid: Using FixedHostKey for verification
func goodFixedKey() {
	var expectedKey ssh.PublicKey
	config := &ssh.ClientConfig{
		User: "admin",
		Auth: []ssh.AuthMethod{
			ssh.Password("password"),
		},
		HostKeyCallback: ssh.FixedHostKey(expectedKey),
	}
	_ = config
}

// Valid: Using a custom callback
func goodCustomCallback() {
	config := &ssh.ClientConfig{
		User:            "admin",
		HostKeyCallback: customHostKeyCallback,
	}
	_ = config
}

var customHostKeyCallback ssh.HostKeyCallback
