package fixtures

func process(item int)   {}
func transform(item int) {}
func save(item int)      {}
func log(item int)       {}
func notify(item int)    {}
func extra(item int)     {}
func isValid(item int) bool { return true }

// Invalid: for-range with single if, no else, body >= 5 statements
func badRange() {
	items := []int{1, 2, 3}
	for _, item := range items {
		if isValid(item) { // MATCH /invert if condition and use continue to reduce nesting/
			process(item)
			transform(item)
			save(item)
			log(item)
			notify(item)
		}
	}
}

// Invalid: for loop with single if, no else, body >= 5 statements
func badFor() {
	items := []int{1, 2, 3}
	for i := 0; i < len(items); i++ {
		if isValid(items[i]) { // MATCH /invert if condition and use continue to reduce nesting/
			process(items[i])
			transform(items[i])
			save(items[i])
			log(items[i])
			notify(items[i])
		}
	}
}

// Invalid: body with more than 5 statements
func badMoreStatements() {
	items := []int{1, 2, 3}
	for _, item := range items {
		if isValid(item) { // MATCH /invert if condition and use continue to reduce nesting/
			process(item)
			transform(item)
			save(item)
			log(item)
			notify(item)
			extra(item)
		}
	}
}

// Valid: already using early continue
func goodEarlyContinue() {
	items := []int{1, 2, 3}
	for _, item := range items {
		if !isValid(item) {
			continue
		}
		process(item)
		transform(item)
		save(item)
		log(item)
		notify(item)
	}
}

// Valid: if body has fewer than 5 statements
func goodShortBody() {
	items := []int{1, 2, 3}
	for _, item := range items {
		if isValid(item) {
			process(item)
			transform(item)
			save(item)
			log(item)
		}
	}
}

// Valid: has else clause
func goodWithElse() {
	items := []int{1, 2, 3}
	for _, item := range items {
		if isValid(item) {
			process(item)
			transform(item)
			save(item)
			log(item)
			notify(item)
		} else {
			extra(item)
		}
	}
}

// Valid: loop body has more than one top-level statement
func goodMultipleStatements() {
	items := []int{1, 2, 3}
	for _, item := range items {
		log(item)
		if isValid(item) {
			process(item)
			transform(item)
			save(item)
			log(item)
			notify(item)
		}
	}
}

// Valid: not inside a loop
func goodNotInLoop() {
	item := 1
	if isValid(item) {
		process(item)
		transform(item)
		save(item)
		log(item)
		notify(item)
	}
}
