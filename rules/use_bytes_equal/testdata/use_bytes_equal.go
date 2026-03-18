package fixtures

import "bytes"

func invalid(a, b []byte) {
	if bytes.Compare(a, b) == 0 { // MATCH /replace bytes.Compare with bytes.Equal/
		// should use bytes.Equal
	}

	if 0 == bytes.Compare(a, b) { // MATCH /replace bytes.Compare with bytes.Equal/
		// reversed comparison
	}

	if bytes.Compare(a, b) != 0 { // MATCH /replace bytes.Compare with bytes.Equal/
		// inequality check
	}

	if 0 != bytes.Compare(a, b) { // MATCH /replace bytes.Compare with bytes.Equal/
		// reversed inequality
	}
}

func valid(a, b []byte) {
	if bytes.Equal(a, b) {
		// correct usage
	}

	if bytes.Compare(a, b) > 0 {
		// not an equality check, this is fine
	}

	if bytes.Compare(a, b) < 0 {
		// not an equality check
	}

	if bytes.Compare(a, b) == 1 {
		// comparing against non-zero, this is fine
	}
}
