package token

import (
	"errors"
)

type Token struct {
	Header struct {

		// ALG - Algorithm
		// The encoding type used within the JWT
		// It is important that the encoding method
		// is checked to be as expected prior to decoding.
		ALG string

		// CTY - Content Type
		// Used only in nested JWT operations
		CTY string

		// TYP - Type
		// Set to "JWT" for JWT operations, allows
		// for the use of encoding tokens for other uses.
		TYP string
	}

	// Payload contains the data stored within the JWT
	// Note information stored here is not secure,
	// it will be transmitted encoded into URLBase64
	Payload struct {

		// *** Registered Claims ***

		// ISS - issuer (string || URI)
		ISS string

		// SUB - subject
		// who the JWT was supplied to.
		// Should be a unique identifier
		SUB string

		// AUD - audience
		// who the JWT is intended for.
		// Should be rejected if the principal processing
		// the claim does not identify itself with
		// the value listed here.
		AUD string

		// EXP - expiration time
		// the time the JWT ceases to be valid
		EXP string

		// NBT - OPTIONAL - not before time
		// the time the begins to be valid
		NBT string

		// IAT - OPTIONAL - issued at time
		// the time the JWT was issued
		IAT string

		// JTI - OPTIONAL - JWT ID
		// The unique identifier for the JWT
		JTI string

		// *** public claims ***
		// Public claims are collision resistant (i.e. URI namespaced)
		// or defined in the "IANA JSON Web Token Registry"
		// https://www.iana.org/assignments/jwt/jwt.xhtml

		// *** private claims ***
		// Custom claims specific to our web app.

	}
}

func New() *Token {

	return &Token{}
}

// Encode creates a token from the jwt struct
func (*Token) Encode() error {
	// check required fields have been completed
	// ISS Issuer - i.e. Server URL
	// SUB UUID or similar unique identifier
	// AUD Unique ID for recipient - i.e. app url
	// EXP The time the token expires
	// convert the Token to JSON
	// convert the JSON to URLbase64
	// Use the ALG to sign the token, add the signature to the end of the token
	return errors.New("Forgot to add code")
}

// Decode checks the validity of the jwt token,
// returns the unique identifer or an error
func (*Token) Decode() (UUID string, err error) {
	return "", errors.New("Forgot to add Code")
}
