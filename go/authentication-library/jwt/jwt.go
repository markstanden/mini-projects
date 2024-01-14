package jwt

/*
	Token is the struct that holds all of the data to be written to the JWT
	Header is an embedded struct containing the header section of the JWT (alg, typ)
	Payload is an embedded struct containing the indentifying information of issuer (iss), user (sub), jwt (jti), and secret key (kid)
*/
type Token struct {
	Header
	Payload
	Config
	Log Log
}

/*
	Header contains the required standard JWT header fields
	They are used when decoding to identify the algorithm used to sign the token,
	and the token type which in other circumstances may not be a jwt
*/
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

type Config struct {
	/*
		validFrom is the Unix time representation of the point the server started issuing tokens
		therefore any token with a timestamp earlier than this should be rejected.
		Defaults to the 1st Jan 2021
	*/
	ValidFrom int64
	/*
		Lifespan (int64) is the duration (in seconds) that the token will be valid for.
		Negative values for lifespan result in an immediately expired token.
	*/
	Lifespan int64
}

/*
	Log is an array of messages added to as the
	token moves through creation or verification.
	Used mainly for fault finding
*/
type Log []string
