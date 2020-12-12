package jwkproxy

import (
	"net/http"
)

func WithFetchURL(u string) func(*Loader) {
	return func(l *Loader) {
		l.fetchURL = u
	}
}

func WithHTTPClient(c *http.Client) func(*Loader) {
	return func(l *Loader) {
		l.http = c
	}
}
