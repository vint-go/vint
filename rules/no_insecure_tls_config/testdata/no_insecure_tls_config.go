package fixtures

import "crypto/tls"

// Invalid: InsecureSkipVerify set to true
func insecureSkipVerify() {
	_ = &tls.Config{
		InsecureSkipVerify: true, // MATCH /TLS InsecureSkipVerify set to true disables certificate verification/
	}
}

// Invalid: MinVersion set to TLS 1.0
func lowMinVersionTLS10() {
	_ = &tls.Config{
		MinVersion: tls.VersionTLS10, // MATCH /TLS MinVersion too low, should be at least tls.VersionTLS12/
	}
}

// Invalid: MinVersion set to TLS 1.1
func lowMinVersionTLS11() {
	_ = &tls.Config{
		MinVersion: tls.VersionTLS11, // MATCH /TLS MinVersion too low, should be at least tls.VersionTLS12/
	}
}

// Invalid: MaxVersion capped to TLS 1.0
func lowMaxVersionTLS10() {
	_ = &tls.Config{
		MaxVersion: tls.VersionTLS10, // MATCH /TLS MaxVersion too low, should be at least tls.VersionTLS12/
	}
}

// Invalid: MaxVersion capped to TLS 1.1
func lowMaxVersionTLS11() {
	_ = &tls.Config{
		MaxVersion: tls.VersionTLS11, // MATCH /TLS MaxVersion too low, should be at least tls.VersionTLS12/
	}
}

// Invalid: Using weak cipher suites
func weakCiphers() {
	_ = &tls.Config{
		CipherSuites: []uint16{
			tls.TLS_RSA_WITH_RC4_128_SHA,          // MATCH /use of weak TLS cipher suite tls.TLS_RSA_WITH_RC4_128_SHA/
			tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,     // MATCH /use of weak TLS cipher suite tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA/
		},
	}
}

// Invalid: Multiple insecure settings combined
func multipleInsecure() {
	_ = &tls.Config{
		InsecureSkipVerify: true,             // MATCH /TLS InsecureSkipVerify set to true disables certificate verification/
		MinVersion:         tls.VersionTLS10, // MATCH /TLS MinVersion too low, should be at least tls.VersionTLS12/
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA, // MATCH /use of weak TLS cipher suite tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA/
		},
	}
}

// Valid: Secure TLS configuration with TLS 1.2
func secureTLS12() {
	_ = &tls.Config{
		MinVersion: tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		},
	}
}

// Valid: Modern TLS 1.3 only configuration
func secureTLS13() {
	_ = &tls.Config{
		MinVersion: tls.VersionTLS13,
	}
}

// Valid: InsecureSkipVerify explicitly set to false
func skipVerifyFalse() {
	_ = &tls.Config{
		InsecureSkipVerify: false,
		MinVersion:         tls.VersionTLS12,
	}
}

// Valid: MaxVersion set to TLS 1.2
func maxVersionTLS12() {
	_ = &tls.Config{
		MinVersion: tls.VersionTLS12,
		MaxVersion: tls.VersionTLS12,
	}
}

// Valid: MaxVersion set to TLS 1.3
func maxVersionTLS13() {
	_ = &tls.Config{
		MinVersion: tls.VersionTLS12,
		MaxVersion: tls.VersionTLS13,
	}
}

// Valid: Empty tls.Config (uses defaults)
func emptyConfig() {
	_ = &tls.Config{}
}
