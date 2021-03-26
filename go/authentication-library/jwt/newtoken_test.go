package token

import (
	"fmt"
	"testing"
)

func TestNewToken(t *testing.T) {
	// create a JWT using data previously parsed on the jwt.io website
	j := jwtioStruct
	validFor := j.ExpirationTime - j.IssuedAtTime
	test := NewToken(j.Issuer, j.Audience, j.UserID, j.JwtID, j.KeyID, validFor)

	// The NewToken func will generate time and expiry dates based on the current time,
	// so check they are being set correctly then override to the supplied data used to
	// create the jwt.io jwt

	now := getUnixTime()
	if now-test.Payload.IssuedAtTime > 1 || now-test.Payload.IssuedAtTime < 0 {
		t.Errorf("issued at time not set correctly: Wanted %v\t Got %v\n", now, test.IssuedAtTime)
	}

	if now-test.Payload.NotBeforeTime > 1 || now-test.Payload.NotBeforeTime < 0 {
		t.Errorf("not before time not set correctly: Wanted %v\t Got %v\n", now, test.NotBeforeTime)
	}

	if now+validFor-test.Payload.ExpirationTime > 1 || now+validFor-test.Payload.ExpirationTime < 0 {
		t.Errorf("Expiration time not set correctly: Wanted %v\t Got %v\n", now, test.ExpirationTime)
	}

	test.IssuedAtTime = j.IssuedAtTime
	test.NotBeforeTime = j.NotBeforeTime
	test.ExpirationTime = j.ExpirationTime

	if test.Payload != j {
		t.Errorf("created struct is not as expected : \nWanted \n%v\n Got \n%v\n", j, test.Payload)
	}

	fmt.Println("Created jwt.io test Token struct OK.")
}
