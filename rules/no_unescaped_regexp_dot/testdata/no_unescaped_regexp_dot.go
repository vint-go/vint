package fixtures

import "regexp"

func invalidRegexps() {
	re := regexp.MustCompile(`google.com`)       // MATCH /unescaped dot in regexp pattern before 'com', use '\.' to match a literal dot/
	re = regexp.MustCompile(`example.org`)        // MATCH /unescaped dot in regexp pattern before 'org', use '\.' to match a literal dot/
	re = regexp.MustCompile(`mysite.net`)          // MATCH /unescaped dot in regexp pattern before 'net', use '\.' to match a literal dot/
	re = regexp.MustCompile(`website.info`)        // MATCH /unescaped dot in regexp pattern before 'info', use '\.' to match a literal dot/
	re = regexp.MustCompile(`school.edu`)          // MATCH /unescaped dot in regexp pattern before 'edu', use '\.' to match a literal dot/
	re = regexp.MustCompile(`startup.io`)          // MATCH /unescaped dot in regexp pattern before 'io', use '\.' to match a literal dot/
	_, _ = regexp.Compile(`domain.com/path`)       // MATCH /unescaped dot in regexp pattern before 'com', use '\.' to match a literal dot/
	_ = re
}

func validRegexps() {
	re := regexp.MustCompile(`google\.com`)    // dot is properly escaped
	re = regexp.MustCompile(`example\.org`)     // dot is properly escaped
	re = regexp.MustCompile(`[.]com`)           // dot inside character class is literal
	re = regexp.MustCompile(`.*`)               // dot followed by quantifier, not a domain
	re = regexp.MustCompile(`.+`)               // dot followed by quantifier, not a domain
	re = regexp.MustCompile(`foo\.bar\.com`)    // all dots escaped
	re = regexp.MustCompile(`\d+\.\d+`)         // numeric pattern with escaped dots
	_ = re
}
