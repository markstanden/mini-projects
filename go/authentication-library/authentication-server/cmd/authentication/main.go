package main

import (
	"fmt"
	"io"
	"log"

	"net/http"
	"os"

	"github.com/markstanden/authentication/deployment/googlecloud"
	"github.com/markstanden/authentication/routes"
	"github.com/markstanden/authentication/tokenservice"
	"github.com/markstanden/authentication/userstore/cache"
	"github.com/markstanden/authentication/userstore/postgres"
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
	gcloud := &googlecloud.DeploymentService{
		Project: "145660875199",
	}

	// open a connection to the database
	db, err := postgres.NewConnection(gcloud.GetSecret("PGPASSWORD"))
	if err != nil {
		return fmt.Errorf("error establishing connection to database: /n %v", err)
	}

	us := postgres.UserService{DB: db}
	ss := postgres.SecretService{DB: db, Lifespan: 3600}

	// create a token service to create authentication tokens for users
	userTokens := &tokenservice.TokenService{
		Issuer:     "markstanden.dev",
		Audience:   "markstanden.dev",
		HoursValid: 24,
		Secret:     ss,
		StartTime:  1617020114,
	}

	// Close the database when the server ends
	defer db.DB.Close()

	// Create a user cache and shadow the db
	cache := cache.NewUserCache(us)

	// Create a handler for our routes, pass in the cache
	http.Handle("/", routes.Home(cache))
	http.Handle("/reset-users-table", routes.ResetUsersTable(cache))
	http.Handle("/reset-keys-table", routes.ResetKeysTable(ss))
	http.Handle("/signin", routes.SignIn(cache))
	http.Handle("/signup", routes.SignUp(cache, userTokens))

	// start the server.
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		return fmt.Errorf("failed to start HTTP server: /n %v", err)
	}

	// return no errors if the app closes normally
	return nil
}
