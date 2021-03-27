package token

import (
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

	sigBS := hash(jwtString, secret)

	jwt = jwtString + "." + encodeBase64(sigBS)

	return jwt, nil
}