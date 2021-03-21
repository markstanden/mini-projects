package token

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"time"
)

// NewToken creates a new token, with sane defaults
func NewToken(secret, issuer, uniqueID, audience, tokenID, keyID string, validFor int64) (token string, err error) {

	// Get the current time and convert to UTC and standardised JSON string
	now := getUnixTime(time.Now())

	h := Header{
		Algorithm: "HS512",
		TokenType: "JWT",
	}
	p := Payload{
		Issuer:         issuer,
		Subject:        uniqueID,
		Audience:       audience,
		ExpirationTime: now + validFor,
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

func getUnixTime(t time.Time) int64 {
	return t.UTC().Unix()
}
