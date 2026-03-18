package fixtures

func badReturnNilVar(x *int) *int {
	if x == nil {
		return x // MATCH /returning a nil variable 'x' instead of explicit nil/
	}
	return x
}

func badReturnNilVarReversed(x *int) *int {
	if nil == x {
		return x // MATCH /returning a nil variable 'x' instead of explicit nil/
	}
	return x
}

func goodExplicitNilReturn(x *int) *int {
	if x == nil {
		return nil
	}
	return x
}

func goodNotEqualCheck(x *int) *int {
	if x != nil {
		return x
	}
	return nil
}

func goodNoNilCheck(x *int) *int {
	if x == new(int) {
		return x
	}
	return x
}

type myError struct{}

func (e *myError) Error() string { return "error" }

func badReturnNilError(err error) error {
	if err == nil {
		return err // MATCH /returning a nil variable 'err' instead of explicit nil/
	}
	return err
}

func goodReturnExplicitNilError(err error) error {
	if err == nil {
		return nil
	}
	return err
}
