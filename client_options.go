package applesignin

import (
	"net/http"
	"time"
)

func WithBaseURL(u string) func(c *Client) {
	return func(c *Client) {
		c.baseURL = u
	}
}

func WithHttpTimeout(t time.Duration) func(c *Client) {
	return func(c *Client) {
		c.http.Timeout = t
	}
}

func WithHttpTransport(t http.RoundTripper) func(c *Client) {
	return func(c *Client) {
		c.http.Transport = t
	}
}
