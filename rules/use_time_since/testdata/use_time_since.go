package fixtures

import "time"

func badUsage(start time.Time) time.Duration {
	return time.Now().Sub(start) // MATCH /replace time.Now().Sub(x) with time.Since(x)/
}

func goodUsage(start time.Time) time.Duration {
	return time.Since(start)
}

func alsoValid() {
	t := time.Now()
	// Sub on a variable, not time.Now()
	_ = t.Sub(time.Now())
}
