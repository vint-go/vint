package fixtures

import "os"

func badDeferInInfiniteForLoop() {
	for {
		f, _ := os.Open("file.txt")
		defer f.Close() // MATCH /defers in infinite loops will never execute/
	}
}

func badDeferInForTrueLoop() {
	for true {
		f, _ := os.Open("file.txt")
		defer f.Close() // MATCH /defers in infinite loops will never execute/
	}
}

func badMultipleDeferInInfiniteLoop() {
	for {
		f, _ := os.Open("file.txt")
		defer f.Close()   // MATCH /defers in infinite loops will never execute/
		defer func() {}() // MATCH /defers in infinite loops will never execute/
	}
}

func badNestedInfiniteLoop() {
	for {
		for {
			f, _ := os.Open("file.txt")
			defer f.Close() // MATCH /defers in infinite loops will never execute/
		}
	}
}

func badDeferInOuterInfiniteLoop() {
	for {
		f, _ := os.Open("file.txt")
		defer f.Close() // MATCH /defers in infinite loops will never execute/
		for i := 0; i < 10; i++ {
			_ = i
		}
	}
}

// Valid: defer outside any loop
func goodDeferOutsideLoop() error {
	f, err := os.Open("file.txt")
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}

// Valid: defer in a closure inside an infinite loop
func goodDeferInClosureInsideInfiniteLoop() {
	for {
		func() {
			f, _ := os.Open("file.txt")
			defer f.Close()
		}()
	}
}

// Valid: defer in a finite for loop (this rule only targets infinite loops)
func goodDeferInFiniteLoop(filenames []string) error {
	for i := 0; i < len(filenames); i++ {
		f, err := os.Open(filenames[i])
		if err != nil {
			return err
		}
		defer f.Close()
	}
	return nil
}

// Valid: defer in a range loop (not infinite)
func goodDeferInRangeLoop(filenames []string) error {
	for _, filename := range filenames {
		f, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer f.Close()
	}
	return nil
}

// Valid: defer in a for loop with a condition
func goodDeferInConditionalLoop() {
	x := true
	for x {
		f, _ := os.Open("file.txt")
		defer f.Close()
		x = false
	}
}

// Valid: extract into a separate function
func processFile() {
	f, _ := os.Open("file.txt")
	defer f.Close()
}

func goodExtractedFunction() {
	for {
		processFile()
	}
}
