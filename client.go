package applesignin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

const (
	DefaultURL       = "https://appleid.apple.com/auth/token"
	DefaultAudience  = "https://appleid.apple.com"
	DefaultUserAgent = "apple-signin-client-go/1.0"
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
	http      *http.Client
}

func New(tid, cid, kid string, opts ...func(c *Client)) *Client {
	c := &Client{
		audience:    DefaultAudience,
		teamID:      tid,
		clientID:    cid,
		keyID:       kid,
		keyData:     nil,
		redirectURL: "",

		userAgent: DefaultUserAgent,
		baseURL:   DefaultURL,
		http:      &http.Client{Transport: http.DefaultTransport},
	}

	for _, setOpt := range opts {
		setOpt(c)
	}

	return c
}

func (c *Client) Auth(code string, ttl time.Duration) (*TokenResponse, error) {
	// test cert
	// if c.AESCert == nil {
	// 	return nil, errors.New("missing cert")
	// }

	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"iss": c.teamID,
		"iat": now.Unix(),
		"exp": now.Add(ttl).Unix(),
		"aud": c.audience,
		"sub": c.clientID,
	})
	token.Header = map[string]interface{}{"kid": c.keyID, "alg": "ES256"}

	secret, _ := token.SignedString(c.keyData)
	v := ValidateTokenRequest{
		ClientId:     c.clientID,
		ClientSecret: secret,
		Code:         code,
		GrantType:    "authorization_code",
		RedirectURI:  c.redirectURL,
	}

	req, err := http.NewRequest(http.MethodPost, c.baseURL, bytes.NewReader(v.Encode()))
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
		respJson := new(TokenResponse)
		if err = json.NewDecoder(resp.Body).Decode(respJson); err != nil {
			return nil, errors.Wrap(err, "cannot parse error response json")
		}
		return respJson, nil

	case http.StatusBadRequest:
		respJson := new(ErrorResponse)
		if err = json.NewDecoder(resp.Body).Decode(respJson); err != nil {
			return nil, errors.Wrap(err, "cannot parse error response json")
		}
		return nil, errors.Wrap(respJson, "auth error")

	default:
		body, err := readResponseBodyText(resp)
		if err != nil {
			body = fmt.Sprintf("body read error: %s", err)
		}
		return nil, errors.Errorf("unexpected http response code %d: %q", resp.StatusCode, body)
	}
}
