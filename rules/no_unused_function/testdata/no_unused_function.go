package fixtures

// Invalid: unexported function never called or referenced from any entry point.

func unusedHelper() int { // MATCH /func unusedHelper is unused/
	return 42
}

func anotherUnused() string { // MATCH /func anotherUnused is unused/
	return "hello"
}

// Invalid: mutually recursive but never called from an entry point.

func mutualA() int { // MATCH /func mutualA is unused/
	return mutualB()
}

func mutualB() int { // MATCH /func mutualB is unused/
	return mutualA()
}

// Valid: exported functions are considered used by the package.

func PublicAPI() string {
	return internalHelper()
}

func internalHelper() string {
	return "data"
}

// Valid: function is called by an exported function (transitive reachability).

func ExportedCaller() string {
	return callee()
}

func callee() string {
	return "result"
}

// Valid: function used as a value, reachable from exported function.

func ExportedUseFunc() int {
	return apply(myFunc)
}

func apply(f func() int) int {
	return f()
}

func myFunc() int {
	return 1
}

// Valid: unexported function called from an exported method should be reachable.

type MyService struct{}

func (s *MyService) Run() string {
	return methodCalledHelper()
}

func methodCalledHelper() string {
	return "called from method"
}

// Valid: unexported function called transitively from an unexported method
// that is called from an exported method.

func (s *MyService) Process() int {
	return unexportedMethodHelper()
}

func unexportedMethodHelper() int {
	return deepHelper()
}

func deepHelper() int {
	return 99
}

// Invalid: unexported function called only from an unexported method on an
// unexported type — the method itself is unreachable.

type hiddenService struct{}

func (h *hiddenService) run() string {
	return unreachableViaMethod()
}

func unreachableViaMethod() string { // MATCH /func unreachableViaMethod is unused/
	return "nobody calls me"
}
