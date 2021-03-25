# token
## JWT creation library to be shared across my apps, with sane defaults

A library for basic jwt creation and parsing written in go

I have gone back to basics reading the JWT spec and getting inspiration from jwt.io and Auth0.com

The idea was to create a standard JWT token to be used across multiple projects, using security best practices and sane defaults.

also, at least initially the idea is to have a minimal number of exportable functions, to allow the workings to invisibly change as required.

Currently the API is planned to be

```go
// NewToken returns a new token object with the provided fields, and time fields filled based on the current time.
NewToken(issuer, userID, audience, tokenID, keyID string, expiresInXSeconds int64) (tokenStruct *Token)

// CreateJWT turns a NewToken() into a signed JWT using HMAC SHA512 using the secret obtained by calling the passwordLookup callback with the keyID value
CreateJWT(passwordLookup func(keyID string)(secret string, err error)) (jwt string)

// Decode turns a signed JWT into a *Token
// but only after checking the validity of the token.
// it also requires a callback to lookup the secret the signature was signed with,
// and a pointer to the object that it needs to fill
Decode(jwt string, passwordLookup func(keyID string)(secret string, err error), trustedTokenObject *Token) (err error)

```
