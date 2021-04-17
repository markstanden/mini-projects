package userstore

import (
	"github.com/markstanden/authentication/datastores/postgres"
)

/*
	Wraps a connection to the test DB in a userservice to allow testing on a real store.
*/
func GetTestUserService(testDB postgres.DataStore) (testUS Userstore) {
	testUS = New(testDB)
	// test.US.XXX = XXX
	return testUS
}
