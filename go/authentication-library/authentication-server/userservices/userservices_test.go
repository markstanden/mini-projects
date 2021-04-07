package userservice_test

import (
	"strings"
	"testing"

	"github.com/markstanden/authentication/datastores/postgres"
	"github.com/markstanden/authentication/datastores/userdatastores/pguserdatastore"
	userservice "github.com/markstanden/authentication/userservices"
)

type test struct {
	name     string
	email    string
	password string
}

/*
	*** TestNewUser ***
	TestNewUser creates a new user from the test slice,
	and checks that the user has been created and added
	to the test datastore.
*/
func TestNewUser(t *testing.T) {
	test := test{
		name:     "Testy McTestFace",
		email:    "testy@mctestface.com",
		password: "livetotest",
	}
	/*
		connect to the test database
	*/
	ds, err := postgres.GetTestConfig().FromEnv().Connect()
	if err != nil {
		t.Fatal("failed to connect to the temp database")
	}

	/*
		create the UserService using the test DB
	*/
	us := userservice.UserService{
		UserDS: pguserdatastore.NewUserService(ds),
		Config: userservice.USConfig{TokenIDSize: 32},
	}

	/*
		Add the user to the database
	*/
	newUser, err := us.NewUser(test.name, test.email, test.password)
	if err != nil {
		t.Error("failed to add user to db")
	}
	/*
		Lookup the just added user
	*/
	foundUser, err := us.UserDS.Find("email", test.email)
	if err != nil {
		t.Error("failed to lookup/find user in db")
	}

	/*
		check for empty strings and that the retieved user
		has the correct information.
	*/
	if newUser.Name == "" ||
		newUser.Email == "" ||
		newUser.HashedPassword == "" ||
		newUser.TokenUserID == "" {
		t.Error("failed to check for empty input strings")
	}

	if foundUser.Name == "" ||
		foundUser.Email == "" ||
		foundUser.HashedPassword == "" ||
		foundUser.TokenUserID == "" {
		t.Error("empty strings returned from db lookup")
	}

	if newUser.Name != foundUser.Name ||
		newUser.Email != foundUser.Email ||
		newUser.HashedPassword != foundUser.HashedPassword ||
		newUser.TokenUserID != foundUser.TokenUserID {
		t.Error("Error adding user")
	}

	if strings.Contains(newUser.HashedPassword, test.password) {
		t.Error("password is stored in plain text!")
	}

	if len(newUser.TokenUserID) != int(us.Config.TokenIDSize) {
		t.Error("TokenUserID is incorrect size")
	}
}
