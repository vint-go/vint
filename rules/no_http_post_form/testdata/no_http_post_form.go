package fixtures

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func badHttpPostForm() {
	resp, err := http.PostForm("https://example.com/form", url.Values{ // MATCH /http.PostForm does not accept a context; use http.NewRequestWithContext and (*http.Client).Do instead/
		"username": {"user"},
		"password": {"pass"},
	})
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println(resp.Status)
}

func goodHttpPostFormWithContext() {
	ctx := context.Background()
	data := url.Values{
		"username": {"user"},
		"password": {"pass"},
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://example.com/form", strings.NewReader(data.Encode()))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println(resp.Status)
}
