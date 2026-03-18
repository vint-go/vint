package fixtures

import (
	"context"
	"net/http"
)

func badNilContext() {
	// Passing nil context
	req, _ := http.NewRequestWithContext(nil, "GET", "http://example.com", nil) // MATCH /nil context passed, use context.TODO or context.Background instead/
	_ = req
}

func goodContext() {
	// Using context.Background()
	req, _ := http.NewRequestWithContext(context.Background(), "GET", "http://example.com", nil)
	_ = req
}

func goodContextTODO() {
	// Using context.TODO()
	req, _ := http.NewRequestWithContext(context.TODO(), "GET", "http://example.com", nil)
	_ = req
}

func takesContext(ctx context.Context) {}

func badDirectCall() {
	takesContext(nil) // MATCH /nil context passed, use context.TODO or context.Background instead/
}

func goodDirectCall() {
	takesContext(context.Background())
}
