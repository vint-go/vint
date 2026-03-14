package fixtures

type Item struct {
	Name string
}

// Invalid: taking address of range value variable
func noRangeVariableAlias_bad1(items []Item) []*Item {
	var result []*Item
	for _, item := range items {
		result = append(result, &item) // MATCH /implicit memory aliasing in for...range: taking address of range variable 'item'/
	}
	return result
}

// Invalid: taking address of range value in map iteration
func noRangeVariableAlias_bad2(m map[string]int) []*int {
	var ptrs []*int
	for _, v := range m {
		ptrs = append(ptrs, &v) // MATCH /implicit memory aliasing in for...range: taking address of range variable 'v'/
	}
	return ptrs
}

// Invalid: assigning address of range value to a variable
func noRangeVariableAlias_bad3(items []Item) []*Item {
	var result []*Item
	for _, item := range items {
		p := &item // MATCH /implicit memory aliasing in for...range: taking address of range variable 'item'/
		result = append(result, p)
	}
	return result
}

// Invalid: taking address of range key variable
func noRangeVariableAlias_bad4(items []Item) []*int {
	var result []*int
	for i := range items {
		result = append(result, &i) // MATCH /implicit memory aliasing in for...range: taking address of range variable 'i'/
	}
	return result
}

// Invalid: taking address of field of range variable
func noRangeVariableAlias_bad5(items []Item) []*string {
	var result []*string
	for _, item := range items {
		result = append(result, &item.Name) // MATCH /implicit memory aliasing in for...range: taking address of range variable 'item.Name'/
	}
	return result
}

// Valid: explicit copy before taking address
func noRangeVariableAlias_good1(items []Item) []*Item {
	var result []*Item
	for _, item := range items {
		item := item
		result = append(result, &item)
	}
	return result
}

// Valid: using index to access element
func noRangeVariableAlias_good2(items []Item) []*Item {
	var result []*Item
	for i := range items {
		result = append(result, &items[i])
	}
	return result
}

// Valid: range variable is not used with address-of
func noRangeVariableAlias_good3(items []Item) []Item {
	var result []Item
	for _, item := range items {
		result = append(result, item)
	}
	return result
}

// Valid: taking address of non-range variable
func noRangeVariableAlias_good4(items []Item) []*Item {
	var result []*Item
	for range items {
		x := Item{Name: "test"}
		result = append(result, &x)
	}
	return result
}

// Valid: range variable is a pointer type
func noRangeVariableAlias_good5(items []*Item) []*Item {
	var result []*Item
	for _, item := range items {
		result = append(result, item)
	}
	return result
}
