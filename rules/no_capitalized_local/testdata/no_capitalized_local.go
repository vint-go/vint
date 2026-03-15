package fixtures

func f(
	IN int, // MATCH /local variable IN should not be capitalized/
	OUT *int, // MATCH /local variable OUT should not be capitalized/
) (
	ERR error, // MATCH /local variable ERR should not be capitalized/
) {
	return nil
}

func g(in int, out *int) (err error) {
	return nil
}

func h(_ int) {
}

func withBlank(_ int, Name string) { // MATCH /local variable Name should not be capitalized/
}
