package fixtures

func redundantNilCheckSlice(s []int) bool {
	return s != nil && len(s) > 0 // MATCH /redundant nil check on slice; len() handles nil slices/
}

func redundantNilCheckSliceNeq(s []int) bool {
	return s != nil && len(s) != 0 // MATCH /redundant nil check on slice; len() handles nil slices/
}

func redundantNilCheckSliceGeq(s []int) bool {
	return s != nil && len(s) >= 1 // MATCH /redundant nil check on slice; len() handles nil slices/
}

func redundantNilCheckMap(m map[string]int) bool {
	return m != nil && len(m) > 0 // MATCH /redundant nil check on slice; len() handles nil slices/
}

func validLenCheckOnly(s []int) bool {
	return len(s) > 0
}

func validNilCheckOnly(s []int) bool {
	return s != nil
}

func validNilCheckWithOtherCondition(s []int) bool {
	return s != nil && s[0] == 1
}

func validDifferentVars(s []int, t []int) bool {
	return s != nil && len(t) > 0
}

func validNilEqualCheck(s []int) bool {
	return s == nil || len(s) == 0
}

func validLenCompareToOther(s []int) bool {
	return s != nil && len(s) > 5
}

func validLenCompareLessThan(s []int) bool {
	return s != nil && len(s) < 10
}
