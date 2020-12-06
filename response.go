package applesignin

// TokenResponse see https://developer.apple.com/documentation/sign_in_with_apple/tokenresponse.
type TokenResponse struct {
	AccessToken  string `json:"access_token"`  // (Reserved for future use) A token used to access allowed data. Currently, no data set has been defined for access.
	ExpiresIn    int    `json:"expires_in"`    // The amount of time, in seconds, before the access token expires.
	IDToken      string `json:"id_token"`      // A JSON Web Token that contains the user’s identity information.
	RefreshToken string `json:"refresh_token"` // The refresh token used to regenerate new access tokens. Store this token securely on your server.
	TokenType    string `json:"token_type"`    // The type of access token. It will always be bearer.
}
