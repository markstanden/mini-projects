package jwt

import (
	"fmt"
	"testing"
)

func TestCreateJWT(t *testing.T) {
	test := Token{
		Header:  Header{Algorithm: "HS512", TokenType: "JWT"},
		Payload: jwtioStruct,
	}

	jwt, err := test.Create(jwtioSecret())
	if err != nil {
		t.FailNow()
	}

	if jwt != jwtioToken {
		t.Fatalf("Failed to create the expected token: \nGot \n%v\nWanted \n%v\n", jwt, jwtioToken)
	}

	fmt.Println("Created jwt.io test token OK.")
}
