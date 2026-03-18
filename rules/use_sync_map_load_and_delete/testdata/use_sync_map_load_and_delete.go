package fixtures

import "sync"

func badLoadThenDelete() {
	var m sync.Map
	key := "foo"

	v, ok := m.Load(key) // MATCH /use m.LoadAndDelete(key) instead of separate Load and Delete calls/
	m.Delete(key)
	_, _ = v, ok
}

func badLoadThenDeletePointer() {
	m := &sync.Map{}
	key := "bar"

	v, ok := m.Load(key) // MATCH /use m.LoadAndDelete(key) instead of separate Load and Delete calls/
	m.Delete(key)
	_, _ = v, ok
}

func goodLoadAndDelete() {
	var m sync.Map
	key := "foo"

	v, ok := m.LoadAndDelete(key)
	_, _ = v, ok
}

func goodDifferentKeys() {
	var m sync.Map
	key1 := "foo"
	key2 := "bar"

	v, ok := m.Load(key1)
	m.Delete(key2)
	_, _ = v, ok
}

func goodDifferentReceivers() {
	var m1, m2 sync.Map
	key := "foo"

	v, ok := m1.Load(key)
	m2.Delete(key)
	_, _ = v, ok
}

func goodLoadOnly() {
	var m sync.Map
	key := "foo"

	v, ok := m.Load(key)
	_, _ = v, ok
}

func goodDeleteOnly() {
	var m sync.Map
	key := "foo"

	m.Delete(key)
}

func badInNestedBlock() {
	var m sync.Map
	key := "foo"

	if true {
		v, ok := m.Load(key) // MATCH /use m.LoadAndDelete(key) instead of separate Load and Delete calls/
		m.Delete(key)
		_, _ = v, ok
	}
}
