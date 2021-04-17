package cache

import (
	"testing"

	"github.com/markstanden/authentication"
	"github.com/markstanden/authentication/datastores/postgres"
	"github.com/markstanden/authentication/datastores/userdatastores/pguserdatastore"
)

var testuser1 = authentication.User{
	Name:           "Test",
	Email:          "SCcedYbuiRCqN86clZ9Jo0KgxOCoIbfKlI6ySPWiOWvcRJFeqoQAp09nfjIQvm7VpFNr4oZiLlVAFR9P93w+dPzeSyY@test.com",
	HashedPassword: "28ZuYpwprirXKRvnvB3nPLd0NA1VbJfv32yoht9x4pZLYrFOD8W4M15n/TRPLvQ7Ey+MEvPQ8Ln61chXP5IsKggIMf3J73N3VORDLTCKhxmwdxpm_XiEALASgk1LijzClRYpRP6Gzeyq",
	TokenUserID:    "uBIH3ZsbyJ5JUiatbeUViBhqNg0p7uaGuBIH3ZsbyJ5JUiatbeUViBhqNg0p7uaG",
}

/*
	Create a new user cache, with a connection to the test DB
*/
func GetTestUserCache(t *testing.T) *UserServiceCache {
	/*
		Create a connection to the test database
	*/
	testdb, err := postgres.GetTestConfig().Connect()
	if err != nil {
		t.Error("Failed to connect to DataStore\n", err)
	}
	/*
		Wrap the database in a new userservice,
		and wrap the userservice in our cache
	*/
	return NewUserCache(pguserdatastore.PGUserDataStore{DB: testdb})
}

func TestNewUserCache(t *testing.T) {

	usc := GetTestUserCache(t)
	if err := usc.FullReset(); err != nil {
		t.Fatalf("failed to initialise test database")
	}

	/*
		Test that the maps have been created correctly
	*/
	if usc.cache.emailCache == nil || usc.cache.tokenUserIDCache == nil {
		t.Error("cache not initialised correctly")
	}

	if len(usc.cache.emailCache) != 0 || len(usc.cache.tokenUserIDCache) != 0 {
		t.Error("cache contains data!")
	}
}

