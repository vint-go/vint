package fixtures

// Invalid: duplicate condition in if/else if chain
func duplicateCondition(x int) string {
	if x > 10 {
		return "big"
	} else if x > 10 { // MATCH /duplicate condition x > 10 in if/else if chain/
		return "also big"
	}
	return "small"
}

// Invalid: duplicate condition with three branches
func duplicateConditionThree(x int) string {
	if x > 10 {
		return "big"
	} else if x > 5 {
		return "medium"
	} else if x > 10 { // MATCH /duplicate condition x > 10 in if/else if chain/
		return "also big"
	}
	return "small"
}

// Invalid: duplicate condition with complex expression
func duplicateComplexCondition(x, y int) string {
	if x > 10 && y < 5 {
		return "a"
	} else if x > 10 && y < 5 { // MATCH /duplicate condition x > 10 && y < 5 in if/else if chain/
		return "b"
	}
	return "c"
}

// Invalid: multiple duplicates in one chain
func multipleDuplicates(x int) string {
	if x > 10 {
		return "a"
	} else if x > 5 {
		return "b"
	} else if x > 10 { // MATCH /duplicate condition x > 10 in if/else if chain/
		return "c"
	} else if x > 5 { // MATCH /duplicate condition x > 5 in if/else if chain/
		return "d"
	}
	return "e"
}

// Valid: all conditions are different
func allDifferent(x int) string {
	if x > 10 {
		return "big"
	} else if x > 5 {
		return "medium"
	}
	return "small"
}

// Valid: separate if statements (not an if/else if chain)
func separateIfs(x int) string {
	if x > 10 {
		return "big"
	}
	if x > 10 {
		return "also big"
	}
	return "small"
}

// Valid: single if with no else
func singleIf(x int) string {
	if x > 10 {
		return "big"
	}
	return "small"
}

// Valid: different conditions
func differentConditions(x int) string {
	if x > 10 {
		return "big"
	} else if x > 5 {
		return "medium"
	} else if x > 0 {
		return "positive"
	}
	return "non-positive"
}
