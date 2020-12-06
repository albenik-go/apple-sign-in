package applesignin_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	applesignin "github.com/albenik/apple-signin-go"
)

func TestValidateTokenRequest_Encode(t *testing.T) {
	r1 := applesignin.ValidateTokenRequest{
		ClientID:     "a1",
		ClientSecret: "b2",
		Code:         "c3",
		GrantType:    "d4",
	}
	r2 := r1
	r2.RedirectURI = "e5"

	const result = "client_id=a1&client_secret=b2&code=c3&grant_type=d4"

	t.Run("WithoutRedirectURL", func(t *testing.T) {
		assert.Equal(t, result, r1.Encode())
	})
	t.Run("WithRedirectURL", func(t *testing.T) {
		assert.Equal(t, result+"&redirect_uri=e5", r2.Encode())
	})
}

func TestRefreshTokenRequest_Encode(t *testing.T) {
	r := applesignin.RefreshTokenRequest{
		ClientID:     "a1",
		ClientSecret: "b2",
		GrantType:    "c3",
		RefreshToken: "d4",
	}
	assert.Equal(t, "client_id=a1&client_secret=b2&grant_type=c3&refresh_token=d4", r.Encode())
}
