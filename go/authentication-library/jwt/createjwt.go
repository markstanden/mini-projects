package token

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

func (t *Token) CreateJWT(passwordLookup func(keyID string) (string, error)) (jwt string, err error) {

	jsonHeader, err := json.Marshal(t.Header)
	if err != nil {
		return "", err
	}
	jsonPayload, err := json.Marshal(t.Payload)
	if err != nil {
		return "", err
	}

	jwtString := base64.RawURLEncoding.EncodeToString(jsonHeader) + "." + base64.RawURLEncoding.EncodeToString(jsonPayload)

	secret, err := passwordLookup(t.KeyID)
	if err != nil {
		return "", fmt.Errorf("failed to obtain secret from callback")
	}

	sig64 := hash(jwtString, secret)

	jwt = jwtString + "." + sig64

	return jwt, nil

}

// hash uses HMAC Sha512 to hash the provided message using
// the provided secret.
// the hash is returned as a URL encoded base64 string
func hash(message, secret string) (hash string) {
	hmac := hmac.New(sha512.New, []byte(secret))
	hmac.Write([]byte(message))
	bs := hmac.Sum(nil)
	return base64.RawURLEncoding.EncodeToString(bs)
}
