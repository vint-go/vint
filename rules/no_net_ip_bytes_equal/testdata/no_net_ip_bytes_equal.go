package fixtures

import (
	"bytes"
	"net"
)

func badBytesEqualNetIP() {
	ip1 := net.ParseIP("127.0.0.1")
	ip2 := net.IPv4(127, 0, 0, 1)
	// Wrong: different byte representations may not match
	if bytes.Equal(ip1, ip2) { // MATCH /use net.IP.Equal to compare net.IP values, not bytes.Equal/
		// ...
	}
}

func goodNetIPEqual() {
	ip1 := net.ParseIP("127.0.0.1")
	ip2 := net.IPv4(127, 0, 0, 1)
	// Correct: handles different representations
	if ip1.Equal(ip2) {
		// ...
	}
}

func goodBytesEqualNonIP() {
	a := []byte{1, 2, 3}
	b := []byte{4, 5, 6}
	// Correct: not comparing net.IP values
	if bytes.Equal(a, b) {
		// ...
	}
}
