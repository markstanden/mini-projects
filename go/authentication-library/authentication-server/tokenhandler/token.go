package token

import (
	"fmt"

	"github.com/markstanden/authentication"
	"github.com/markstanden/authentication/random"
	"github.com/markstanden/token"
)

func Create(u *authentication.User, secret string) (jwt string, err error) {

	// set the values for the token
	issuer := "markstanden.dev"
	uniqueID := u.Token
	audience := "markstanden.dev"
	tokenID, err := random.String(64)
	if err != nil {
		return "", fmt.Errorf("failed to create 'jti' :\n%v", err)
	}
	keyID, err := random.String(64)
	if err != nil {
		return "", fmt.Errorf("failed to create 'kid' :\n%v", err)
	}
	validFor := hoursValid(24)

	//create the token, and return
	return token.Create(secret, issuer, uniqueID, audience, tokenID, keyID, validFor)
}

func Decode(jwt, secret string) (map[string]interface{}, error) {
	return token.Decode(jwt, secret)
}

// returns hours in seconds
func hoursValid(hours int) (secs int64) {
	return int64(hours * 60 * 60)
}
