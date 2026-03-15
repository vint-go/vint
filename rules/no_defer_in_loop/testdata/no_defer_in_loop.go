package fixtures

import "os"

func badDeferInForLoop(filenames []string) error {
	for i := 0; i < len(filenames); i++ {
		f, err := os.Open(filenames[i])
		if err != nil {
			return err
		}
		defer f.Close() // MATCH /defer statement inside a loop; deferred call executes only when the function returns, not on each iteration/
	}
	return nil
}

func badDeferInRangeLoop(filenames []string) error {
	for _, filename := range filenames {
		f, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer f.Close() // MATCH /defer statement inside a loop; deferred call executes only when the function returns, not on each iteration/
	}
	return nil
}

func badDeferInNestedLoop(filenames []string) error {
	for _, filename := range filenames {
		for i := 0; i < 3; i++ {
			f, err := os.Open(filename)
			if err != nil {
				return err
			}
			defer f.Close() // MATCH /defer statement inside a loop; deferred call executes only when the function returns, not on each iteration/
		}
	}
	return nil
}

func goodDeferOutsideLoop() error {
	f, err := os.Open("file.txt")
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}

func goodDeferInClosureInsideLoop(filenames []string) error {
	for _, filename := range filenames {
		func() {
			f, err := os.Open(filename)
			if err != nil {
				return
			}
			defer f.Close()
		}()
	}
	return nil
}

func goodDeferInSeparateFunction(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}
