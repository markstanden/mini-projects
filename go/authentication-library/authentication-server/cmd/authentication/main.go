package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/markstanden/authentication/cache"
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
	

	// get secrets from the secret store
	keys := []string{"PGPASSWORD"}
	secrets, err := googlecloudsecrets.getSecrets(keys)

	// open a connection to the database
	db := postgres.NewConnection()

	// check the database connection is up and running
	err := db.DB.Ping()
	if err != nil {
		fmt.Println("Connection Failure", err)
	}
	// Close the database when the server ends
	defer db.DB.Close()

	// Create a user cache and shadow the db
	c = cache.NewUserCache(db)

	// Create a handler for our routes, pass in the cache
	http.Handle("/signin", routes.SignIn(c))

	// start the server.
	if err := http.ListenAndServe(":8080", nil); err != nil {
		return errors.New("Failed to Start HTTP server: " + err.Error())
	}
	// return no errors if the app closes normally
	return nil
}
