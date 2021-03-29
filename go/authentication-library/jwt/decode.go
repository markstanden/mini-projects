package jwt

import (
	"crypto/hmac"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/markstanden/jwt/b64"
	"github.com/markstanden/jwt/hash"
	"github.com/markstanden/jwt/time"
)

func addError(err error, new string) error {
	return errors.New(err.Error() + "\n" + new)
}

/*
	Decode takes an untrusted JWT and checks it for validity:
		- Checks structure, and for invalid characters
		- Checks Header
		- Checks Payload, checking timestamps are valid
		- Checks signature, by calling the secret callback with the key version encoded within the JWT
	it returns:
		- nil error if the token is valid and has not expired
		- ErrInvalidToken if the token fails any of the validity checks
		- ErrExpiredToken if the token is valid, but has expired
		- ErrFailedSecret if the callback failed to return a secret, or the secret was an empty string
*/
func Decode(untrustedJWT string, passwordLookup func(key string) (secret string), token *Token) (err error) {

	err = errors.New("decoding JWT")

	/*
		ValidFrom is the official time that the server started issuing tokens.
		Any tokens with an time stamp before this will be discarded as invalid, since
		we cannot have issued them before this date.
		If the value hasn't been set in the supplied struct we will default it to the start of 2021.
	*/
	if token.Config.ValidFrom == 0 {
		/* 01 Jan 2021 00:00 UTC */
		token.Config.ValidFrom = 1609459200
	}

	ut := Token{}

	/*
		Check the JWT as a whole
		uthe JWT should be three base64 URL encoded strings,
		separated by full stops.
	*/
	jwtSection := strings.Split(untrustedJWT, ".")

	if !checkJwtValid(jwtSection) {
		return addError(err, "Failed checkJwtValid")
	}

	/*
		assign each of the sections to
		variables to increase readability
	*/
	header := jwtSection[0]
	payload := jwtSection[1]
	signature := jwtSection[2]
	jwtSection = nil

	/*
		Check header section
	*/

	if err := unmarshalJWT(header, &ut); err != nil {
		return addError(err, "Failed to unmarshalJWT header")
	}

	if !checkHeaderValid(ut.Header) {
		return addError(err, "Failed checkHeaderValid")
	}

	/*
		Check payload section
	*/

	if err := unmarshalJWT(payload, &ut); err != nil {
		return addError(err, "Failed to unmarshalJWT payload")
	}

	tokenInvalid, tokenExpired := checkTimeValidity(
		ut.IssuedAtTime,
		ut.NotBeforeTime,
		ut.ExpirationTime,
		token.ValidFrom,
		token.lifespan)

	if tokenInvalid {
		return addError(err, "Failed checkTimeValidity")
	}

	/*
		Check Signature
	*/
	fmt.Println(passwordLookup(ut.KeyID))
	if err := signatureValid(header, payload, signature, passwordLookup(ut.KeyID)); err != nil {
		return err
	}

	/*
		Copy the data from the untrusted struct
		into the supplied pointer to a struct
		as the data has now been validated
	*/
	*token = ut

	/*
		Token is valid but expired, so return data with expiry error
	*/
	if tokenExpired {
		return ErrExpiredToken
	}

	return nil
}

/*
	checkHeaderValid performs tests on the contents of the JWT header
	returns true only if all tests pass
*/
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

/*
	checkJwtValid checks the basic make up of each part of the JWT
	The idea is these are cheap checks, compared with >300ms for a hash check,
	returns true only if all tests pass
*/
func checkJwtValid(jwt []string) bool {

	/*
		The allowable characters in a URL flavoured base64 encoded string.
		shamelessly stolen from the base64 package, as it is not exported.
	*/
	const encodeURL = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"

	/*
		Our JWTs will contain three sections,
		a Header, Payload, and Signature
	*/
	if len(jwt) != 3 {
		return false
	}

	for i := range jwt {

		/* Check for invalid characters in the input string */
		if !checkRunes(jwt[i], encodeURL) {
			return false
		}

		/* Check for empty sections */
		if len(jwt[i]) == 0 {
			return false
		}
	}

	return true
}

