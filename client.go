package applesignin

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

const (
	baseURL   = "https://appleid.apple.com/auth/token"
	audience  = "https://appleid.apple.com"
	userAgent = "apple-signin-client-go/1.0"
)

type Client struct {
	audience    string
	teamID      string
	clientID    string
	keyID       string
	keyData     interface{}
	redirectURL string

	userAgent string
	baseURL   string
	http      http.Client
}

func New(tid, cid, kid string, key interface{}) Client {
	return Client{
		audience:    audience,
		teamID:      tid,
		clientID:    cid,
		keyID:       kid,
		keyData:     key,
		redirectURL: "",

		userAgent: userAgent,
		baseURL:   baseURL,
		http:      http.Client{Transport: http.DefaultTransport},
	}
}

func (c Client) WithBaseURL(u string) Client {
	c.baseURL = u
	return c
}

func (c Client) WithHTTPTimeout(t time.Duration) Client {
	c.http.Timeout = t
	return c
}

func (c Client) WithHTTPTransport(t http.RoundTripper) Client {
	c.http.Transport = t
	return c
}

func (c *Client) Auth(code string, ttl time.Duration) (*TokenResponse, error) {
	return c.AuthContext(context.Background(), code, ttl)
}

func (c *Client) AuthContext(ctx context.Context, code string, ttl time.Duration) (*TokenResponse, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"iss": c.teamID,
		"iat": now.Unix(),
		"exp": now.Add(ttl).Unix(),
		"aud": c.audience,
		"sub": c.clientID,
	})
	token.Header["kid"] = c.keyID

	secret, err := token.SignedString(c.keyData)
	if err != nil {
		return nil, errors.Wrap(err, "token signed string error")
	}
	v := ValidateTokenRequest{
		ClientID:     c.clientID,
		ClientSecret: secret,
		Code:         code,
		GrantType:    "authorization_code",
		RedirectURI:  c.redirectURL,
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL, strings.NewReader(v.Encode()))
	if err != nil {
		return nil, errors.Wrap(err, "cannot prepare request")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept-Header", "application/json")
	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "http request failed")
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		payload := new(TokenResponse)
		if err = jsoniter.NewDecoder(resp.Body).Decode(payload); err != nil {
			return nil, errors.Wrap(err, "cannot parse error response json")
		}
		return payload, nil

	case http.StatusBadRequest:
		payload := new(ErrorResponse)
		if err = jsoniter.NewDecoder(resp.Body).Decode(payload); err != nil {
			return nil, errors.Wrap(err, "cannot parse error response json")
		}
		return nil, errors.Wrap(payload, "auth error")

	default:
		body, err := readResponseBodyText(resp)
		if err != nil {
			body = fmt.Sprintf("body read error: %s", err)
		}
		return nil, errors.Errorf("unexpected http response code %d: %q", resp.StatusCode, body)
	}
}
