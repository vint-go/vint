package fixtures

import "context"

type contextKey string

const myKey contextKey = "key"

func badStringKey() {
	ctx := context.WithValue(context.Background(), "key", "value") // MATCH /should not use basic type string as key in context.WithValue/
	_ = ctx
}

func badIntKey() {
	ctx := context.WithValue(context.Background(), 42, "value") // MATCH /should not use basic type int as key in context.WithValue/
	_ = ctx
}

func goodCustomTypeKey() {
	// Using custom type as key - no violation
	ctx := context.WithValue(context.Background(), myKey, "value")
	_ = ctx
}

func goodCustomTypeKeyInline() {
	// Using custom type value directly
	ctx := context.WithValue(context.Background(), contextKey("inline"), "value")
	_ = ctx
}
