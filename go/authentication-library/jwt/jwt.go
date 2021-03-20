package jwt

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

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

func getTime() time.Time {
	return time.Now()
}

func format(t time.Time) (string, error) {
	sb, err := t.UTC().MarshalText()
	if err != nil {
		return "", err
	}
	return string(sb), nil
}

// NewToken creates a new token, with sane defaults
func NewToken(secret, issuer, uniqueID, audience, tokenID, keyID string) (token string, err error) {

	// Get the current time and convert to UTC and standardised JSON string
	now, err := format(time.Now())
	if err != nil {
		return "", fmt.Errorf("incorrect time set: \n%v", err)
	}
	expires, err := format(
		time.Now().
			AddDate(
				0, /* years */
				1, /* months */
				0 /* days */))
	if err != nil {
		return "", fmt.Errorf("incorrect expiry set: \n%v", err)
	}

	h := Header{
		Algorithm: "HS512",
		TokenType: "JWT",
	}
	p := Payload{
		Issuer:         issuer,
		Subject:        uniqueID,
		Audience:       audience,
		ExpirationTime: expires,
		NotBeforeTime:  now,
		IssuedAtTime:   now,
		TokenID:        tokenID,
		KeyID:			keyID,
	}
	jsonHeader, err := json.Marshal(h)
	if err != nil {
		return "", err
	}
	jsonPayload, err := json.Marshal(p)
	if err != nil {
		return "", err
	}

	jwtString := base64.RawURLEncoding.EncodeToString(jsonHeader) + "." + base64.RawURLEncoding.EncodeToString(jsonPayload)

	hmac := hmac.New(sha512.New, []byte(secret))
	hmac.Write([]byte(jwtString))
	signature := hmac.Sum(nil)

	signature64 := base64.RawURLEncoding.EncodeToString(signature)
	signedToken := jwtString + "." + signature64

	return signedToken, nil

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
	return errors.New("forgot to add code")
}

// Decode checks the validity of the jwt token,
// returns the unique identifer or an error
func Validate(unverified string, secret string) (data map[string]string, err error) {
	// from jwt spec

	// Verify that the JWT contains at least one period ('.') character.
	if !strings.Contains(unverified, ".") {
		return nil, fmt.Errorf("invalid token")
	}

	// Split the dot separated base64 URL encoded string into 2 or 3 segments depending on whether it
	// contains a signature
	splitUnv := strings.SplitAfter(unverified, ".")

	if len(splitUnv) == 0 || len(splitUnv) > 3 {
		return nil, fmt.Errorf("token split error, incorrect number of sections")
	}
	header64 := strings.TrimSuffix(splitUnv[0], ".")
	payload64 := strings.TrimSuffix(splitUnv[1], ".")
	signature64 := strings.TrimSuffix(splitUnv[2], ".")

	//(HMACSHA512(base64UrlEncode(header) + "." + base64UrlEncode(payload), secret)

	headerBytes, err := base64.RawURLEncoding.DecodeString(header64)
	if err != nil {
		return nil, err
	}
	payloadBytes, err := base64.RawURLEncoding.DecodeString(payload64)
	if err != nil {
		return nil, err
	}
	signatureBytes, err := base64.RawURLEncoding.DecodeString(signature64)
	if err != nil {
		return nil, err
	}

	// create the slice of bytes to encode
	// {json header} + "." + {json payload}
	var toEncode []byte
	toEncode = append(toEncode, header64...)
	toEncode = append(toEncode, []byte(".")...)
	toEncode = append(toEncode, payload64...)

	// this will take a while to compute
	h := hmac.New(sha512.New, []byte(secret))
	h.Write(toEncode)
	testBytes := h.Sum(nil)

	if hmac.Equal(testBytes, signatureBytes) {
		fmt.Println("Signature Verified")
	} else {
		fmt.Println("Signature invalid")
	}

	fmt.Println(string(headerBytes))
	fmt.Println(string(payloadBytes))
	// Let the Encoded JOSE Header be the portion of the JWT before the first period ('.') character.
	return data, nil
}

/*
7.2.  Validating a JWT

   When validating a JWT, the following steps are performed.  The order
   of the steps is not significant in cases where there are no
   dependencies between the inputs and outputs of the steps.  If any of
   the listed steps fail, then the JWT MUST be rejected -- that is,
   treated by the application as an invalid input.

   1.   Verify that the JWT contains at least one period ('.')
        character.

   2.   Let the Encoded JOSE Header be the portion of the JWT before the
        first period ('.') character.

   3.   Base64url decode the Encoded JOSE Header following the
        restriction that no line breaks, whitespace, or other additional
        characters have been used.

   4.   Verify that the resulting octet sequence is a UTF-8-encoded
        representation of a completely valid JSON object conforming to
        RFC 7159 [RFC7159]; let the JOSE Header be this JSON object.

   5.   Verify that the resulting JOSE Header includes only parameters
        and values whose syntax and semantics are both understood and
        supported or that are specified as being ignored when not
        understood.

   6.   Determine whether the JWT is a JWS or a JWE using any of the
        methods described in Section 9 of [JWE].

	7.   Depending upon whether the JWT is a JWS or JWE, there are two
        cases:

        *  If the JWT is a JWS, follow the steps specified in [JWS] for
           validating a JWS.  Let the Message be the result of base64url
           decoding the JWS Payload.

        *  Else, if the JWT is a JWE, follow the steps specified in
           [JWE] for validating a JWE.  Let the Message be the resulting
           plaintext.

   8.   If the JOSE Header contains a "cty" (content type) value of
        "JWT", then the Message is a JWT that was the subject of nested
        signing or encryption operations.  In this case, return to Step
        1, using the Message as the JWT.

   9.   Otherwise, base64url decode the Message following the
        restriction that no line breaks, whitespace, or other additional
        characters have been used.

   10.  Verify that the resulting octet sequence is a UTF-8-encoded
        representation of a completely valid JSON object conforming to
        RFC 7159 [RFC7159]; let the JWT Claims Set be this JSON object.

*/
