# Sign in with Apple Golang client library

**Early version may be unstable**

## Install

```bash
go get -u github.com/albenik/apple-signin-go
```

## Test server

Very usefull with [ngrok](https://ngrok.com)

Run:

```bash
go run github.com/albenik/apple-signin-go/cmd/test-server -aud <audience> -team <team_id> -client <client_id> -key <key_id> -keyfile <pem_file_path> -redirect https://<ngrok_subdomain>.ngrok.io/callback
```

Then open http://localhost:8080 and follow instructions

## Resources

* [Sign in with Apple official documentation](https://developer.apple.com/documentation/sign_in_with_apple)
* [What the Heck is Sign In with Apple?](https://developer.okta.com/blog/2019/06/04/what-the-heck-is-sign-in-with-apple) by
  {okta}
* [Ngrok](https://ngrok.com) — Secure introspectable tunnels to localhost
