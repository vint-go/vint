package fixtures

import (
	"crypto/tls"
	"net"
)

func bindAllInterfacesExplicit() {
	listener, err := net.Listen("tcp", "0.0.0.0:8080") // MATCH /binding to all interfaces is a security risk, bind to a specific interface instead/
	_ = listener
	_ = err
}

func bindAllInterfacesShorthand() {
	listener, err := net.Listen("tcp", ":8080") // MATCH /binding to all interfaces is a security risk, bind to a specific interface instead/
	_ = listener
	_ = err
}

func bindAllInterfacesTLS(tlsConfig *tls.Config) {
	listener, err := tls.Listen("tcp", "0.0.0.0:443", tlsConfig) // MATCH /binding to all interfaces is a security risk, bind to a specific interface instead/
	_ = listener
	_ = err
}

func bindLocalhostOnly() {
	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	_ = listener
	_ = err
}

func bindSpecificInterface() {
	listener, err := net.Listen("tcp", "192.168.1.10:8080")
	_ = listener
	_ = err
}

func bindLocalhostTLS(tlsConfig *tls.Config) {
	listener, err := tls.Listen("tcp", "127.0.0.1:443", tlsConfig)
	_ = listener
	_ = err
}
