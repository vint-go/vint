package fixtures

func checkBoolComparisons(b bool) {
	if b == true { // MATCH /omit explicit comparison to boolean constant, can be simplified to b/
		// ...
	}
	if b == false { // MATCH /omit explicit comparison to boolean constant, can be simplified to !b/
		// ...
	}
	if b != true { // MATCH /omit explicit comparison to boolean constant, can be simplified to !b/
		// ...
	}
	if b != false { // MATCH /omit explicit comparison to boolean constant, can be simplified to b/
		// ...
	}
}

func checkBoolComparisonsReversed(b bool) {
	if true == b { // MATCH /omit explicit comparison to boolean constant, can be simplified to b/
		// ...
	}
	if false == b { // MATCH /omit explicit comparison to boolean constant, can be simplified to !b/
		// ...
	}
}

func validBoolUsage(b bool) {
	if b {
		// ...
	}
	if !b {
		// ...
	}
}

func checkExpressions(x, y int) {
	if x == y {
		// not a bool comparison
	}
	if x != y {
		// not a bool comparison
	}
}

func checkMethodCallComparison() {
	b := isOk()
	if b == true { // MATCH /omit explicit comparison to boolean constant, can be simplified to b/
		// ...
	}
}

func isOk() bool {
	return true
}
