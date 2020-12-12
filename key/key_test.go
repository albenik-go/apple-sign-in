package key_test

import (
	"crypto/ecdsa"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/albenik/apple-signin-go/key"
)

func TestParsePrivateKey(t *testing.T) {
	const raw = `-----BEGIN PRIVATE KEY-----
MIGTAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBHkwdwIBAQQgusZ/Y029MmQ4mXWn
fnzXUMI/DgtJIJdvG3cZtOsL3pmgCgYIKoZIzj0DAQehRANCAASQloEXsIF31S59
n5/2YdbDaijlx2eIyIfkv7tre3GxgG8NILwvNCrg6L9Tm9JkVjsLucwXcQ+ezINf
YJBJn/t2
-----END PRIVATE KEY-----`

	k, err := key.ParsePrivateFromPEM([]byte(raw))
	require.NoError(t, err)
	assert.IsType(t, (*ecdsa.PrivateKey)(nil), k)
	assert.NotNil(t, k)
}

func TestParsePrivateKey_InvalidKey(t *testing.T) {
	keys := [][]byte{
		nil,
		[]byte(""),
		[]byte("xyz"),
	}

	for i, k := range keys {
		t.Run("key#"+strconv.Itoa(i), func(kk []byte) func(t *testing.T) {
			return func(t *testing.T) {
				_, err := key.ParsePrivateFromPEM(kk)
				assert.Error(t, err)
			}
		}(k))
	}
}
