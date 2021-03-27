package token

import (
	"encoding/base64"
	"encoding/json"
)

func (t *Token) CreateJWT(passwordLookup func(keyID string) string) (jwt string, err error) {

	jsonHeader, err := json.Marshal(t.Header)
	if err != nil {
		return "", err
	}
	jsonPayload, err := json.Marshal(t.Payload)
	if err != nil {
		return "", err
	}

	jwtString := base64.RawURLEncoding.EncodeToString(jsonHeader) + "." + base64.RawURLEncoding.EncodeToString(jsonPayload)

	secret := passwordLookup(t.KeyID)
	if secret == "" {
		return "", ErrFailedSecret
	}

	sigBS := hash(jwtString, secret)

	jwt = jwtString + "." + encodeBase64(sigBS)

	return jwt, nil
}
