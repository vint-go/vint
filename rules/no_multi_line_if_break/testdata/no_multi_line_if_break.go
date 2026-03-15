package fixtures

func doSomething()      {}
func grantAccess(u any) {}

var a, b, c int

type user struct{}

func (u user) IsActive() bool            { return false }
func (u user) HasPermission(s string) bool { return false }

func invalidIfBreak1() {
	if // MATCH /multi-line if condition should start on the same line as the if keyword/
		a == 1 &&
		b == 2 &&
		c == 3 {
		doSomething()
	}
}

func invalidIfBreak2() {
	u := user{}
	if // MATCH /multi-line if condition should start on the same line as the if keyword/
		u.IsActive() &&
		u.HasPermission("admin") {
		grantAccess(u)
	}
}

func validIfMultiLine1() {
	if a == 1 &&
		b == 2 &&
		c == 3 {
		doSomething()
	}
}

func validIfMultiLine2() {
	u := user{}
	if u.IsActive() &&
		u.HasPermission("admin") {
		grantAccess(u)
	}
}

func validIfSingleLine() {
	if a == 1 {
		doSomething()
	}
}
