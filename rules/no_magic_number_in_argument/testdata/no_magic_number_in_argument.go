package fixtures

import (
	"context"
	"math"
	"net/http"
	"os"
	"sync"
	"time"
)

// Invalid: Magic number 9.5 as argument to math.Abs (math.* is ignored by default,
// but we test with a non-default function to trigger).
// Note: math.* is in the default ignored-functions list, so we use a custom function.

func calculateArea() float64 {
	return customAbs(9.5) // MATCH /magic number: 9.5, in <argument> detected/
}

func customAbs(f float64) float64 { return f }

// Invalid: Magic number 200 as argument to http.StatusText.
// Note: http.StatusText is in the default ignored-functions list, so we use a different function.
func getStatusCode() int {
	return lookupCode(200) // MATCH /magic number: 200, in <argument> detected/
}

func lookupCode(code int) int { return code }

// Invalid: Magic number 5 in a binary expression used as argument.
func doWork() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // MATCH /magic number: 5, in <argument> detected/
	defer cancel()
	_ = ctx
}

// Invalid: Magic number 3 as second argument.
func process() {
	foobar(0, 3) // MATCH /magic number: 3, in <argument> detected/
}

func foobar(a, b int) {}

// Invalid: Magic number 500 as argument to make.
func createPool() {
	ch := make(chan int, 500) // MATCH /magic number: 500, in <argument> detected/
	_ = ch
}

// Valid: Using a named constant as argument.
const threshold = 9.5

func calculateAreaGood() float64 {
	return customAbs(threshold)
}

// Valid: Using a named constant as argument.
const statusOK = 200

func getStatusTextGood() string {
	return http.StatusText(statusOK)
}

// Valid: 1 is excluded by default.
func shutdown() {
	os.Exit(1)
}

// Valid: 0 and 1 are excluded by default.
func processGood() {
	foobar(0, 1)
}

// Valid: time.Date is ignored by default.
func createDate() time.Time {
	return time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC)
}

// Valid: math.Abs is ignored by default (math.* pattern).
func mathUsage() float64 {
	return math.Abs(42.5)
}

// Ensure imports are used.
var _ sync.Mutex
