package jwt

import (
	"fmt"
	"log"

	"github.com/markstanden/authentication"
	"github.com/markstanden/securerandom"
	"github.com/markstanden/token"
)

type TokenService struct {
	Issuer     string
	Audience   string
	HoursValid int
	SecretKey  string
	Secrets    authentication.SecretService
}

func (ts *TokenService) passwordLookup() func(KeyID string) (secret string, err error) {
	//return ts.Secrets.GetSecret(ts.SecretKey)
	log.Println("tokenservice/jwt: password lookup called")
	return ts.Secrets.GetSecret(ts.SecretKey)
}

func (ts *TokenService) Create(u *authentication.User) (jwt string, err error) {

	// set the values for the token

	// The user the token refers to
	uniqueID := u.Token

	// The unique identifier for this token
	tokenID, err := securerandom.String(64)
	if err != nil {
		return "", fmt.Errorf("failed to create 'jti' :\n%v", err)
	}

	// the unique identifier for the secret
	// keyID, err := securerandom.String(64)
	// if err != nil {
	// 	return "", fmt.Errorf("failed to create 'kid' :\n%v", err)
	// }
	keyID := "1"

	// the number of seconds the token is valid for
	validFor := hoursToSeconds(ts.HoursValid)

	//create the token, and return
	t := token.NewToken(ts.Issuer, uniqueID, ts.Audience, tokenID, keyID, validFor)
	return t.CreateJWT(ts.passwordLookup())
}

func (ts *TokenService) Decode(jwt string) (userTokenID string, err error) {
	data := &token.Token{}
	err = token.Decode(jwt, ts.passwordLookup(), data)
	if err != nil {
		return "", fmt.Errorf("authentication/tokenhandler/token: Failed to decode JWT: \n%v", err)
	}
	if data.Audience != ts.Audience {
		return "", fmt.Errorf("authentication/tokenhandler/token: token is invalid - incorrect audience: \nWanted: %v, Got: %v", ts.Audience, data.Audience)
	}
	if data.Issuer != ts.Issuer {
		return "", fmt.Errorf("authentication/tokenhandler/token: token is invalid - incorrect issuer: \nWanted: %v, Got: %v", ts.Issuer, data.Issuer)
	}

	return data.Subject, nil
}

// returns hours in seconds
func hoursToSeconds(hours int) (secs int64) {
	return int64(hours * 60 * 60)
}
