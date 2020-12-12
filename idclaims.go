package applesignin

import (
	"bytes"

	"github.com/dgrijalva/jwt-go/v4"
	"go.uber.org/multierr"
)

var trueStr = []byte("true")

type IDTokenClaims struct {
	Iss            string      `json:"iss"`
	Sub            string      `json:"sub"`
	Aud            string      `json:"aud"`
	Iat            int64       `json:"iat"`
	Exp            interface{} `json:"exp"`
	Nonce          string      `json:"nonce"`
	NonceSupported bool        `json:"nonce_supported"`
	AuthTime       int64       `json:"auth_time"`
	AtHash         string      `json:"at_hash"`
	Email          string      `json:"email"`
	EmailVerified  FlexBool    `json:"email_verified"`   // apple returns as string "true"
	EmailPrivate   FlexBool    `json:"is_private_email"` // apple returns as string "true"
	RealUserStatus int8        `json:"real_user_status"`
}

// Valid validates standard claims using jwt.ValidationHelper
// Validates time based claims "exp" (see: jwt.WithLeeway)
// Validates "aud" if present in claims. (see: jwt.WithAudience, jwt.WithoutAudienceValidation)
// Validates "iss" if option is provided (see: jwt.WithIssuer).
func (c *IDTokenClaims) Valid(h *jwt.ValidationHelper) error {
	var vErr error

	if h == nil {
		h = jwt.DefaultValidationHelper
	}

	exp, err := jwt.ParseTime(c.Exp)
	if err != nil {
		return err
	}

	if err = h.ValidateExpiresAt(exp); err != nil {
		vErr = multierr.Append(err, vErr)
	}

	var aud jwt.ClaimStrings
	if aud, err = jwt.ParseClaimStrings(c.Aud); err == nil && aud != nil {
		// If it's present and well formed, validate
		if err = h.ValidateAudience(aud); err != nil {
			vErr = multierr.Append(err, vErr)
		}
	} else if err != nil {
		// If it's present and not well formed, return an error
		return &jwt.MalformedTokenError{Message: "couldn't parse 'aud' value"}
	}

	if err = h.ValidateIssuer(c.Iss); err != nil {
		vErr = multierr.Append(err, vErr)
	}

	return vErr
}

type FlexBool bool

func (b *FlexBool) UnmarshalText(s []byte) error {
	*b = FlexBool(bytes.Equal(bytes.ToLower(bytes.Trim(s, "\"")), trueStr))
	return nil
}