/* Test the add and Find functions by first adding to the cache/store and then retreiving, first from the store, then the cache */
func TestAddThenFind(t *testing.T) {

	usc := GetTestUserCache(t)
	if err := usc.FullReset(); err != nil {
		t.Fatalf("failed to initialise test database")
	}

	table := []struct {
		desc    string
		user    authentication.User
		isValid bool
	}{
		{
			desc:    "Empty User",
			user:    authentication.User{},
			isValid: false,
		},
		{
			desc: "Full User",
			user: authentication.User{
				Name:           "Full User",
				Email:          "SCcedYbuiRCqN86clZ9Jo0KgxOCoIbfKlI6ySPWiOWvcRJFeqoQAp09nfjIQvm7VpFNr4oZiLlVAFR9P93w+dPzeSyY@test.com",
				HashedPassword: "28ZuYpwprirXKRvnvB3nPLd0NA1VbJfv32yoht9x4pZLYrFOD8W4M15n/TRPLvQ7Ey+MEvPQ8Ln61chXP5IsKggIMf3J73N3VORDLTCKhxmwdxpm_XiEALASgk1LijzClRYpRP6Gzeyq",
				TokenUserID:    "uBIH3ZsbyJ5JUiatbeUViBhqNg0p7uaGuBIH3ZsbyJ5JUiatbeUViBhqNg0p7uaG"},
			isValid: true,
		},
		{
			desc: "Empty Email",
			user: authentication.User{
				Name:           "Empty Email User",
				Email:          "",
				HashedPassword: "28ZuYpwprirXKRvnvB3asdfasdfnPLd0NA1VbJfv32yoht9x4pZLYrFOD8W4M15n/TRPLvQ7Ey+MEvPQ8Ln61chXP5IsKggIMf3J73N3VORDLTCKhxmwdxpm_XiEALASgk1LijzClRYpRP6Gzeyq",
				TokenUserID:    "uBIH3ZsbyJ5JUiasdfasdfatbeUViBhqNg0p7uaGuBIH3ZsbyJ5JUiatbeUViBhqNg0p7uaG"},
			isValid: false,
		},
		{
			desc: "Missing Email",
			user: authentication.User{
				Name:           "Missing Email User",
				HashedPassword: "28ZuYpwprirXKRvnvB3asdfasdfnPLd0NasdfA1VbJfv32yoht9x4pZLYrFOD8W4M15n/TRPLvQ7Ey+MEvPQ8Ln61chXP5IsKggIMf3J73N3VORDLTCasdfasdfKhxmwdxpm_XiEALASgk1LijzClRYpRP6Gzeyq",
				TokenUserID:    "uBIH3ZsbyJ5JUasdfasiasdfasdfatbeUViBhqNg0p7uaGuBIH3ZsbyJ5JUiatbeUViBhqNg0p7uaG"},
			isValid: false,
		},
		{
			desc: "Empty HashedPassword",
			user: authentication.User{
				Name:           "Empty HashedPassword User",
				Email:          "SCcedYbuiRCqN86clZ9Jo0KgxOCo345ggbIbfKlI6ySPWiOWvcRJFeqoQAp09nfjIQvm7VpFNr4oZiLlVAFR9P93w+dPzeSyY@test.com",
				HashedPassword: "",
				TokenUserID:    "uBIH3ZsbyJ5JUiatbe345qfasdfUViBhqNg0p7uaGuBIH3ZsbyJ5JUiatbeUViBhqNg0p7uaG"},
			isValid: false,
		},
		{
			desc: "Missing HashedPassword",
			user: authentication.User{
				Name:        "Missing HashedPassword User",
				Email:       "SCcedYbuiRCqN86clZ9Jo0KgxOCoIbfKnuymylI6ySPWiOWvcRJFeqoQAp09nfjIQvm7VpFNr4oZiLlVAFR9P93w+dPzeSyY@test.com",
				TokenUserID: "uBIH3ZsbyJ5JsdfgUiatbeUViBhqNg0p7uaGuBIH3ZsbyJ5JUiatbeUViBhqNg0p7uaG"},
			isValid: false,
		},
		{
			desc: "Empty TokenUserID",
			user: authentication.User{
				Name:           "Empty TokenUserID User",
				Email:          "SCcedYbuiRCqN86clZ9Jo0KgxOCsdfgsdfgoIbfKlI6ySPWiOWvcRJFeqoQAp09nfjIQvm7VpFNr4oZiLlVAFR9P93w+dPzeSyY@test.com",
				HashedPassword: "28ZuYpwprirXKRvnvB3nPLd0NA1VbJfv32yoht9x4pZLYrFOD8W4M15n/TRPLvQ7Ey+MEvPQ8Ln61chXP5IsKggIMf3asdgadfgJ73N3VORDLTCKhxmwdxpm_XiEALASgk1LijzClRYpRP6Gzeyq",
				TokenUserID:    ""},
			isValid: false,
		},
		{
			desc: "Missing TokenUserID",
			user: authentication.User{
				Name:           "Missing TokenUserID User",
				Email:          "SCcedYbuiRCqN86clZ9Jo0KgxOCoIbfKlI6ySPWiOWvcRJFeqoQAhjlhjklp09nfjIQvm7VpFNr4oZiLlVAFR9P93w+dPzeSyY@test.com",
				HashedPassword: "28ZuYpwprirXKhjklhjklvnvB3nPLd0NA1VbJfv32yoht9x4pZLYrFOD8W4M15n/TRPLvQ7Ey+MEvPQ8Ln61chXP5IsKggIMf3J73N3VORDLTCKhxmwdxpm_XiEALASgk1LijzClRYpRP6Gzeyq"},
			isValid: false,
		},
	}
	for _, test := range table {

		/*
			Add from the test table to the store
		*/
		t.Run("Add User to DB", func(t *testing.T) {
			if err := usc.Add(&test.user); err != nil && test.isValid {
				t.Fatal(test.desc, " - failed to add test user to datastore")
			} else if err == nil && !test.isValid {
				t.Error(test.desc, " - added an unexpected/invalid user from datastore")
			}
		})
	}

	for _, test := range table {
		/*
			Retrieve the entries from the store/cache as appropriate.
		*/
		t.Run("Find User", func(t *testing.T) {
			t.Run("Find User in email cache", func(t *testing.T) {
				/* user should not be in the cache yet so will need to be fetched from the store */
				if _, ok := usc.cache.emailCache[test.user.Email]; ok {
					/* User data is already present in the cache */
					t.Error(test.desc, " - user already present in the cache")
				}

				if _, err := usc.Find("email", test.user.Email); err != nil && test.isValid {
					t.Error(test.desc, " - failed to retrieve test user from datastore")
				} else if err == nil && !test.isValid {
					t.Error(test.desc, " - retrieved an unexpected/invalid user from datastore")
				}

				/* user should now have been fetched from the store and added to the cache */
				if _, ok := usc.cache.emailCache[test.user.Email]; !ok && test.isValid {
					/* User data is not present in the cache */
					t.Error(test.desc, " - user not present in the cache, after being successfully found in the datastore")
				} else if ok && !test.isValid {
					t.Error(test.desc, " - retrieved an unexpected/invalid user from cache")
				}
			})
			t.Run("Find User in token cache", func(t *testing.T) {
				/* user should not be in the cache yet so will need to be fetched from the store */
				if _, ok := usc.cache.tokenUserIDCache[test.user.TokenUserID]; ok {
					/* User data is already present in the cache */
					t.Error(test.desc, " - user already present in the cache")
				}

				if _, err := usc.Find("tokenuserid", test.user.TokenUserID); err != nil && test.isValid {
					t.Error(test.desc, " - failed to retrieve test user from datastore")
				} else if err == nil && !test.isValid {
					t.Error(test.desc, " - retrieved an unexpected/invalid user from datastore")
				}

				/* user should now have been fetched from the store and added to the cache */
				if _, ok := usc.cache.tokenUserIDCache[test.user.TokenUserID]; !ok && test.isValid {
					/* User data is not present in the cache */
					t.Error(test.desc, " - user not present in the cache, after being successfully found in the datastore")
				} else if ok && !test.isValid {
					t.Error(test.desc, " - retrieved an unexpected/invalid user from cache")
				}
			})
		})
	}
}

