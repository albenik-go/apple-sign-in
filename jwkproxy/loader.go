package jwkproxy

import (
	"context"
	"net/http"

	"github.com/lestrrat-go/jwx/jwk"
)

const appleKeysURL = "https://appleid.apple.com/auth/keys"

type Loader struct {
	fetchURL string
	http     *http.Client
}

func NewLoader(opts ...func(*Loader)) *Loader {
	l := &Loader{
		fetchURL: appleKeysURL,
		http:     &http.Client{},
	}

	for _, opt := range opts {
		opt(l)
	}

	return l
}

func (l *Loader) FetchContext(ctx context.Context) (jwk.Set, error) {
	return jwk.Fetch(ctx, l.fetchURL, jwk.WithHTTPClient(l.http))
}
