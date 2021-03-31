package postgres_test

import (
	"testing"

	"github.com/markstanden/authentication/userstore/postgres"
)

/*
	Wraps a connection to the test DB in a userservice to allow testing on a real store.
*/
func GetTestUserService(testDB postgres.DataStore) (testUS postgres.UserService) {
	testUS = postgres.NewUserService(testDB)
	// test.US.XXX = XXX
	return testUS
}

func Test(t *testing.T) {
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
