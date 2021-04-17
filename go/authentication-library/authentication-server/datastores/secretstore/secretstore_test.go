package secretstore_test

import (
	"testing"

	"github.com/markstanden/authentication"
	"github.com/markstanden/authentication/datastores/postgres"
	"github.com/markstanden/authentication/datastores/secretstore"
)

func GetTestSecretService(testDB postgres.DataStore) (testSS authentication.SecretDataStore) {
	testSS = secretstore.New(testDB, 5 /* seconds */)
	// test.US.XXX = XXX
	return testSS
}

func TestSecretService(t *testing.T) {
	testCases := []struct {
		desc string
	}{
		{
			desc: "",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

		})
	}
}
