package applesignin

import (
	"net/http"

	"github.com/dgrijalva/jwt-go/v4"

	"github.com/albenik/apple-signin-go/jwkproxy"
)

func WithBaseURL(u string) func(*Client) {
	return func(c *Client) {
		c.baseURL = u
	}
}

func WithHTTPClient(h *http.Client) func(*Client) {
	return func(c *Client) {
		c.http = h
	}
}

func WithJWTParser(p *jwt.Parser) func(*Client) {
	return func(c *Client) {
		c.jwtparser = p
	}
}

func WithJWKProxy(p jwkproxy.Interface) func(*Client) {
	return func(c *Client) {
		c.jwk = p
	}
}
