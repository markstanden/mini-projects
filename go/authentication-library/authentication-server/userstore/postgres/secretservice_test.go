package postgres_test

import (
	"testing"

	"github.com/markstanden/authentication/userstore/postgres"
)

func GetTestSecretService(testDB postgres.DataStore) (testUS postgres.SecretService) {
	testUS = postgres.NewSecretService(testDB, 5 /* seconds */)
	// test.US.XXX = XXX
	return testUS
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
