package example

func check(b bool) bool {
	if b == true { //nolint:staticcheck
		return false
	}
	return true
}
