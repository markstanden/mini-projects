package token

import (
	"time"
)

// NewToken creates a new token, with sane defaults for header and payload time values,
func NewToken(issuer, uniqueID, audience, tokenID, keyID string, validFor int64) (token *Token) {

	h := Header{
		Algorithm: "HS512",
		TokenType: "JWT",
	}
	p := Payload{
		Issuer:         issuer,
		Subject:        uniqueID,
		Audience:       audience,
		ExpirationTime: getUnixTime() + validFor,
		NotBeforeTime:  getUnixTime(),
		IssuedAtTime:   getUnixTime(),
		TokenID:        tokenID,
		KeyID:          keyID,
	}
	return &Token{h, p}
}

func getUnixTime() int64 {
	return time.Now().UTC().Unix()
}
