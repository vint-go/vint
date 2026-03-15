package fixtures

import (
	"context"
	"fmt"
	"net/http"
)

func badHttpNewRequest() {
	req, err := http.NewRequest(http.MethodGet, "https://example.com", nil) // MATCH /http.NewRequest does not accept a context; use http.NewRequestWithContext instead/
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println(resp.Status)
}

func goodHttpNewRequestWithContext() {
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://example.com", nil)
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println(resp.Status)
}
