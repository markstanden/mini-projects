package token

// Token is the struct that holds all of the data to be written to the JWT
type Token struct {
	Header
	Payload
}

// Header contains the required standard JWT fields
// Header.ALG (Algorithm) The encoding type used within the JWT
// It is important that the encoding method is checked to be as expected prior to decoding.
// Header.CTY (Content Type) Used only in nested JWT operations
// Header.TYP (Type) Set to "JWT" for JWT operations, allows for the use of encoding tokens for other uses.
type Header struct {
	Algorithm   string `json:"alg,omitempty"`
	ContentType string `json:"cty,omitempty"`
	TokenType   string `json:"typ,omitempty"`
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
	Issuer string `json:"iss,omitempty"`

	// SUB - subject
	// who the JWT was supplied to.
	// Should be a unique identifier
	Subject string `json:"sub,omitempty"`

	// AUD - audience
	// who the JWT is intended for.
	// Should be rejected if the principal processing
	// the claim does not identify itself with
	// the value listed here.
	Audience string `json:"aud,omitempty"`

	// EXP - expiration time
	// the time the JWT ceases to be valid
	ExpirationTime string `json:"exp,omitempty"`

	// NBT - OPTIONAL - not before time
	// the time the begins to be valid
	NotBeforeTime string `json:"nbt,omitempty"`

	// IAT - OPTIONAL - issued at time
	// the time the JWT was issued
	IssuedAtTime string `json:"iat,omitempty"`

	// JTI - OPTIONAL - JWT ID
	// The unique identifier for the JWT
	TokenID string `json:"jti,omitempty"`

	// *** public claims ***
	// Public claims are collision resistant (i.e. URI namespaced)
	// or defined in the "IANA JSON Web Token Registry"
	// https://www.iana.org/assignments/jwt/jwt.xhtml
	KeyID string `json:"kid,omitempty"`

	// *** private claims ***
	// Custom claims specific to our web app.

}
