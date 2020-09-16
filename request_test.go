package applesignin_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	applesignin "github.com/albenik/apple-signin-go"
)

func TestValidateTokenRequest_Encode(t *testing.T) {
	r := &applesignin.ValidateTokenRequest{
		ClientId:     "a1",
		ClientSecret: "b2",
		Code:         "c3",
		GrantType:    "d4",
		RedirectURI:  "e5",
	}
	assert.Equal(t, "client_id=a1&client_secret=b2&code=c3&grant_type=d4&redirect_uri=e5", string(r.Encode()))
}

func TestRefreshTokenRequest_Encode(t *testing.T) {
	r := applesignin.RefreshTokenRequest{
		ClientId:     "a1",
		ClientSecret: "b2",
		GrantType:    "c3",
		RefreshToken: "d4",
	}
	assert.Equal(t, "client_id=a1&client_secret=b2&grant_type=c3&refresh_token=d4", string(r.Encode()))
}
