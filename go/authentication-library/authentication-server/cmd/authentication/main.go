package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/markstanden/authentication/cache"
	"github.com/markstanden/authentication/deploy/googlecloud"
	"github.com/markstanden/authentication/postgres"
	"github.com/markstanden/authentication/routes"
)

func main() {
	if err := run(os.Args, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(args []string, stdout io.Writer) error {
	
	// Attempt to get port to listen on from ENV variables
	port := os.Getenv("PORT")
	if port == "" {
		// $PORT env variable not set, assume local dev environment
		log.Printf("authentication/main: $PORT Env Variable not set, setting to default")
		port = "8080"
	}

	// create a secret store to pass to the UserStore
	secrets := &googlecloud.SecretStore{}

	// open a connection to the database
	db, err := postgres.NewConnection(secrets)
	if err != nil {
		return fmt.Errorf("error esablishing connection to database: /n %v", err)
	}

	// Close the database when the server ends
	defer db.DB.Close()

	// Create a user cache and shadow the db
	c := cache.NewUserCache(db)

	// Create a handler for our routes, pass in the cache
	http.Handle("/", routes.Home(c))
	http.Handle("/create-users-table", routes.CreateUsersTable(db))
	http.Handle("/signin", routes.SignIn(c))
	http.Handle("/signup", routes.SignUp(c))

	// start the server.
	if err := http.ListenAndServe(":" + port, nil); err != nil {
		return fmt.Errorf("failed to start HTTP server: /n %v", err)
	}
	
	
	/*
	certFile := ""
	keyFile := ""
	if err := http.ListenAndServeTLS(":" + port, certFile, keyFile, nil); err != nil {
		return fmt.Errorf("failed to start HTTPS server: /n %v", err)
	}
	*/
	
	// return no errors if the app closes normally
	return nil
}
