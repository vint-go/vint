package fixtures

import (
	"fmt"
	"os"
	"log"
)

// Invalid: code after return is unreachable.
func afterReturn() int {
	return 42 // MATCH /unreachable code after this statement/
	fmt.Println("this is unreachable")
}

// Invalid: code after panic is unreachable.
func afterPanic() {
	panic("fatal error") // MATCH /unreachable code after this statement/
	cleanup()
}

// Invalid: code after os.Exit is unreachable.
func afterOsExit() {
	os.Exit(1) // MATCH /unreachable code after this statement/
	fmt.Println("unreachable")
}

// Invalid: code after log.Fatal is unreachable.
func afterLogFatal() {
	log.Fatal("fatal") // MATCH /unreachable code after this statement/
	fmt.Println("unreachable")
}

// Invalid: code after infinite loop with no break.
func afterInfiniteLoop() {
	for { // MATCH /unreachable code after this statement/
		doWork()
	}
	fmt.Println("unreachable")
}

// Valid: code after return is the last statement in the function.
func lastReturn() int {
	fmt.Println("computing result")
	return 42
}

// Valid: reachable in both branches.
func conditionalReturn(x int) int {
	if x > 0 {
		return x
	}
	return -x
}

// Valid: return after os.Exit to satisfy function signature.
func exitWithReturn() int {
	os.Exit(1)
	return 0
}

// Valid: return after panic to satisfy function signature.
func panicWithReturn() int {
	panic("error")
	return 0
}

// Valid: loop with a break is not necessarily infinite.
func loopWithBreak() {
	for {
		if someCondition() {
			break
		}
		doWork()
	}
	fmt.Println("reachable after loop with break")
}

// Valid: for loop with a condition.
func conditionalLoop() {
	for i := 0; i < 10; i++ {
		doWork()
	}
	fmt.Println("reachable after conditional loop")
}

// Valid: return after log.Fatal to satisfy function signature.
func logFatalWithReturn() int {
	log.Fatal("fatal")
	return 0
}

// Valid: local variable named "log" calling Panic is NOT stdlib log.Panic.
// This should NOT be flagged as unreachable (e.g., *zap.SugaredLogger).
func localVarNamedLog() {
	log := getLogger()
	log.Panic("something went wrong")
	cleanup() // reachable: log is a local variable, not the log package
}

// Valid: local variable named "t" calling Fatal is NOT testing.T.Fatal.
func localVarNamedT() {
	t := getTracer()
	t.Fatal("trace error")
	cleanup() // reachable: t is a local variable, not *testing.T
}

// helper stubs
func cleanup()              {}
func doWork()               {}
func someCondition() bool   { return false }
func getLogger() interface{ Panic(args ...interface{}) } { return nil }
func getTracer() interface{ Fatal(args ...interface{}) } { return nil }
