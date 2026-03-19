package commentspacings

//revive:disable:cyclomatic
type some struct{}

//revive:disable:cyclomatic High complexity score but easy to understand
type some1 struct{}

// Some nice comment
//
// With an empty line in the middle will make the rule panic!
func itsATrap() {}

// This is a well formed comment
type hello struct {
	random `json:random`
}

//This comment does not respect the spacing rule!
// MATCH /no space between comment delimiter and comment text/

var a string

//myOwnDirective: do something

/*
Should be valid
*/

//	Tabs between comment delimiter and comment text should be fine

/*Not valid*/ // MATCH /no space between comment delimiter and comment text/

/*	valid
 */

/* valid
 */

//nolint:staticcheck // nolint should be in the default list of acceptable comments.
var b string

type c struct {
	//+optional
	d *int `json:"d,omitempty"`
}

//extern open
//export MyFunction

//nolint:gochecknoglobals

//this is a regular command that's incorrectly formatted //nolint:foobar // because one two three
// MATCH /no space between comment delimiter and comment text/
