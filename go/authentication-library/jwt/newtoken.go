package jwt

import "github.com/markstanden/jwt/time"

/*
	NewToken creates a new jwt token struct, with sane defaults for header and payload time values.
	issuer (string) is assigned to the Token.Issuer and represents the identity of the creating site/service.
	audience (string) is assigned to the Token.Audience and represents the identity of the site/service the token is being created for
	userID (string) is assigned to Token.UserID and is a unique identifier to represent the user the token is to authenticate.
	jwtID (string) is a unique key to identify this particular jwt.
	keyID (string) is the key, or version of the secret that the signature of the JWT is to be encoded with.
	validFor (int64)
*/
func NewToken(issuer, audience, userID, jwtID, keyID string, validFor int64) (token *Token) {

	/*
		If a negative value of expiration time is provided
		set to zero to avoid potential strange behaviour
	*/
	if validFor < 0 {
		validFor = 0
	}

	h := Header{
		Algorithm: "HS512",
		TokenType: "JWT",
	}
	p := Payload{
		Issuer:         issuer,
		Audience:       audience,
		UserID:         userID,
		JwtID:          jwtID,
		KeyID:          keyID,
		IssuedAtTime:   time.GetUnix(),
		NotBeforeTime:  time.GetUnix(),
		ExpirationTime: time.GetUnix() + validFor,
	}
	c := Config{
		Lifespan: validFor,
	}
	return &Token{h, p, c}
}
