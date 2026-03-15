package fixtures

import "reflect"

func badEqualityComparison(a, b reflect.Value) bool {
	return a == b // MATCH /avoid using == with reflect.Value, use reflect.Value.Equal or compare with reflect.Value.Interface()/
}

func badDeepEqual(a, b reflect.Value) bool {
	return reflect.DeepEqual(a, b) // MATCH /avoid using reflect.DeepEqual with reflect.Value, use reflect.Value.Equal or compare with reflect.Value.Interface()/
}

func goodEqual(a, b reflect.Value) bool {
	// Good: compares underlying values
	return a.Equal(b)
}

func goodInterface(a, b reflect.Value) bool {
	// Good: extract interface values and compare
	return reflect.DeepEqual(a.Interface(), b.Interface())
}

func goodNonReflectEqual(a, b int) bool {
	// Good: not reflect.Value
	return a == b
}

func goodDeepEqualNonReflect(a, b interface{}) bool {
	// Good: not reflect.Value
	return reflect.DeepEqual(a, b)
}