func TestUpdate(t *testing.T) {
	startuser := authentication.User{
		Name:                "Test",
		Email:               "test@test.com",
		HashedPassword:      "Cm#oG8JTTMbr%CcY!#Ky8yD*KMM!v%LwC^YY889!eaG3s4pzVKT6&dBwrzVK5GdBUm%6i$cL7tUg3M@^3MD$zsPyFhdmojwkkHEc$$7*UZZwLQvVnX%hi327Tcb7AsDo",
		TokenUserID:         "CFJvAi9moQFqteznLkceR5xvnWB7d3bPwPEy3ao6hvhQyYEdN5z8ZREiggESLJbJ",
		CurrentRefreshToken: "hdfkahsdskldjhfalksdhfkaljsdhflshdflahsldfhahdfkajhsdfkhasllksjdhflakshdflahsdlfhasl",
	}
	enduser := authentication.User{
		/*
			The store will sequentially assign numbers to the user,
			and this is the first user, so it's ID will be 1
		*/
		UniqueID:            1,
		Name:                "Testy",
		Email:               "testy@testing.com",
		HashedPassword:      "Cm#oG8JTTMbr%CcY!#Ky8yD*KMM!v%LwC^YY889!eaG3s4pzVKT6&dBwrzVK5GdBUm%6i$cL7tUg3M@^3MD$zsPyFhdmojwkkHEc$$7*UZZwLQvVnX%hi327Tcb7AsDo",
		TokenUserID:         "cP7Pd9RiZyWpuZEweCpDnzSk7zB7aJKj9ZcGgAyJMVBzKMgymh2GVajWxn3hEZ5b",
		CurrentRefreshToken: "qlifakne5mncakvnlaiejflautua3sfdhv,zdgklsjzvckzjxdjchzskuefheaksuhfasjhdfjnzxjdvhdlz",
	}

	var usc *UserServiceCache

	t.Run("Init", func(t *testing.T) {
		usc = GetTestUserCache(t)
		if err := usc.FullReset(); err != nil {
			t.Fatalf("failed to initialise test database")
		}
		/* Create a user in the datastore */
		if err := usc.Add(&startuser); err != nil {
			t.Fatalf("failed to add user to UserStore")
		}
	})

	t.Run("Prep", func(t *testing.T) {
		/* Load the user in the cache */
		userbyemail, err := usc.Find("email", startuser.Email)
		if err != nil {
			t.Fatalf("failed to find user in UserStore by email")
		}
		userbytoken, err := usc.Find("tokenuserid", startuser.TokenUserID)
		if err != nil {
			t.Fatalf("failed to find user in UserStore by tokenUserID")
		}
		if *userbyemail != *userbytoken {
			t.Fatalf("Users are not the same")
		}
	})

	t.Run("Update", func(t *testing.T) {
		/* Update the user.  This should update in the UserStore and delete from the cache */
		if err := usc.Update(&startuser, enduser); err != nil {
			t.Fatalf("Failed to update user")
		}

		if err := usc.UpdateRefreshToken(&startuser, enduser.CurrentRefreshToken); err != nil {
			t.Fatalf("Failed to update user")
		}

		/* Check the user is no longer in the cache */
		//if _, ok := usc.cache.emailCache[startuser.Email]; ok {
		/* user is still in the cache! */
		//	t.Fatalf("Failed to remove user from emailcache on user update")
		//}
		//if _, ok := usc.cache.tokenUserIDCache[startuser.TokenUserID]; ok {
		/* user is still in the cache! */
		//	t.Fatalf("Failed to remove user from tokencache on user update")
		//}
	})
	t.Run("Check Update", func(t *testing.T) {
		/*
			Find the new user in the datastore
			to prove the update has been made,
			and load the user into the cache
		*/
		if u, err := usc.Find("email", enduser.Email); err != nil {
			t.Fatalf("failed to find user in UserStore by email")
		} else if *u != enduser {
			t.Fatalf("returned user is not as expected.\nWanted:\n%v\nGot:\n%v", enduser, *u)
		}
		if u, err := usc.Find("tokenuserid", enduser.TokenUserID); err != nil {
			t.Fatalf("failed to find updated user in UserStore by tokenUserID")
		} else if *u != enduser {
			t.Fatalf("returned user is not updated as expected.\nWanted:\n%v\nGot:\n%v", enduser, *u)
		}
	})
}

