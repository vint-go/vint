package fixtures

func firstItem(items []string) string {
	for _, item := range items { // MATCH /loop exits unconditionally after one iteration/
		return item
	}
	return ""
}

func firstItemGood(items []string) string {
	if len(items) > 0 {
		return items[0]
	}
	return ""
}

func breakInForLoop() {
	items := []int{1, 2, 3}
	for _, v := range items { // MATCH /loop exits unconditionally after one iteration/
		_ = v
		break
	}
}

func conditionalBreak() {
	items := []int{1, 2, 3}
	for _, v := range items {
		if v > 2 {
			break
		}
	}
}

func conditionalReturn(items []string) string {
	for _, item := range items {
		if item == "target" {
			return item
		}
	}
	return ""
}

func forStmtAlwaysReturn() int {
	for i := 0; i < 10; i++ { // MATCH /loop exits unconditionally after one iteration/
		return i
	}
	return -1
}

func forStmtAlwaysBreak() {
	for i := 0; i < 10; i++ { // MATCH /loop exits unconditionally after one iteration/
		_ = i
		break
	}
}

func normalLoop() {
	for i := 0; i < 10; i++ {
		_ = i
	}
}

func loopWithContinue() {
	for i := 0; i < 10; i++ {
		if i%2 == 0 {
			continue
		}
		_ = i
	}
}

func ifElseBothExit(items []int) int {
	for _, v := range items { // MATCH /loop exits unconditionally after one iteration/
		if v > 0 {
			return v
		} else {
			return -v
		}
	}
	return 0
}

func ifElseOnlyThenExits(items []int) int {
	for _, v := range items {
		if v > 0 {
			return v
		} else {
			_ = v
		}
	}
	return 0
}

func emptyLoop() {
	for i := 0; i < 10; i++ {
	}
}

func infiniteLoop() {
	for {
		break
	}
}

func rangeWithPanic(items []string) {
	for range items { // MATCH /loop exits unconditionally after one iteration/
		panic("always panics")
	}
}
