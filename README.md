# Sign in with Apple Golang client library

## Manual auth check

```
curl -v https://appleid.apple.com/auth/authorize?client_id=[CLIENT_ID]&response_type=code&scope=name%20email&response_mode=query&state=[STATE]&redirect_uri=[REDIRECT_URI]
```