func TestDelete(t *testing.T) {
	startuser := authentication.User{
		Name:                "Test",
		Email:               "test@test.com",
		HashedPassword:      "Cm#oG8JTTMbr%CcY!#Ky8yD*KMM!v%LwC^YY889!eaG3s4pzVKT6&dBwrzVK5GdBUm%6i$cL7tUg3M@^3MD$zsPyFhdmojwkkHEc$$7*UZZwLQvVnX%hi327Tcb7AsDo",
		TokenUserID:         "CFJvAi9moQFqteznLkceR5xvnWB7d3bPwPEy3ao6hvhQyYEdN5z8ZREiggESLJbJ",
		CurrentRefreshToken: "jhlasdjfhlasshdflaehqruipowerupoqiweriuqywetuoahsdfasdfasfdjl321asd14asdaf1",
	}

	var usc *UserServiceCache
	var founduser authentication.User

	t.Run("Init", func(t *testing.T) {
		usc = GetTestUserCache(t)
		if err := usc.FullReset(); err != nil {
			t.Fatalf("failed to initialise test database")
		}
		/* Create a user in the datastore */
		if err := usc.Add(&startuser); err != nil {
			t.Fatalf("failed to add user to UserStore")
		}
	})

	t.Run("Prep", func(t *testing.T) {
		/* Load the user in the cache */
		userbyemail, err := usc.Find("email", startuser.Email)
		if err != nil {
			t.Fatalf("failed to find user in UserStore by email")
		}
		userbytoken, err := usc.Find("tokenuserid", startuser.TokenUserID)
		if err != nil {
			t.Fatalf("failed to find user in UserStore by tokenUserID")
		}
		if *userbyemail != *userbytoken {
			t.Fatalf("Users are not the same")
		}
		founduser = *userbyemail
	})

	/*
		Delete the user.
		This should delete in the UserStore
		and delete from the cache
	*/
	t.Run("Delete", func(t *testing.T) {
		if err := usc.Delete(&founduser); err != nil {
			t.Fatalf("Failed to delete user")
		}
	})
	t.Run("Check Deleted from cache", func(t *testing.T) {
		/* Check the user is no longer in the cache */
		if _, ok := usc.cache.emailCache[founduser.Email]; ok {
			/* user is still in the cache! */
			t.Fatalf("Failed to remove user from emailcache on user delete")
		}
		if _, ok := usc.cache.tokenUserIDCache[founduser.TokenUserID]; ok {
			/* user is still in the cache! */
			t.Fatalf("Failed to remove user from tokencache on user delete")
		}
	})

	t.Run("Check main Store and Cache", func(t *testing.T) {
		/*
			Attempt to find the new user in the datastore
			to prove the delete has been made,
			obviously we would expect this to fail
		*/
		if _, err := usc.Find("email", founduser.Email); err == nil {
			t.Fatalf("found user in UserStore by email")
		}
		if _, err := usc.Find("tokenUserID", founduser.TokenUserID); err == nil {
			t.Fatalf("found user in UserStore by tokenUserID")
		}
	})
}
