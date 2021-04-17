package main

import (
	"fmt"
	"io"

	"net/http"
	"os"

	"github.com/markstanden/authentication/accesstoken"
	"github.com/markstanden/authentication/datastores/postgres"
	"github.com/markstanden/authentication/datastores/secretstore"
	"github.com/markstanden/authentication/datastores/usercache"
	"github.com/markstanden/authentication/datastores/userstore"
	"github.com/markstanden/authentication/deployment/googlecloud"
	"github.com/markstanden/authentication/routes"
	"github.com/markstanden/authentication/userservice"
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
	pgConfig := postgres.NewConfig().FromEnv()

	/*
		Attempt to retrieve a password from the GCP password store.
		If this fails it returns a empty string
	*/
	if dbPW := gcloud.GetSecret("PGPASSWORD")("latest"); dbPW != "" {
		pgConfig = pgConfig.Password(dbPW)
	}

	/*
		Use the config and connect to the database.
		Panic on failure to connect.
		Close the database when the server ends
	*/
	authdb, err := pgConfig.DBName("authentication").Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to the DataStore...\n %v", err)
	}
	defer authdb.DB.Close()

	// create the userstore instance and connect it to the postgres DB
	userDB := userstore.New(authdb)
	// create the secretstore instance and connect it to our database
	ss := secretstore.New(authdb, 3600)

	// create the userservice
	us := userservice.UserService{

		/* Create a userstore cache and shadow the main userstore */
		UserDS: usercache.New(userDB),

		/* Create a SecretStore to handle our rotating keys */
		SecretDS: ss,

		/* create a token service to create authentication tokens for users */
		AccessTS: &accesstoken.AccessToken{
			Issuer:    "markstanden.dev",
			Audience:  "markstanden.dev",
			MinsValid: 60,
			Secret:    ss,
			StartTime: 1617020114,
		},
	}

	/* Create a handler for our routes, pass in the cache */
	http.Handle("/", routes.Home(us.UserDS))
	http.Handle("/reset-users-table", routes.ResetUsersTable(us.UserDS))
	http.Handle("/reset-keys-table", routes.ResetKeysTable(us.SecretDS))
	http.Handle("/signin", routes.SignIn(us))
	http.Handle("/signup", routes.SignUp(us))

	/* start the server. */
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		return fmt.Errorf("failed to start HTTP server: /n %v", err)
	}
	fmt.Fprintf(stdout, "HTTP Server listening on port :%v", port)

	/* return no errors if the app closes normally */
	return nil
}
