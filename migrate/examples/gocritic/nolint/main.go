package example

func example() {
	x := 1
	switch true { //nolint:gocritic
	case x == 1:
		_ = x
	}
}
