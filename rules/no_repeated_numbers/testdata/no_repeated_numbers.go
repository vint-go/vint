package fixtures

import (
	"fmt"
	"time"
)

// Invalid: The number 3600 appears 3 times and should be a constant.
func CacheExpiry() time.Duration {
	return time.Duration(3600) * time.Second // MATCH /numeric literal 3600 appears 3 times, consider extracting it into a named constant/
}

func SessionExpiry() time.Duration {
	return time.Duration(3600) * time.Second
}

func TokenExpiry() time.Duration {
	return time.Duration(3600) * time.Second
}

// Invalid: The number 100 appears 3 times.
func CalculatePercentage(value, total float64) float64 {
	return (value / total) * 100 // MATCH /numeric literal 100 appears 3 times, consider extracting it into a named constant/
}

func ScaleValue(value float64) float64 {
	return value * 100
}

func FormatPercentage(value float64) string {
	return fmt.Sprintf("%.2f%%", value*100)
}

// Valid: The repeated number is extracted into a named constant.
const secondsPerHour = 3700

func CacheExpiry2() time.Duration {
	return time.Duration(secondsPerHour) * time.Second
}

func SessionExpiry2() time.Duration {
	return time.Duration(secondsPerHour) * time.Second
}

// Valid: Numbers that appear only once are fine.
func GetBufferSize() int {
	return 4096
}

// Valid: Numbers that appear only twice are fine (below default threshold of 3).
func TwoTimes1() int {
	return 999
}

func TwoTimes2() int {
	return 999
}
