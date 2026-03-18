package fixtures

import "time"

func manualConversions() {
	t := time.Now()

	_ = t.UnixNano() / 1000000 // MATCH /use UnixMilli() instead of manual time conversion/

	_ = t.UnixNano() / 1000 // MATCH /use UnixMicro() instead of manual time conversion/
}

func validUsages() {
	t := time.Now()

	_ = t.UnixMilli()

	_ = t.UnixMicro()

	_ = t.UnixNano()

	// Dividing by other numbers is fine
	_ = t.UnixNano() / 100

	_ = t.UnixNano() / 1000000000
}
