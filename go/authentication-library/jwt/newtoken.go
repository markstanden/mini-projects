package token

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
)

// NewToken creates a new token, with sane defaults
func NewToken(secret, issuer, uniqueID, audience, tokenID, keyID string) (token string, err error) {

	// Get the current time and convert to UTC and standardised JSON string
	now, err := format(getTime())
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
		KeyID:          keyID,
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
