package fixtures

import "runtime"

type Resource struct {
	name string
}

func cleanup(r *Resource) {}

// Bad: closure captures r, preventing GC
func badFinalizerCapture() *Resource {
	r := &Resource{name: "test"}
	runtime.SetFinalizer(r, func(_ *Resource) { // MATCH /the finalizer references the finalized object, preventing garbage collection/
		cleanup(r)
	})
	return r
}

// Good: uses the parameter, not the closure variable
func goodFinalizerUsesParam() *Resource {
	r := &Resource{name: "test"}
	runtime.SetFinalizer(r, func(r *Resource) {
		cleanup(r)
	})
	return r
}

// Good: uses a non-closure function
func goodFinalizerNonClosure() *Resource {
	r := &Resource{name: "test"}
	runtime.SetFinalizer(r, cleanup)
	return r
}

// Good: closure does not reference the finalized object
func goodFinalizerNoCapture() *Resource {
	r := &Resource{name: "test"}
	other := &Resource{name: "other"}
	runtime.SetFinalizer(r, func(_ *Resource) {
		cleanup(other)
	})
	return r
}
