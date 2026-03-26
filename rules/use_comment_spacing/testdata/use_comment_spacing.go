package fixtures

//This comment has no space after the slashes
// MATCH /no space between comment delimiter and comment text/

func f() {}

//TODO: fix this later

// This comment has proper spacing
func g() {}

// TODO: fix this later

//go:generate stringer -type=Status

//nolint:foo

//noinspection GoUnusedExportedFunction

//region MyRegion

//endregion MyRegion

//+build linux

//-build linux

//#build linux

//!build linux

//line foo.go:10

/* block comment without space is fine */

// short comment
//