/*
	checkValidity checks whether the supplied string contains
	runes outside of the valid runes set supplied
*/
func checkRunes(toCheck, validRunes string) (valid bool) {
	for _, r := range toCheck {
		valid = strings.ContainsRune(validRunes, r)
		if !valid {
			/* invalid rune found */
			return false
		}
	}
	return true
}

/*
	unmarshalJWT decodes a section of the JWT and
	unmarshals the JSON data into the provided *Token
*/
func unmarshalJWT(jwtSection string, t *Token) error {

	bytes, err := b64.ToBytes(jwtSection)
	if err != nil {
		return ErrInvalidToken
	}

	if err := json.Unmarshal(bytes, t); err != nil {
		return ErrInvalidToken
	}
	fmt.Println(t.KeyID)
	return nil
}

/*
	checkTimeValidity checks that the Issued at time, Not before time and Expiry
	are set to values that could have been set by our server.
	and that they are within the expiry window
*/
func checkTimeValidity(iat, nbf, exp, firstIssuedToken, lifespan int64) (tokenInvalid, tokenExpired bool) {

	now := time.GetUnix()
	min := now - lifespan
	max := iat + lifespan

	fmt.Println("iat", iat)
	fmt.Println("nbf", nbf)
	fmt.Println("exp", exp)
	fmt.Println("now", now)
	fmt.Println("min", min)
	fmt.Println("max", max)

	/*
		the IssuedAtTime must be verified first as the value is used
		to verify the expiration time of the token.
		if the issued at time is in the future, or before the project started the
		time stamp is invalid.
		if the time so far in the past that the token would have expired anyway, but not
		before the project started the token should be marked valid but expired
	*/
	fmt.Println("tokenExpired:", tokenExpired, "tokenInvalid", tokenInvalid)
	if time.WithinRange(iat, firstIssuedToken, now) {
		if iat < min {
			/*
				token could have been made by our server,
				but was too long ago to not have expired
			*/
			tokenExpired = true
		}
	} else {
		/*
			token was made before the project began,
			or in the future, so must be invalid
		*/
		tokenInvalid = true
	}
	fmt.Println("tokenExpired:", tokenExpired, "tokenInvalid", tokenInvalid)
	/*
		The not before time restricts access until a point in time has been reached,
		so if that time has not been reached yet, the token is invalid.
	*/
	if time.WithinRange(nbf, firstIssuedToken, max) {
		/* It could have been issued by our server */
		if nbf > now {
			/* Token is not yet valid */
			tokenInvalid = true
		}
	} else {
		/*
			token was made before the project began,
			or in the future, so must be invalid
		*/
		tokenInvalid = true
	}
	fmt.Println("tokenExpired:", tokenExpired, "tokenInvalid", tokenInvalid)
	/*
		First check that the token could have been made by the server,
		and that the expiry date is not too far in the future
	*/
	if time.WithinRange(exp, firstIssuedToken, max) {
		/* It could have been issued by our server */
		if exp < now {
			/* Token has expired */
			tokenExpired = true
		}
	} else {
		/*
			token was made before the project began,
			or too far in the future, so must be invalid
		*/
		tokenInvalid = true
	}
	fmt.Println("tokenExpired:", tokenExpired, "tokenInvalid", tokenInvalid)
	return tokenInvalid, tokenExpired
}

func signatureValid(header, payload, signature, secret string) (err error) {

	/*
		If the secret is empty we have failed to obtain the secret from the secret store,
		or the secret is intentionally empty.
	*/
	if secret == "" {
		return ErrFailedSecret
	}

	bodyBytes := header + "." + payload

	signatureBytes, err := b64.ToBytes(signature)
	if err != nil {
		return ErrInvalidToken
	}

	/*
		Re-Hash the body section, using our secret,
		and compare to the signature supplied
		by the JWT using a time safe comparison.
	*/
	hashedBody := hash.HS512(bodyBytes, secret)
	if hmac.Equal(hashedBody, signatureBytes) {
		/* signature is valid */
		return nil
	}

	/* signature is invalid */
	return ErrInvalidToken
}
