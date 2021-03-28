package jwt

import (
	"encoding/json"

	"github.com/markstanden/jwt/b64"
	"github.com/markstanden/jwt/hash"
)

func (t *Token) CreateJWT(getRemoteSecret func(keyID string) string) (jwt string, err error) {

	jsonHeader, err := json.Marshal(t.Header)
	if err != nil {
		return "", err
	}
	jsonPayload, err := json.Marshal(t.Payload)
	if err != nil {
		return "", err
	}

	jwtBody := b64.FromBytes(jsonHeader) + "." + b64.FromBytes(jsonPayload)

	secret := getRemoteSecret(t.KeyID)
	if secret == "" {
		return "", ErrFailedSecret
	}

	sigBS := hash.HS512(jwtBody, secret)

	jwt = jwtBody + "." + b64.FromBytes(sigBS)

	return jwt, nil
}
