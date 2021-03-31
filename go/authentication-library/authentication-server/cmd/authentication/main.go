package main

import (
	"fmt"
	"io"

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

	/*
		Set default port for server to listen on
	*/
	port := "8080"

	/*
		Attempt to get port to listen on from ENV variables
	*/
	if portENV, ok := os.LookupEnv("PORT"); ok {
		/* $PORT env variable set */
		fmt.Fprintf(stdout, "authentication/main: $PORT Env Variable set, setting to %v", portENV)
		port = portENV
	}

	/*
		create a secret store to pass to the UserStore
	*/
	gcloud := &googlecloud.DeploymentService{
		ProjectID: "145660875199",
	}

	/*
		initialise and build the PGConfig configuration struct.
		set the default DB name to authentication, and then add the FromEnv
		to override the defaults if required
	*/
	pgConfig := postgres.NewConfig().DBName("authentication").FromEnv()
	/*
		Attempt to retrieve a password from the GCP password store.
		If this fails it returns a empty string
	*/
	if dbPW := gcloud.GetSecret("PGPASSWORD")("latest"); dbPW != "" {
		pgConfig = pgConfig.Password(dbPW)
	}
	/*
		Use the config and connect to the database
	*/
	db, err := pgConfig.Connect()
	if err != nil {
		return fmt.Errorf("error establishing connection to database: /n %v", err)
	}

	/*
		Close the database when the server ends
	*/
	defer db.DB.Close()

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
	fmt.Fprintf(stdout, "HTTP Server listening on port :%v", port)

	// return no errors if the app closes normally
	return nil
}
