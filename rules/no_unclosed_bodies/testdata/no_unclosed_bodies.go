package fixtures

import (
	"io"
	"net/http"
)

func noUnclosedBodiesGetClosed(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func noUnclosedBodiesGetUnclosed(url string) ([]byte, error) {
	resp, err := http.Get(url) // MATCH /response body must be closed/
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func noUnclosedBodiesClientDoClosed(client *http.Client, req *http.Request) error {
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_ = resp.StatusCode
	return nil
}

func noUnclosedBodiesClientDoUnclosed(client *http.Client, req *http.Request) error {
	resp, err := client.Do(req) // MATCH /response body must be closed/
	if err != nil {
		return err
	}
	_ = resp.StatusCode
	return nil
}

func noUnclosedBodiesPostClosed(url string) error {
	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func noUnclosedBodiesPostUnclosed(url string) error {
	resp, err := http.Post(url, "application/json", nil) // MATCH /response body must be closed/
	if err != nil {
		return err
	}
	_ = resp
	return nil
}

func noUnclosedBodiesBlankIdentifier(url string) error {
	_, err := http.Get(url) // MATCH /response body must be closed/
	return err
}

func noUnclosedBodiesBothBlank(url string) {
	_, _ = http.Get(url) // MATCH /response body must be closed/
}

func noUnclosedBodiesDirectClose(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func noUnclosedBodiesDiscardedGo(url string) {
	go http.Get(url) // MATCH /response body must be closed/
}

func noUnclosedBodiesDiscardedDefer(url string) {
	defer http.Get(url) // MATCH /response body must be closed/
}

func noUnclosedBodiesDiscardedWrapper(client *http.Client, req *http.Request) {
	go client.Do(req) // MATCH /response body must be closed/
}
