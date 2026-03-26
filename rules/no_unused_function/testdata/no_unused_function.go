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

// Valid: unexported function called from an unexported method that is itself
// called via a selector expression from an exported method.
// This tests that obj.method() calls are resolved in the call graph.

type processor struct{}

func (p *processor) doProcess() int {
	return selectorCalledHelper()
}

func NewProcessor() int {
	p := &processor{}
	return p.doProcess()
}

func selectorCalledHelper() int {
	return 77
}

// Valid: chain of selector calls — exported method calls unexported method
// via selector, which calls a plain function.

type pipeline struct{}

func (p *pipeline) step1() string {
	return p.step2()
}

func (p *pipeline) step2() string {
	return pipelineHelper()
}

func (p *pipeline) Run() string {
	return p.step1()
}

func pipelineHelper() string {
	return "pipeline done"
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

// Invalid: unexported function called only from an unreachable unexported
// method via selector expression — still unreachable.

type isolated struct{}

func (iso *isolated) secret() int {
	return isolatedHelper()
}

func isolatedHelper() int { // MATCH /func isolatedHelper is unused/
	return 0
}
