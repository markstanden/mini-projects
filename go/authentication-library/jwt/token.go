package token

// Token is the struct that holds all of the data to be written to the JWT
type Token struct {
	Header
	Payload
	lifespan int64
}

// Header contains the required standard JWT fields
type Header struct {
	// Algorithm - "alg" - The encoding algorithm used to sign the token
	// This is "HS512" and is set automatically
	Algorithm string `json:"alg"`

	// TokenType - "typ" - The type of token to be produced
	// This is set to "JWT" automatically
	TokenType string `json:"typ"`
}

// Payload contains the data stored within the JWT
// Note information stored here is not secure,
// it will be transmitted encoded into URLBase64
type Payload struct {

	// *** Registered Claims ***

	// Issuer - "iss" - issuer (string || URI)
	// The top level domain that issues the token
	Issuer string `json:"iss"`

	// Audience - "aud" - audience
	// who the JWT is intended for.
	// The token will be rejected if the principal processing
	// the claim does not identify itself with
	// the value listed here.
	Audience string `json:"aud"`

	// UserID - "sub" - subject
	// who the JWT was supplied to.
	// Should be a unique identifier
	UserID string `json:"sub"`

	// JwtID - "jti" - JWT ID
	// The unique identifier for this particular token
	JwtID string `json:"jti"`

	// KeyID - "kid" - Key ID
	// ** Public Claim **
	// The version of the secret used to hash the signature.
	KeyID string `json:"kid"`

	// IssuedAtTime - "iat" - issued at time
	// the time the JWT was issued
	// Represented as UNIX time int64 as seconds since the epoch
	IssuedAtTime int64 `json:"iat"`

	// NotBeforeTime - "nbf" - not before time
	// the time the token begins to be valid
	// Represented as UNIX time int64 as seconds since the epoch
	NotBeforeTime int64 `json:"nbf"`

	// ExpirationTime - "exp" - expiration time
	// the time the JWT ceases to be valid
	// Represented as UNIX time int64 as seconds since the epoch
	ExpirationTime int64 `json:"exp"`
}
