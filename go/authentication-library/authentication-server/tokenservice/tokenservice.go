package tokenservice

import (
	"fmt"

	"github.com/markstanden/jwt"
	"github.com/markstanden/securerandom"
)

type TokenService struct {
	/*
		Issuer is the Issuer's ID/Name/URL.
		This gets set to the "iss" field within the payload of the JWT
	*/
	Issuer string

	/*
		Audience is the URL that the token is intended for.
		This gets set to the "aud" field within the payload of the JWT
		The JWT spec requires tokens to be rejected if the
		Audience does not match the expected value
	*/
	Audience string

	/*
		HoursValid is the number of hours until the JWT expires.
		This value is used in both creating and verifying tokens,
		and tokens *made* "iat" or *not before* "nbf" outside of this range will be rejected.
		Its primary role is to set the expiry date/time "exp" field on the JWT which defaults
		to the current time plus the value of HoursValid (in Hours)
	*/
	HoursValid int

	/*
		SecretCallback is the callback function to be invoked by the
		Create/Decode methods to obtain a specific version of the
		secret used to encode/decode the token.
		if err != nil the err should be returned, with empty data
	*/
	SecretCallback func(KeyID string) (secret string)
}

func (ts *TokenService) Create(userID string) (jwtString, jwtID string, err error) {

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
	t := jwt.NewToken(ts.Issuer, ts.Audience, userID, jwtID, keyID, validFor)

	jwtString, err = t.CreateJWT(ts.SecretCallback)
	if err != nil {
		return "", "", err
	}

	return jwtString, jwtID, nil

}

func (ts *TokenService) Decode(jwtString string) (userID, jwtID string, err error) {
	data := &jwt.Token{}
	err = jwt.Decode(jwtString, ts.SecretCallback, data)
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

// returns hours in seconds
func hoursToSeconds(hours int) (secs int64) {
	return int64(hours * 60 * 60)
}
