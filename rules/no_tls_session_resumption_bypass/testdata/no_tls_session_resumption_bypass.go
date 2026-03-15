package fixtures

import (
	"crypto/tls"
	"crypto/x509"
)

// Invalid: VerifyPeerCertificate without VerifyConnection or SessionTicketsDisabled
func verifyPeerOnly() {
	_ = &tls.Config{ // MATCH /VerifyPeerCertificate is not called during TLS session resumption; set VerifyConnection or disable session tickets/
		VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
			return nil
		},
	}
}

// Valid: Both VerifyPeerCertificate and VerifyConnection set
func verifyPeerAndConnection() {
	_ = &tls.Config{
		VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
			return nil
		},
		VerifyConnection: func(state tls.ConnectionState) error {
			return nil
		},
	}
}

// Valid: VerifyPeerCertificate with SessionTicketsDisabled
func verifyPeerWithSessionDisabled() {
	_ = &tls.Config{
		SessionTicketsDisabled: true,
		VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
			return nil
		},
	}
}

// Valid: No VerifyPeerCertificate set at all
func noVerifyPeer() {
	_ = &tls.Config{
		MinVersion: tls.VersionTLS12,
	}
}

// Valid: Empty tls.Config
func emptyConfig() {
	_ = &tls.Config{}
}

// Valid: Only VerifyConnection set (no VerifyPeerCertificate)
func onlyVerifyConnection() {
	_ = &tls.Config{
		VerifyConnection: func(state tls.ConnectionState) error {
			return nil
		},
	}
}
