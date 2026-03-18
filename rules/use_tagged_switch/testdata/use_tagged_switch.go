package fixtures

func processInvalid(x int) string {
	switch { // MATCH /could convert untagged switch to tagged switch on x/
	case x == 1:
		return "one"
	case x == 2:
		return "two"
	default:
		return "other"
	}
}

func processValid(x int) string {
	switch x {
	case 1:
		return "one"
	case 2:
		return "two"
	default:
		return "other"
	}
}

func mixedConditions(x int, y int) string {
	// Different variables in conditions - should not trigger
	switch {
	case x == 1:
		return "one"
	case y == 2:
		return "two"
	default:
		return "other"
	}
}

func nonEqualComparisons(x int) string {
	// Non-== comparisons - should not trigger
	switch {
	case x > 1:
		return "big"
	case x < 0:
		return "negative"
	default:
		return "other"
	}
}

func defaultOnly() string {
	// Only default case - should not trigger
	switch {
	default:
		return "default"
	}
}

func selectorExpr(s struct{ Field int }) string {
	switch { // MATCH /could convert untagged switch to tagged switch on s.Field/
	case s.Field == 1:
		return "one"
	case s.Field == 2:
		return "two"
	default:
		return "other"
	}
}

func multipleCaseValues(x int) string {
	switch { // MATCH /could convert untagged switch to tagged switch on x/
	case x == 1, x == 2:
		return "one or two"
	case x == 3:
		return "three"
	default:
		return "other"
	}
}

func noDefaultCase(x int) string {
	switch { // MATCH /could convert untagged switch to tagged switch on x/
	case x == 1:
		return "one"
	case x == 2:
		return "two"
	}
	return "other"
}

func boolExpression(b bool) {
	// Boolean expression, not == comparison - should not trigger
	switch {
	case b:
		_ = 1
	case !b:
		_ = 2
	}
}

func functionCallInCase(x int) string {
	// Non-binary expression in case - should not trigger
	switch {
	case x == 1:
		return "one"
	case x > 2:
		return "big"
	default:
		return "other"
	}
}

func taggedSwitchAlready(x int) string {
	// Already a tagged switch - should not trigger
	switch x {
	case 1:
		return "one"
	case 2:
		return "two"
	default:
		return "other"
	}
}

func emptySwitch() {
	// Empty switch body - should not trigger
	switch {
	}
}

func stringComparison(s string) string {
	switch { // MATCH /could convert untagged switch to tagged switch on s/
	case s == "hello":
		return "greeting"
	case s == "bye":
		return "farewell"
	default:
		return "unknown"
	}
}
