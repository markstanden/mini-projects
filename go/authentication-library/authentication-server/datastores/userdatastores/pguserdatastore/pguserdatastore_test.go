package pguserdatastore_test

import (
	"testing"

	"github.com/markstanden/authentication/datastores/postgres"
	"github.com/markstanden/authentication/datastores/userdatastores/pguserdatastore"
)

/*
	Wraps a connection to the test DB in a userservice to allow testing on a real store.
*/
func GetTestUserService(testDB postgres.DataStore) (testUS pguserdatastore.PGUserDataStore) {
	testUS = pguserdatastore.NewUserService(testDB)
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
