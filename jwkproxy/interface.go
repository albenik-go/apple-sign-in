package jwkproxy

import (
	"context"

	"github.com/lestrrat-go/jwx/jwk"
)

type Interface interface {
	Refresh() error
	RefreshContext(ctx context.Context) error

	FetchKeys() (*jwk.Set, error)
	FetchKeysContext(ctx context.Context) (*jwk.Set, error)
}
