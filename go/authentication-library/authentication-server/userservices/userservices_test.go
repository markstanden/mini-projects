package userservice_test

import (
	"strings"
	"testing"

	"github.com/markstanden/authentication"
	"github.com/markstanden/authentication/datastores/postgres"
	"github.com/markstanden/authentication/datastores/userdatastores/pguserdatastore"
	userservice "github.com/markstanden/authentication/userservices"
)

/*
	contains the parameters for the NewUser function
*/
type newUser struct {
	name     string
	email    string
	password string
}
type test struct {
	/*
		A test description to aid debuggin if a test should fail.
	*/
	desc string

	/*
		the user to be tested
	*/
	user newUser

	/*
		flag identifying whether the user under test is valid and should pass,
		or is invalid and should create an error
	*/
	isValid bool
}

/*
	tests is a slice of test data to be run against the userservice methods
*/
var tests = []test{

	{
		desc: "valid user",
		user: newUser{
			name:     "Testy 'first' McTestface",
			email:    "userservice@mctestface.com",
			password: "livetotest",
		},
		isValid: true,
	},
	{
		desc: "empty name",
		user: newUser{
			name:     "",
			email:    "emptyname@mctestface.com",
			password: "livetotest",
		},
		isValid: false,
	},

	{
		desc: "empty email",
		user: newUser{
			name:     "Testy McTestface",
			email:    "",
			password: "livetotest",
		},
		isValid: false,
	},

	{
		desc: "empty password",
		user: newUser{
			name:     "Testy McTestface",
			email:    "emptypassword@mctestface.com",
			password: "",
		},
		isValid: false,
	},

	{
		desc: "duplicate user",
		user: newUser{
			name:     "Testy 'second' McTestface",
			email:    "userservice@mctestface.com",
			password: "livetotesttwice",
		},
		isValid: false,
	},

	{
		desc: "invalid email - no @",
		user: newUser{
			name:     "Testy McTestface",
			email:    "testymctestface.com",
			password: "livetotest",
		},
		isValid: false,
	},

	{
		desc: "invalid email - no top level domain",
		user: newUser{
			name:     "Testy McTestface",
			email:    "testy@mctestfacecom",
			password: "livetotest",
		},
		isValid: false,
	},
}

/*
	*** TestNewUser ***
	TestNewUser creates a new user from the test slice,
	and checks that the user has been created and added
	to the test datastore.
*/
func TestNewUserAndLogin(t *testing.T) {

	/*
		connect to the test database
	*/
	ds, err := postgres.GetTestConfig().FromEnv().Connect()
	if err != nil {
		t.Fatal("failed to connect to the temp database")
	}

	userStore := pguserdatastore.NewUserService(ds)
	//	userStore.FullReset()
	/*
		create the UserService using the test DB
	*/
	us := userservice.UserService{
		UserDS: userStore,
		Config: userservice.USConfig{TokenIDSize: 32},
	}

	var usersCreated []authentication.User

	for _, test := range tests {

		var newUser, foundUser *authentication.User

		t.Run(test.desc, func(t *testing.T) {

			t.Run("NewUser", func(t *testing.T) {

				/*
					Add the user to the database
				*/
				newUser, err = us.NewUser(test.user.name, test.user.email, test.user.password)
				if test.isValid && err != nil {
					t.Fatalf("failed to add user (%v) to db", test.user.email)
				} else if !test.isValid && err == nil {
					t.Fatalf("failed to return an error for an invalid user (%v)", test.user.email)
				}

				if err == nil {
					usersCreated = append(usersCreated, *newUser)
				}
			})

			t.Run("Login", func(t *testing.T) {
				/*
					Lookup the added users
				*/

				foundUser, err = us.Login(test.user.email, test.user.password)
				if test.isValid && err != nil {
					t.Fatalf("failed to find user (%v) in db", test.user.email)
				} else if !test.isValid && err == nil {
					t.Fatalf("failed to return an error for an invalid user (%v) lookup", test.user.email)
				}

			})

			t.Run("Validation Tests", func(t *testing.T) {
				/*
					check for empty strings and that the retieved user
					has the correct information.
				*/
				if !test.isValid && foundUser != nil {
					if err == nil {
						foundUser = newUser
					}
					t.Fatalf("invalid test returned a found user")
				}

				if test.isValid {

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

					if strings.Contains(newUser.HashedPassword, test.user.password) {
						t.Error("password is stored in plain text!")
					}

					if len(newUser.TokenUserID) != int(us.Config.TokenIDSize) {
						t.Error("TokenUserID is incorrect size")
					}
				}
			})

		})
	}
	/*
		Delete the created user account
	*/
	for _, u := range usersCreated {
		t.Run("Delete User", func(t *testing.T) {
			if u.Email != "" {
				us.UserDS.Delete(&u)
			}
		})

	}
}

/***
*
*	Benchmarks
*
***/

/*
	required to keep the bechmark tests from being collected
*/
var (
	benchResult *authentication.User
)

func BenchmarkNewUser(b *testing.B) {

	ds, err := postgres.GetTestConfig().FromEnv().Connect()
	if err != nil {
		b.Fatal("failed to connect to the temp database")
	}

	userStore := pguserdatastore.NewUserService(ds)
	us := userservice.UserService{
		UserDS: userStore,
		Config: userservice.USConfig{TokenIDSize: 32},
	}

	var user *authentication.User
	for n := 0; n < b.N; n++ {
		user, _ = us.NewUser("name", "email@address.com", "password")
	}

	benchResult = user
}
