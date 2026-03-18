package fixtures

import "time"

func suspiciousTimeSleep() {
	time.Sleep(1) // MATCH /suspiciously small untyped constant 1 in time.Sleep/
	time.Sleep(0) // MATCH /suspiciously small untyped constant 0 in time.Sleep/
	time.Sleep(100) // MATCH /suspiciously small untyped constant 100 in time.Sleep/
	time.Sleep(1000) // MATCH /suspiciously small untyped constant 1000 in time.Sleep/
}

func validTimeSleep() {
	time.Sleep(1 * time.Second)
	time.Sleep(500 * time.Millisecond)
	time.Sleep(time.Duration(5) * time.Second)
	time.Sleep(2 * time.Minute)

	d := 5 * time.Second
	time.Sleep(d)
}
