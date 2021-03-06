package applesignin

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	jsoniter "github.com/json-iterator/go"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/pkg/errors"

	"github.com/albenik-go/apple-sign-in/jwkproxy"
)

const (
	baseURL   = "https://appleid.apple.com"
	authURL   = baseURL + "/auth/authorize"
	tokenURL  = baseURL + "/auth/token"
	audience  = baseURL
	userAgent = "apple-signin-client-go/1.0"

	MaxExpiration = 15777000 * time.Second // half a year

	ResponseModeQuery = "query"
	ResponseModePost  = "form_post"

	ScopeEmail = "email"
	ScopeName  = "name"
)

type Client struct {
	audience    string
	teamID      string
	clientID    string
	keyID       string
	keyData     interface{}
	RedirectURL string

	userAgent string
	baseURL   string
	http      *http.Client
	jwtparser *jwt.Parser
	jwk       jwkproxy.Interface
}

// New instatinates a new client.
// Arguments: tid — teamID, cid — clientID, kid — keyID.
func New(tid, cid, kid string, key interface{}, opts ...func(c *Client)) *Client {
	c := &Client{
		audience:    audience,
		teamID:      tid,
		clientID:    cid,
		keyID:       kid,
		keyData:     key,
		RedirectURL: "",

		userAgent: userAgent,
		baseURL:   tokenURL,
		http:      &http.Client{},
		jwtparser: jwt.NewParser(jwt.WithoutAudienceValidation()),
		jwk:       jwkproxy.NewInMemory(jwkproxy.NewLoader(), time.Hour),
	}

	for _, o := range opts {
		o(c)
	}

	return c
}

func (c *Client) AuthURL(mode string, scopes []string, state, nonce string) string {
	q := url.Values{
		"response_type": []string{"code"},
		"response_mode": []string{mode},
		"client_id":     []string{c.clientID},
		"scope":         []string{strings.Join(scopes, " ")}, // "name email"
		"redirect_uri":  []string{c.RedirectURL},
		"state":         []string{state},
		"nonce":         []string{nonce},
	}

	return authURL + "?" + q.Encode()
}

func (c *Client) ValidateCode(code, nonce string, exp time.Duration) (*TokenResponse, error) {
	return c.ValidateCodeContext(context.Background(), code, nonce, exp)
}

func (c *Client) ValidateCodeContext(ctx context.Context, code, nonce string, exp time.Duration) (*TokenResponse, error) { //nolint:lll
	const errmsg = "authorization code validation error"

	res, err := c.doRequest(ctx, &validateTokenRequest{
		apiRequest: apiRequest{
			ClientID:  c.clientID,
			GrantType: "authorization_code",
		},
		Code:        code,
		RedirectURI: c.RedirectURL,
	}, exp)
	if err != nil {
		return nil, errors.Wrap(err, errmsg)
	}

	if res.IDToken == nil { // not impossible by docs, but all possible in real life
		return nil, errors.Wrap(ErrIDTokenMissing, errmsg)
	}

	if res.IDToken.NonceSupported && res.IDToken.Nonce != nonce {
		return nil, errors.Wrap(ErrNonceMismatch, errmsg)
	}

	return res, nil
}

func (c *Client) ValidateRefreshToken(token string, exp time.Duration) (*TokenResponse, error) {
	return c.ValidateRefreshTokenContext(context.Background(), token, exp)
}

func (c *Client) ValidateRefreshTokenContext(ctx context.Context, token string, exp time.Duration) (*TokenResponse, error) { //nolint:lll
	token2, err := c.doRequest(ctx, &refreshTokenRequest{
		apiRequest: apiRequest{
			ClientID:  c.clientID,
			GrantType: "refresh_token",
		},
		RefreshToken: token,
	}, exp)
	if err != nil {
		return nil, errors.Wrap(err, "refresh token validation error")
	}

	return token2, nil
}

func (c *Client) ParseIDToken(token string) (*IDTokenClaims, error) {
	return c.ParseIDTokenContext(context.Background(), token)
}

func (c *Client) ParseIDTokenContext(ctx context.Context, token string) (*IDTokenClaims, error) {
	claims := new(IDTokenClaims)
	if _, err := c.jwtparser.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		keys, err := c.jwk.FetchKeysContext(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "cannot load apple public keys")
		}

		if kid, ok := token.Header["kid"]; ok {
			if alg, ok := token.Header["alg"]; ok {
				for iter := keys.Iterate(ctx); iter.Next(ctx); {
					pair := iter.Pair()
					key, ok := pair.Value.(jwk.Key)
					if ok && key.KeyID() == kid && key.Algorithm() == alg {
						var keydata interface{}
						if err := key.Raw(&keydata); err != nil {
							return nil, errors.Wrap(err, "cannot decode key data")
						}
						return keydata, nil
					}
				}
			}
		}

		return nil, ErrNoSuitableJWK
	}); err != nil {
		return nil, errors.Wrap(err, "ID token parse error")
	}

	return claims, nil
}

func (c *Client) doRequest(ctx context.Context, req request, exp time.Duration) (*TokenResponse, error) {
	if exp > MaxExpiration {
		return nil, errors.Wrap(ErrSecretExpirationTimeTooFar, "cannot retrieve token")
	}

	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"iss": c.teamID,
		"iat": now.Unix(),
		"exp": now.Add(exp).Unix(),
		"aud": c.audience,
		"sub": c.clientID,
	})
	token.Header["kid"] = c.keyID

	secret, err := token.SignedString(c.keyData)
	if err != nil {
		return nil, errors.Wrap(err, "token signed string error")
	}
	req.SetSecret(secret) // nolint:wsl

	httpreq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL, strings.NewReader(req.Encode()))
	if err != nil {
		return nil, errors.Wrap(err, "cannot prepare request")
	}

	httpreq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	httpreq.Header.Set("Accept-Header", "application/json")
	httpreq.Header.Set("User-Agent", c.userAgent)

	resp, err := c.http.Do(httpreq)
	if err != nil {
		return nil, errors.Wrap(err, "http request failed")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var t tokenResponseRaw
		if err = jsoniter.NewDecoder(resp.Body).Decode(&t); err != nil {
			return nil, errors.Wrap(err, "cannot parse error response json")
		}

		var id *IDTokenClaims
		if id, err = c.ParseIDTokenContext(ctx, t.IDToken); err != nil { //nolint:wsl
			return nil, err
		}

		return &TokenResponse{
			ExpiresIn:    t.ExpiresIn,
			IDToken:      id,
			AccessToken:  t.AccessToken,
			RefreshToken: t.RefreshToken,
			TokenType:    t.TokenType,
		}, nil
	}

	if resp.StatusCode == http.StatusBadRequest {
		payload := new(ErrorResponse)
		if err = jsoniter.NewDecoder(resp.Body).Decode(payload); err != nil {
			return nil, errors.Wrap(err, "cannot parse error response json")
		}
		return nil, errors.Wrap(payload, "auth error") //nolint:wsl
	}

	body, err := readResponseBodyText(resp)
	if err != nil {
		body = fmt.Sprintf("body read error: %s", err)
	}
	return nil, errors.Errorf("unexpected http response code %d: %q", resp.StatusCode, body) //nolint:wsl
}
