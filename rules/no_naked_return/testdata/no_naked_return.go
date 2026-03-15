package fixtures

import "fmt"

// longFuncWithNakedReturn has > 5 lines and naked returns.
func longFuncWithNakedReturn(input int) (result int, err error) {
	if input < 0 {
		err = fmt.Errorf("negative")
		return // MATCH /naked return in func with result, err as named returns, use explicit return values/
	}
	result = input * 2
	return // MATCH /naked return in func with result, err as named returns, use explicit return values/
}

// shortFuncWithNakedReturn has <= 5 lines and a naked return - OK.
func shortFuncWithNakedReturn(a, b int) (sum int) {
	sum = a + b
	return
}

// longFuncExplicitReturn has > 5 lines but uses explicit returns - OK.
func longFuncExplicitReturn(input int) (result int, err error) {
	if input < 0 {
		return 0, fmt.Errorf("negative")
	}
	result = input * 2
	result = result + 1
	return result, nil
}

// funcWithoutNamedReturns cannot have naked returns - OK.
func funcWithoutNamedReturns(a, b int) int {
	if a > b {
		return a
	}
	c := a + b
	_ = c
	return b
}

// longClosureWithNakedReturn has a closure exceeding the threshold.
func longClosureWithNakedReturn() {
	_ = func() (count int, err error) {
		count = 1
		count = count + 1
		count = count + 2
		err = fmt.Errorf("oops")
		return // MATCH /naked return in func with count, err as named returns, use explicit return values/
	}
}

// multipleNakedReturns flags each naked return separately.
func multipleNakedReturns(x int) (a int, b string) {
	if x == 0 {
		return // MATCH /naked return in func with a, b as named returns, use explicit return values/
	}
	a = x
	b = "hello"
	return // MATCH /naked return in func with a, b as named returns, use explicit return values/
}
