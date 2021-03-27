package jwt

import (
	"fmt"

	"github.com/markstanden/securerandom"
	"github.com/markstanden/token"
)

type TokenService struct {
	// The Issuer's URL
	Issuer string

	// The URL that the token is intended for
	Audience string

	// The number of hours until the JWT expires
	HoursValid int

	// Callback to be invoked by the Create/Decode methods to obtain a specific version of the
	// secret used to encode/decode the token.
	// if err != nil the err should be returned, with empty data
	SecretCallback func(KeyID string) (secret string, err error)
}

func (ts *TokenService) Create(userID string) (jwt, jwtID string, err error) {

	// The unique identifier for this particular token
	jwtID = securerandom.String(64)
	if jwtID == "" {
		return "", "", fmt.Errorf("failed to create 'jti' :\n%v", err)
	}

	// the unique identifier for the secret
	// keyID, err := securerandom.String(64)
	// if err != nil {
	// 	return "", fmt.Errorf("failed to create 'kid' :\n%v", err)
	// }
	// Temp override until keystore is implemented
	keyID := "1"

	// the number of seconds the token is valid for
	validFor := hoursToSeconds(ts.HoursValid)

	//create the token, and return
	t := token.NewToken(ts.Issuer, ts.Audience, userID, jwtID, keyID, validFor)

	jwt, err = t.CreateJWT(ts.GetSecret)
	if err != nil {
		return "", "", err
	}

	return jwt, jwtID, nil

}

func (ts *TokenService) Decode(jwt string) (userID, jwtID string, err error) {
	data := &token.Token{}
	err = token.Decode(jwt, ts.GetSecret, data)
	if err != nil {
		return "", "", fmt.Errorf("authentication/tokenhandler/token: Failed to decode JWT: \n%v", err)
	}
	if data.Audience != ts.Audience {
		return "", "", fmt.Errorf("authentication/tokenhandler/token: token is invalid - incorrect audience: \nWanted: %v, Got: %v", ts.Audience, data.Audience)
	}
	if data.Issuer != ts.Issuer {
		return "", "", fmt.Errorf("authentication/tokenhandler/token: token is invalid - incorrect issuer: \nWanted: %v, Got: %v", ts.Issuer, data.Issuer)
	}

	return data.UserID, data.JwtID, nil
}

func (ts *TokenService) GetSecret(KeyID string) (secret string, err error) {
	return ts.SecretCallback(KeyID)
}

// returns hours in seconds
func hoursToSeconds(hours int) (secs int64) {
	return int64(hours * 60 * 60)
}
