package token

// NewToken creates a new token, with sane defaults for header and payload time values,
func NewToken(issuer, audience, userID, jwtID, keyID string, validFor int64) (token *Token) {

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
		IssuedAtTime:   getUnixTime(),
		NotBeforeTime:  getUnixTime(),
		ExpirationTime: getUnixTime() + validFor,
	}
	return &Token{h, p, validFor}
}
