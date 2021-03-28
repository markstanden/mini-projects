package jwt

import (
	"fmt"
	"testing"
)

// test with a valid SHA512 JWT created from jwt.io website
func TestDecodeJWTIO(t *testing.T) {

	//create an empty struct
	got := Token{
		lifespan: 100000000,
	}

	secret := jwtioSecret()

	// Decode the test data
	err := Decode(jwtioToken, secret, &got)
	if err != nil {
		fmt.Printf("failed to decode JWT \n%t", err)
	}
	if got.Payload != jwtioStruct {
		t.Errorf("Want: \n%v\nGot: \n%v\n", jwtioStruct, got.Payload)
	}

	fmt.Println("Decoded jwt.io test token OK.")
}
