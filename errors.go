package applesignin

// nolint:depguard
import (
	"errors"
)

var (
	ErrSecretExpirationTimeTooFar = errors.New("exp is too far from now")
	ErrNonceMismatch              = errors.New("nonce mismatch")
	ErrNoSuitableJWK              = errors.New("no suitable JWK")
)
