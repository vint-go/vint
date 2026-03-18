package fixtures

// Invalid: error is not the last return value
func process() (error, string) { // MATCH /error should be the last type when returning multiple items/
	return nil, "ok"
}

// Invalid: error is first of three return values
func processThree() (error, string, int) { // MATCH /error should be the last type when returning multiple items/
	return nil, "ok", 0
}

// Invalid: error is in the middle
func processMiddle() (string, error, int) { // MATCH /error should be the last type when returning multiple items/
	return "ok", nil, 0
}

// Valid: error is the last return value
func goodProcess() (string, error) {
	return "ok", nil
}

// Valid: single error return
func singleError() error {
	return nil
}

// Valid: no return values
func noReturn() {
}

// Valid: single non-error return
func singleInt() int {
	return 0
}

// Valid: error is the last of three return values
func threeReturns() (string, int, error) {
	return "", 0, nil
}

// Valid: only one return value (not error)
func oneReturn() string {
	return ""
}
