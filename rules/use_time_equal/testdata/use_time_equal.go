package fixtures

import "time"

func compareTimeValues() bool {
	t := time.Now()
	u := t

	if !t.After(u) {
		return t == u // MATCH /use t.Equal(u) instead of "==" operator/
	}

	return t != u // MATCH /use !t.Equal(u) instead of "!=" operator/
}

func isNow(t time.Time) bool    { return t == time.Now() }    // MATCH /use t.Equal(time.Now()) instead of "==" operator/
func isNotNow(t time.Time) bool { return time.Now() != t }    // MATCH /use !time.Now().Equal(t) instead of "!=" operator/

// Valid: using Equal method
func isSame(a, b time.Time) bool {
	return a.Equal(b)
}

// Valid: comparing non-time values
func compareInts(a, b int) bool {
	return a == b
}
