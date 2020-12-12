package jwkproxy

import (
	"context"
	"time"

	"github.com/lestrrat-go/jwx/jwk"
	"github.com/pkg/errors"
)

var _ Interface = (*InMemory)(nil)

type InMemory struct {
	loader *Loader
	ttl    time.Duration
	errTTL time.Duration
	set    *jwk.Set
	expire time.Time
}

func NewInMemory(ldr *Loader, ttl time.Duration) *InMemory {
	return &InMemory{
		loader: ldr,
		ttl:    ttl,
		errTTL: time.Minute,
		set:    nil,
		expire: time.Now().Add(-time.Second),
	}
}

func (m *InMemory) Refresh() error {
	return m.RefreshContext(context.Background())
}

func (m *InMemory) RefreshContext(ctx context.Context) error {
	set, err := m.loader.FetchContext(ctx)
	if err != nil {
		m.expire = time.Now().Add(m.errTTL)
		return errors.Wrap(err, "refresh failed")
	}

	m.expire = time.Now().Add(m.ttl)
	m.set = set
	return nil
}

func (m *InMemory) FetchKeys() (*jwk.Set, error) {
	return m.FetchKeysContext(context.Background())
}

func (m *InMemory) FetchKeysContext(ctx context.Context) (*jwk.Set, error) {
	if time.Now().After(m.expire) {
		if err := m.RefreshContext(ctx); err != nil && m.set == nil {
			return nil, errors.Wrap(err, "cannot fetch keys")
		}
	}
	return m.set, nil
}
