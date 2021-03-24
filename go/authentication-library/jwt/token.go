package token

// Token is the struct that holds all of the data to be written to the JWT
type Token struct {
	Header
	Payload
	lifespan int64
}

// Header contains the required standard JWT fields
// Header.ALG (Algorithm) The encoding type used within the JWT
// It is important that the encoding method is checked to be as expected prior to decoding.
// Header.CTY (Content Type) Used only in nested JWT operations
// Header.TYP (Type) Set to "JWT" for JWT operations, allows for the use of encoding tokens for other uses.
type Header struct {
	Algorithm string `json:"alg"`
	TokenType string `json:"typ"`
}

// Payload contains the data stored within the JWT
// Note information stored here is not secure,
// it will be transmitted encoded into URLBase64
// ISS - issuer (string || URI),
// SUB (subject) who the JWT was supplied to. (Should be a unique identifier),
// AUD (audience).  Who the JWT is intended for,
// EXP (expiration time) - the time the JWT ceases to be valid,
type Payload struct {

	// *** Registered Claims ***

	// ISS - issuer (string || URI)
	Issuer string `json:"iss"`

	// SUB - subject
	// who the JWT was supplied to.
	// Should be a unique identifier
	Subject string `json:"sub"`

	// AUD - audience
	// who the JWT is intended for.
	// Should be rejected if the principal processing
	// the claim does not identify itself with
	// the value listed here.
	Audience string `json:"aud"`

	// EXP - expiration time
	// the time the JWT ceases to be valid
	ExpirationTime int64 `json:"exp"`

	// NBF - OPTIONAL - not before time
	// the time the begins to be valid
	NotBeforeTime int64 `json:"nbf"`

	// IAT - OPTIONAL - issued at time
	// the time the JWT was issued
	IssuedAtTime int64 `json:"iat"`

	// JTI - OPTIONAL - JWT ID
	// The unique identifier for the JWT
	TokenID string `json:"jti"`

	// *** public claims ***
	// Public claims are collision resistant (i.e. URI namespaced)
	// or defined in the "IANA JSON Web Token Registry"
	// https://www.iana.org/assignments/jwt/jwt.xhtml
	KeyID string `json:"kid"`

	// *** private claims ***
	// Custom claims specific to our web app.

}
