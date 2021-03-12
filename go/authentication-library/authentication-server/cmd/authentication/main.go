package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/markstanden/authentication/cache"
	"github.com/markstanden/authentication/deploy/googlecloud"
	"github.com/markstanden/authentication/postgres"
	"github.com/markstanden/authentication/routes"
)

var c *cache.UserCache

func main() {
	if err := run(os.Args, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(args []string, stdout io.Writer) error {

	// create a secret store to pass to the UserStore
	secrets := &googlecloud.SecretStore{}

	// open a connection to the database
	db, err := postgres.NewConnection(secrets)
	if err != nil {
		return fmt.Errorf("error esablishing connection to database: /n %v", err)
	}
	
	db.Create()
	
	// check the database connection is up and running
	err = db.DB.Ping()
	if err != nil {
		return fmt.Errorf("error checking connection to database: /n %v", err)
	}
	// Close the database when the server ends
	defer db.DB.Close()

	// Create a user cache and shadow the db
	c = cache.NewUserCache(db)

	// Create a handler for our routes, pass in the cache
	http.Handle("/signin", routes.SignIn(c))
	http.Handle("/signup", routes.SignUp(c))

	// start the server.
	if err := http.ListenAndServe(":8080", nil); err != nil {
		return fmt.Errorf("failed to start HTTP server: /n %v", err)
	}
	// return no errors if the app closes normally
	return nil
}
