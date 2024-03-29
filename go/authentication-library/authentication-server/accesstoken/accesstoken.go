package accesstoken

import (
	"errors"
	"fmt"

	"github.com/markstanden/authentication"
	"github.com/markstanden/jwt"
	"github.com/markstanden/securerandom"
)

/*
	** AccessToken **
	This struct holds the config for the creation, verification, and decoding
	of the short lived access tokens used for stateless authentication
*/
type AccessToken struct {

	/*
		Audience is the URL that the token is intended for.
		This gets set to the "aud" field within the payload of the JWT
		The JWT spec requires tokens to be rejected if the
		Audience does not match the expected value
	*/
	Audience string

	/*
		Issuer is the Issuer's ID/Name/URL.
		This gets set to the "iss" field within the payload of the JWT
	*/
	Issuer string

	/*
		AuthLevel is the authentication level of the user.
	*/
	AuthLevel string

	/*
		MinsValid is the number of mins until the JWT expires.
		This value is used in both creating and verifying tokens,
		and tokens *made* "iat" or *not before* "nbf" outside of this range will be rejected.
		Its primary role is to set the expiry date/time "exp" field on the JWT which defaults
		to the current time plus the value of MinsValid (in mins)
	*/
	MinsValid int

	/*
		SecretCallback is the callback function to be invoked by the
		Create/Decode methods to obtain a specific version of the
		secret used to encode/decode the token.
		if err != nil the secret should be returned empty
	*/
	Secret authentication.SecretDataStore

	/*
		StartTime is time the server started issuing tokens in Unix time UTC
		It is used as the earliest possible time a token is valid.  Tokens made
		before this will time will be automatically discarded
	*/
	StartTime int64
}

func New() (at AccessToken) {
	return at
}

/*
	** Create **
	Create takes the provided (anonymous) userID and creates a new JWT
	returning the new JWT, and the unique identifer for the token itself.
*/
func (ts *AccessToken) Create(userID string) (jwtString, jwtID string, err error) {

	// Check for invalid input
	if userID == "" {
		return "", "", errors.New("invalid userID")
	}

	// The unique identifier for this particular token
	jwtID = securerandom.String(64)
	if jwtID == "" {
		return "", "", fmt.Errorf("failed to create 'jti' :\n%v", err)
	}

	// the unique identifier for the secret
	keyID := ts.Secret.GetKeyID("JWT")

	// the number of seconds the token is valid for
	validFor := minsToSeconds(ts.MinsValid)

	//create the token, and return
	t := jwt.NewToken(ts.Issuer, ts.Audience, userID, jwtID, keyID, validFor)

	jwtString, err = t.Create(ts.Secret.GetSecret("JWT"))
	if err != nil {
		return "", "", err
	}

	return jwtString, jwtID, nil

}

/*
	** Decode **
	Decode takes a JWT string, verifies it and returns the embedded UserID and JwtID.
	The UserID is an (anonymous) identifier for the issued to user,
	and the JwtID is the identifier for the token itself.
*/
func (ts *AccessToken) Decode(jwtString string) (tokenUserID, jwtID string, err error) {

	data := new(jwt.Token)
	data.Config.Lifespan = minsToSeconds(ts.MinsValid)

	err = jwt.Decode(jwtString, ts.Secret.GetSecret("JWT"), data)
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

/*
	** minsToSeconds **
	minsToSeconds returns the minutes valid (int) in seconds (int64)
	so it can be easily used as a unix time measurement
*/
func minsToSeconds(mins int) (secs int64) {
	return int64(mins * 60)
}
