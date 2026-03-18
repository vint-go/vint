package fixtures

import "time"

func badUsage(deadline time.Time) time.Duration {
	return deadline.Sub(time.Now()) // MATCH /replace x.Sub(time.Now()) with time.Until(x)/
}

func goodUsage(deadline time.Time) time.Duration {
	return time.Until(deadline)
}

func alsoValid() {
	now := time.Now()
	// time.Now().Sub(x) is a different pattern (useTimeSince)
	_ = time.Now().Sub(now)
}
