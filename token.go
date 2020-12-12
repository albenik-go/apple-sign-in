package applesignin

import (
	"github.com/dgrijalva/jwt-go/v4"
)

type Token struct {
	ExpiresIn    int        `json:"expires_in"`
	IDToken      *jwt.Token `json:"id_token"`
	AccessToken  string     `json:"access_token"`
	RefreshToken string     `json:"refresh_token"`
	TokenType    string     `json:"token_type"`
}
