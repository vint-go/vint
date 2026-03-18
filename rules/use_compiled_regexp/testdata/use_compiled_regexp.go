package fixtures

import "regexp"

func badMatchStringInForLoop(items []string) []bool {
	results := make([]bool, len(items))
	for i := 0; i < len(items); i++ {
		results[i], _ = regexp.MatchString(`\d+`, items[i]) // MATCH /calling regexp.MatchString in a loop; compile the regexp once outside the loop with regexp.Compile/
	}
	return results
}

func badMatchStringInRangeLoop(items []string) []bool {
	results := make([]bool, len(items))
	for i, item := range items {
		results[i], _ = regexp.MatchString(`\d+`, item) // MATCH /calling regexp.MatchString in a loop; compile the regexp once outside the loop with regexp.Compile/
	}
	return results
}

func badMatchInForLoop(items [][]byte) []bool {
	results := make([]bool, len(items))
	for i := 0; i < len(items); i++ {
		results[i], _ = regexp.Match(`\d+`, items[i]) // MATCH /calling regexp.Match in a loop; compile the regexp once outside the loop with regexp.Compile/
	}
	return results
}

func badMatchInNestedLoop(items [][]byte) []bool {
	results := make([]bool, len(items))
	for _, item := range items {
		for i := 0; i < 3; i++ {
			results[i], _ = regexp.Match(`\d+`, item) // MATCH /calling regexp.Match in a loop; compile the regexp once outside the loop with regexp.Compile/
		}
	}
	return results
}

func goodCompiledRegexp(items []string) []bool {
	re := regexp.MustCompile(`\d+`)
	results := make([]bool, len(items))
	for i, item := range items {
		results[i] = re.MatchString(item)
	}
	return results
}

func goodMatchStringOutsideLoop() {
	_, _ = regexp.MatchString(`\d+`, "test123")
}

func goodMatchInClosure(items []string) []bool {
	results := make([]bool, len(items))
	for i, item := range items {
		func() {
			results[i], _ = regexp.MatchString(`\d+`, item)
		}()
	}
	return results
}

func goodCompileInLoop(items []string) {
	for _, item := range items {
		re := regexp.MustCompile(item)
		_ = re
	}
}
