package fixtures

// Invalid examples - TODO/FIXME/FIX/BUG comments without detail

// TODO
// MATCH /TODO comment without detail or assignee/
func processData() {}

// FIXME
// MATCH /FIXME comment without detail or assignee/
func brokenFunction() {}

// FIX
// MATCH /FIX comment without detail or assignee/
func fixMe() {}

// BUG
// MATCH /BUG comment without detail or assignee/
func buggyFunction() {}

// TODO:
// MATCH /TODO comment without detail or assignee/
func noDetailAfterColon() {}

// FIXME -
// MATCH /FIXME comment without detail or assignee/
func noDetailAfterDash() {}

// Valid examples - TODO/FIXME comments with detail or assignee

// TODO(jsmith): implement error handling for edge cases
func validProcessData() {}

// FIXME(team): this function panics when input is empty, needs guard clause
func validBrokenFunction() {}

// TODO: implement caching layer for database queries
func validTodoWithDetail() {}

// BUG(go.dev/issue/1234): known issue with concurrent access
func validBugWithDetail() {}

// This is a normal comment mentioning TODOS or BUGFIX
func normalComment() {}
