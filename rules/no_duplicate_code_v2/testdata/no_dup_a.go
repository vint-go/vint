package fixtures

import "fmt"

func uniqueAddition(a, b int) string {
	sum := a + b
	return fmt.Sprintf("%d + %d = %d", a, b, sum)
}
