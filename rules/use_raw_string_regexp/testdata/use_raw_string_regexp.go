package fixtures

import "regexp"

func invalidRegexps() {
	re := regexp.MustCompile("\\d+\\.\\d+")          // MATCH /regexp pattern can be simplified by using a raw string literal (backtick)/
	re = regexp.MustCompile("\\bfoo\\b")              // MATCH /regexp pattern can be simplified by using a raw string literal (backtick)/
	re = regexp.MustCompile("^\\w+$")                 // MATCH /regexp pattern can be simplified by using a raw string literal (backtick)/
	_, _ = regexp.Compile("hello\\s+world")            // MATCH /regexp pattern can be simplified by using a raw string literal (backtick)/
	_ = regexp.MustCompile("\\[section\\]")            // MATCH /regexp pattern can be simplified by using a raw string literal (backtick)/
	_ = re
}

func validRegexps() {
	re := regexp.MustCompile(`\d+\.\d+`)    // raw string literal - no issue
	re = regexp.MustCompile(`\bfoo\b`)       // raw string literal - no issue
	re = regexp.MustCompile("hello world")   // no backslash escapes - no issue
	re = regexp.MustCompile("^[a-z]+$")      // no backslash escapes - no issue
	re = regexp.MustCompile("")               // empty string - no issue
	_ = re
}
