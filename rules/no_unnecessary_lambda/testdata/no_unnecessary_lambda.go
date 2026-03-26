package fixtures

import (
	"fmt"
	"sort"
)

func less(i, j int) bool {
	return i < j
}

func doWork() {
	fmt.Println("work")
}

func process(i, j int) {
	fmt.Println(i, j)
}

// Invalid: unnecessary lambda wrapping sort comparison
func badSortSlice() {
	xs := []int{3, 1, 2}
	sort.Slice(xs, func(i, j int) bool { // MATCH /unnecessary lambda, use the function directly/
		return less(i, j)
	})
}

// Invalid: unnecessary lambda wrapping a no-arg function in go statement
func badGoStatement() {
	go func() { doWork() }() // MATCH /unnecessary lambda, use the function directly/
}

// Invalid: unnecessary lambda wrapping a function call
func badSimpleWrapper() {
	f := func() { doWork() } // MATCH /unnecessary lambda, use the function directly/
	f()
}

// Invalid: unnecessary lambda wrapping a two-param function
func badTwoParams() {
	f := func(i, j int) { // MATCH /unnecessary lambda, use the function directly/
		process(i, j)
	}
	f(1, 2)
}

// Valid: direct function reference
func goodSortSlice() {
	xs := []int{3, 1, 2}
	sort.Slice(xs, less)
}

// Valid: direct go statement
func goodGoStatement() {
	go doWork()
}

// Valid: lambda body has more than one statement
func goodMultipleStatements() {
	f := func() {
		fmt.Println("before")
		doWork()
	}
	f()
}

// Valid: lambda body is empty
func goodEmptyBody() {
	f := func() {}
	f()
}

// Valid: arguments are not passed in order
func goodDifferentOrder() {
	f := func(i, j int) {
		process(j, i)
	}
	f(1, 2)
}

// Valid: lambda does extra work (different arguments)
func goodExtraArgs() {
	f := func(i int) {
		fmt.Println(i, "extra")
	}
	f(1)
}

// Valid: lambda body is not a call expression (assignment)
func goodNotACallExpr() {
	x := 0
	f := func() {
		x = 1
	}
	f()
	_ = x
}

// Valid: lambda has return type but body is expression statement
func goodWithReturnType() {
	f := func() error {
		return fmt.Errorf("error")
	}
	_ = f()
}

// Valid: different number of args than params
func goodDifferentArgCount() {
	f := func(i, j int) {
		fmt.Println(i)
	}
	f(1, 2)
}

// Valid: inner call uses variadic expansion
func goodVariadicExpansion() {
	f := func(args ...int) {
		fmt.Println(args...)
	}
	f(1, 2, 3)
}

type service struct{}

func (s *service) handle(i int) {
	fmt.Println(i)
}

// Valid: method call on receiver variable (SelectorExpr, not a plain function)
func goodMethodCallOnReceiver() {
	s := &service{}
	f := func(i int) {
		s.handle(i)
	}
	f(1)
}

// Valid: call through a parameter variable (callee is a captured variable)
func goodCallThroughParameter() {
	var next func(int) error
	_ = next
	f := func(next func(int) error, i int) error {
		return next(i)
	}
	_ = f
}

// Valid: callee is a parameter of the enclosing function (captured variable)
func goodCallThroughEnclosingParam() {
	type HandlerFunc func(int) error
	type MiddlewareFunc func(HandlerFunc) HandlerFunc

	var middleware MiddlewareFunc = func(next HandlerFunc) HandlerFunc {
		return func(i int) error {
			return next(i) // next is captured from enclosing func — not a direct ref
		}
	}
	_ = middleware
}

// Valid: callee is a parameter two levels up (deeply nested capture)
func goodCallThroughDeepEnclosingParam() {
	outer := func(handler func(int)) {
		_ = func() {
			inner := func(i int) {
				handler(i) // handler is captured from two scopes up
			}
			inner(1)
		}
	}
	outer(func(int) {})
}

// Invalid: nested lambda wrapping a package-level function (enclosing params don't save it)
func badNestedButPackageLevel() {
	_ = func(next func(int)) {
		f := func() { doWork() } // MATCH /unnecessary lambda, use the function directly/
		_ = f
		_ = next
	}
}

type waiter struct{}

func (w *waiter) WithPort(port string) *waiter { return w }
func (w *waiter) WaitUntilReady(i, j int)      {}

// Valid: method chain (SelectorExpr, eagerly evaluated chain)
func goodMethodChain() {
	w := &waiter{}
	port := "8080"
	f := func(i, j int) {
		w.WithPort(port).WaitUntilReady(i, j)
	}
	f(1, 2)
}
