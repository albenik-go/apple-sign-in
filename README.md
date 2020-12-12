# Sign in with Apple Golang client library

## Manual auth check

```
curl -v https://appleid.apple.com/auth/authorize?response_type=code&response_mode=query&scope=name%20email&client_id=[CLIENT_ID]&state=[STATE]&redirect_uri=[REDIRECT_URI]
```

CLIENT_ID: 
