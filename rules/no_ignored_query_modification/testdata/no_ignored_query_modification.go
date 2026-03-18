package fixtures

import "net/url"

func modifyQueryCopyBad(u *url.URL) {
	u.Query().Set("key", "value") // MATCH /(*net/url.URL).Query returns a copy, modifying it doesn't change the URL/
}

func modifyQueryCopyBadAdd(u *url.URL) {
	u.Query().Add("key", "value") // MATCH /(*net/url.URL).Query returns a copy, modifying it doesn't change the URL/
}

func modifyQueryCopyBadDel(u *url.URL) {
	u.Query().Del("key") // MATCH /(*net/url.URL).Query returns a copy, modifying it doesn't change the URL/
}

func modifyQueryGood(u *url.URL) {
	q := u.Query()
	q.Set("key", "value")
	u.RawQuery = q.Encode()
}

func modifyQueryGoodAdd(u *url.URL) {
	q := u.Query()
	q.Add("key", "value")
	u.RawQuery = q.Encode()
}

func readQueryGood(u *url.URL) string {
	return u.Query().Get("key")
}

func readQueryGoodEncode(u *url.URL) string {
	return u.Query().Encode()
}
