package fixtures

import (
	"log"
	"os"
)

// Invalid: os.Exit in a function that uses defer.
func badOsExit() {
	f, _ := os.Open("file.txt")
	defer f.Close()
	os.Exit(1) // MATCH /calling os.Exit in a function that uses defer will skip deferred cleanup/
}

// Invalid: log.Fatal in a function that uses defer.
func badLogFatal() {
	f, _ := os.Open("file.txt")
	defer f.Close()
	log.Fatal("error") // MATCH /calling log.Fatal in a function that uses defer will skip deferred cleanup/
}

// Invalid: log.Fatalf in a function that uses defer.
func badLogFatalf() {
	f, _ := os.Open("file.txt")
	defer f.Close()
	log.Fatalf("error: %v", "reason") // MATCH /calling log.Fatalf in a function that uses defer will skip deferred cleanup/
}

// Invalid: log.Fatalln in a function that uses defer.
func badLogFatalln() {
	f, _ := os.Open("file.txt")
	defer f.Close()
	log.Fatalln("error") // MATCH /calling log.Fatalln in a function that uses defer will skip deferred cleanup/
}

// Invalid: multiple exit calls in one function with defer.
func badMultiple() {
	f, _ := os.Open("file.txt")
	defer f.Close()
	if f == nil {
		os.Exit(1) // MATCH /calling os.Exit in a function that uses defer will skip deferred cleanup/
	}
	log.Fatal("done") // MATCH /calling log.Fatal in a function that uses defer will skip deferred cleanup/
}

// Valid: no defer in this function, so os.Exit is fine.
func goodNoDefer() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

// Valid: no defer in this function.
func goodOsExitNoDefer() {
	os.Exit(0)
}

// Valid: the defer is in a separate function.
func goodSeparateFunction() {
	if err := doWork(); err != nil {
		os.Exit(1)
	}
}

func doWork() error {
	f, _ := os.Open("file.txt")
	defer f.Close()
	return nil
}

// Valid: exit call is inside a nested function literal, which does not have defer.
func goodNestedLiteral() {
	defer func() {}()
	fn := func() {
		os.Exit(1) // no defer in the closure itself
	}
	_ = fn
}

// helper stubs
func run() error { return nil }
