package jwt

import (
	"crypto/hmac"
	"encoding/json"
	"log"
	"strings"

	"github.com/markstanden/jwt/b64"
	"github.com/markstanden/jwt/hash"
	"github.com/markstanden/jwt/time"
)

// Decode turns a signed JWT into a map[string]interface (or returns an error)
// but only after checking the validity of the jwt.
func Decode(untrustedJWT string, passwordLookup func(key string) (secret string), token *Token) (err error) {

	ut := Token{}

	// untrustedJWT should be three base64 URL encoded strings, separated by full stops.
	jwtSection := strings.Split(untrustedJWT, ".")

	if !checkJwtValid(jwtSection) {
		return ErrInvalidToken
	}

	header := jwtSection[0]
	payload := jwtSection[1]

	signature, err := b64.ToBytes(jwtSection[2])
	if err != nil {
		return ErrInvalidToken
	}

	// Check header section
	if err := unmarshalJWT(header, &ut); err != nil {
		return ErrInvalidToken
	}

	if !checkHeaderValid(ut.Header) {
		return ErrInvalidToken
	}

	// Check payload section
	if err := unmarshalJWT(payload, &ut); err != nil {
		return ErrInvalidToken
	}

	now := time.GetUnix()
	min := now - token.lifespan
	max := ut.IssuedAtTime + token.lifespan

	tokenInvalid := false
	tokenExpired := false

	/*
		if the expiry date is not between now
		and point of expiry of the token it has either expired or
		it is invalid.
	*/
	if !time.WithinRange(ut.ExpirationTime, now, max){
		tokenExpired = tokenExpired || ut.ExpirationTime < now
		tokenInvalid = tokenInvalid || ut.ExpirationTime > max
	}

	/*
		if the issued at time is in the future,
		or so far in the past that the token would have expired anyway
	*/
	!time.WithinRange(ut.IssuedAtTime, min, now)
		/*
			if the 
		*/
		!time.WithinRange(ut.NotBeforeTime, min, now) 	{
			return ErrExpiredToken
	}
}

	secret := passwordLookup(ut.KeyID)
	if secret == "" {
		return ErrFailedSecret
	}

	jwtBody := header + "." + payload

	testBytes := hash.HS512(jwtBody, secret)
	// check that the hashed body is equal to the decoded signature supplied by the jwt
	if !hmac.Equal(testBytes, signature) {
		return ErrInvalidToken
	}

	// Copy the data from the untrusted struct into the supplied struct
	*token = ut

	return nil
}

// checkHeaderValid performs tests on the contents of the JWT header
// returns true only if all tests pass
func checkHeaderValid(h Header) bool {

	// jwt vulnerability where the signature can be
	// bypassed by setting the alg to none.
	// if this is attempted log it
	if h.Algorithm == "none" {
		log.Println(`"alg": "none" present in header`)
		return false
	}

	if h.Algorithm != "HS512" {
		return false
	}

	if h.TokenType != "JWT" {
		return false
	}

	return true
}

// checkJwtValid checks the basic make up of each part of the JWT
// The idea is these are cheap checks, compared with >300ms for a hash check,
// returns true only if all tests pass
func checkJwtValid(jwt []string) bool {

	// The allowable characters in a URL flavoured base64 encoded string.
	// shamelessly stolen from the base64 package, as it is not exported.
	const encodeURL = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"

	// Our JWTs will contain three sections, a Header, Payload, and Signature
	if len(jwt) != 3 {
		return false
	}

	for i := range jwt {

		// Check for invalid characters in the input string
		if !checkRunes(jwt[i], encodeURL) {
			return false
		}

		// Check for empty sections
		if len(jwt[i]) == 0 {
			return false
		}
	}

	return true
}

// checkValidity checks whether the supplied string contains runes outside of the valid runes set supplied
func checkRunes(toCheck, validRunes string) (valid bool) {
	for _, r := range toCheck {
		valid = strings.ContainsRune(validRunes, r)
		if !valid {
			// invalid rune found
			return false
		}
	}
	return true
}

// unmarshalJWT decodes a section of the JWT and
// unmarshals the JSON data into the provided *Token
func unmarshalJWT(jwtSection string, t *Token) error {

	bytes, err := b64.ToBytes(jwtSection)
	if err != nil {
		return ErrInvalidToken
	}

	if err := json.Unmarshal(bytes, t); err != nil {
		return ErrInvalidToken
	}

	return nil
}
