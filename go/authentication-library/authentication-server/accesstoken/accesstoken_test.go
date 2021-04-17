package accesstoken

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/markstanden/authentication/datastores/postgres"
	"github.com/markstanden/authentication/datastores/secretstore"
	"github.com/markstanden/securerandom"
)

/*
	Tests for the creation of a valid JWT
*/
func TestCreate(t *testing.T) {

	ds, err := postgres.GetTestConfig().FromEnv().Connect()
	secretStore := secretstore.New(ds, int64(60))
	secretStore.FullReset()

	tokenLength := uint(64)

	if err != nil {
		t.Fatal("failed to connect to test database")
	}
	testCases := []struct {
		desc        string
		tokenUserID string
		isValid     bool
	}{
		{
			desc:        "Empty ID",
			tokenUserID: "",
			isValid:     false,
		},
		{
			desc:        "Valid TokenUserID",
			tokenUserID: securerandom.String(tokenLength),
			isValid:     true,
		},
	}

	/* Create the token handler */
	ts := New()
	ts.Audience = "Test"
	ts.Issuer = "Test Issuer"
	ts.MinsValid = 10
	ts.Secret = secretStore
	ts.StartTime = time.Now().UTC().Unix() - int64(time.Minute)

	/* The string to be created in step one, and decoded in step 2 */
	var jwtString, jwtID string

	for _, test := range testCases {
		t.Run("Create", func(t *testing.T) {
			t.Run(test.desc, func(t *testing.T) {

				jwtString, jwtID, err = ts.Create(test.tokenUserID)

				if err == nil {
					log.Printf("\nAccess Token created successfully for UserID (%v).  \nTokenID: %v\nJWT: %v\nErr: %v", test.tokenUserID, jwtID, jwtString, err)
				}
				if test.isValid && err != nil {
					t.Fatal("failed to create token:\n", err)
				}

				if !test.isValid && err == nil {
					t.Fatal("failed to error with invalid inputs")
				}

			})

		})

		t.Run("Decode", func(t *testing.T) {
			decodedUserID, decodedJWTID, err := ts.Decode(jwtString)
			if test.isValid && err != nil {
				t.Fatal("failed to decode created JWT")
			}
			if !test.isValid && err == nil {
				t.Fatal("managed to return a token for an invalid test!??!")
			}
			if test.isValid && decodedJWTID != jwtID {
				t.Fatal("failed to obtain correct jwtID")
			}
			if test.isValid && decodedUserID != test.tokenUserID {
				t.Fatal("failed to decode correct UserID from token")
			}
			fmt.Println("Token Decoded OK, UserID :", decodedUserID)
		})
	}
}
