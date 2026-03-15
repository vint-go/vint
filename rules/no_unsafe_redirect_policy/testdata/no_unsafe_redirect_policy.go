package fixtures

import "net/http"

// Invalid: Client that blindly follows redirects with all headers
func badCopyAllHeaders() {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error { // MATCH /unsafe redirect policy may propagate sensitive headers to different domains/
			for key, val := range via[0].Header {
				req.Header[key] = val
			}
			return nil
		},
	}
	_ = client
}

// Valid: Client that strips sensitive headers on cross-domain redirects
func goodStripSensitiveHeaders() {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) > 0 && req.URL.Host != via[0].URL.Host {
				req.Header.Del("Authorization")
				req.Header.Del("Cookie")
			}
			return nil
		},
	}
	_ = client
}

// Valid: Disabling automatic redirects entirely
func goodDisableRedirects() {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	_ = client
}

// Valid: Client with no CheckRedirect (uses default behavior)
func goodDefaultClient() {
	client := &http.Client{}
	_ = client
}

// Valid: Client that copies headers but checks host first
func goodCopyWithHostCheck() {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if req.URL.Host == via[0].URL.Host {
				for key, val := range via[0].Header {
					req.Header[key] = val
				}
			}
			return nil
		},
	}
	_ = client
}
