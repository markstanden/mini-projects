package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

type DataStore struct {
	*sql.DB
}

// PGConfig is the Postgres configuration options struct
type pgconfig struct {
	host     string
	port     string
	user     string
	password string
	dbname   string
}

// GetPostgresEnvConfig returns a populated struct with the postgres config options
func getPostgresEnvConfig() (config pgconfig) {

	// Host is the name of host to connect to
	if config.host = os.Getenv("PGHOST"); config.host == "" {
		log.Println("authentication/postgres: PGHOST environment variable not set, using default instead")
		config.host = "localhost"
	}

	// Port is the port number to connect to at the server host, or socket file name extension for Unix-domain connections.
	if config.port = os.Getenv("PGPORT"); config.port == "" {
		log.Println("authentication/postgres: PGPORT environment variable not set, using default instead")
		config.port = "5432"
	}

	// User is the PostgreSQL user name to connect as. Defaults to be the same as the operating system name of the user running the application.
	if config.user = os.Getenv("PGUSER"); config.user == "" {
		log.Println("authentication/postgres: PGUSER environment variable not set, using default instead")
		config.user = "postgres"
	}

	// DBName is the database name. Defaults to be the same as the user name.
	if config.dbname = os.Getenv("PGDATABASE"); config.dbname == "" {
		log.Println("authentication/postgres: PGDATABASE environment variable not set, using default instead")
		//config.dbname = config.user
		config.dbname = "authentication"
	}

	// password is the PostgreSQL user password to connect with. defaults to postgres
	if config.password = os.Getenv("PGPASSWORD"); config.password == "" {
		config.password = "postgres"
	}
	return config
}

// NewConnection returns a new Postgres DB instance
func NewConnection(getPassword func(version string) string) (ds DataStore, err error) {

	// create the config object, taking the non-secret info from the env variables
	config := getPostgresEnvConfig()

	// Password is the password to be used if the server demands password authentication.
	//secretPW, err := us.secrets.GetSecret("145660875199", "PGPASSWORD", "latest")
	secretPW := getPassword("latest")
	if secretPW == "" {
		log.Println("authentication/postgres: PGPASSWORD secret variable not available, using ENV or default instead", err)
	} else {
		// set the config password to the password from the secret store
		log.Println("authentication/postgres: PGPASSWORD secret obtained from the secret store.")
		config.password = secretPW
	}
	// Create a connection string without a password argument
	connectionString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable",
		config.host, config.port, config.user, config.dbname)

	if config.password != "" {
		connectionString = fmt.Sprintf("%s password=%s", connectionString, config.password)
	}

	// Connect to postgres using the connection string
	ds.DB, err = sql.Open("postgres", connectionString)

	if err != nil {
		return ds, err
	}
	return ds, nil
}
