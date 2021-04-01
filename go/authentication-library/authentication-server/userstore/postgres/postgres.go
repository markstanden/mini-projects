package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
)

type DataStore struct {
	*sql.DB
}

/*
	pgConfig is the Postgres configuration options struct
*/
type PGConfig struct {
	/*
		host - The host url for the postgres server
		defaults to "localhost"
	*/
	host string

	/*
		port - the port that the postgres server is listening on
		default is port "5432"
	*/
	port string

	/*
		user - The username to connect as
		defaults to "postgres"
	*/
	user string

	/*
		password - The password to authenticate the user with the postgres server
		defaults to "postgres" but can also be unset ""
	*/
	password string

	/*
		the name of the database to connect to.
		defaults to "postgres"
	*/
	dbname string
}

/*
	Connect returns a DataStore object containing a connection to a postgres database taking connection options from the configuration object.
*/
func (config PGConfig) Connect() (ds DataStore, err error) {

	/*
		Create a connection string without a password argument
	*/
	connectionString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable",
		config.host, config.port, config.user, config.dbname)

	if config.password != "" {
		connectionString = fmt.Sprintf("%s password=%s", connectionString, config.password)
	}

	/*
		Connect to postgres using the connection string
	*/
	ds.DB, err = sql.Open("postgres", connectionString)
	return ds, err
}

/*
	NewConfig returns a populated struct with the default postgres config options.
	The defaults can be overridden using the builder methods
*/
func NewConfig() PGConfig {
	return PGConfig{
		host:     "localhost",
		port:     "5432",
		user:     "postgres",
		password: "postgres",
		dbname:   "postgres",
	}
}

/*
	FromEnv attempts to take the connnection options from the ENV variables.  if not set, the values remain as previously set.  This enables the function to be placed wherever most convenient in the builder chain.
*/
func (config PGConfig) FromEnv() PGConfig {

	/*
		PGHOST is the name of host to connect to
	*/
	if host, ok := os.LookupEnv("PGHOST"); ok {
		config.host = host
	} else {
		log.Println("authentication/postgres: PGHOST environment variable not set, host unchanged")
	}

	/*
		PGPORT is the port number to connect to at the server host, or socket file name extension for Unix-domain connections.
	*/
	if port, ok := os.LookupEnv("PGPORT"); ok {
		if _, err := strconv.ParseInt(port, 10, 32); err == nil {
			config.port = port
		} else {
			log.Println("authentication/postgres: invalid PGPORT environment variable, port unchanged")
		}
	} else {
		log.Println("authentication/postgres: PGPORT environment variable not set, port unchanged")
	}

	/*
		PGUSER is the PostgreSQL user name to connect as. Defaults to be the same as the operating system name of the user running the application.
	*/
	if user, ok := os.LookupEnv("PGUSER"); ok {
		config.user = user
	} else {
		log.Println("authentication/postgres: PGUSER environment variable not set, username unchanged")
	}

	/*
		PGDATABASE is the database name to connect to
	*/
	if db, ok := os.LookupEnv("PGDATABASE"); ok {
		config.dbname = db
	} else {
		log.Println("authentication/postgres: PGDATABASE environment variable not set, database name unchanged")
	}

	/*
		PGPASSWORD is the PostgreSQL user password to connect with.
	*/
	if pw, ok := os.LookupEnv("PGPASSWORD"); ok {
		config.password = pw
	} else {
		log.Println("authentication/postgres: PGPASSWORD environment variable not set, password unchanged")
	}

	return config
}

/*
	Host changes the config hostname to the supplied string.
*/
func (config PGConfig) Host(host string) PGConfig {
	if host != "" {
		config.host = host
	}
	return config
}

/*
	Host changes the config port to the supplied string.
*/
func (config PGConfig) Port(port string) PGConfig {
	if num, err := strconv.ParseInt(port, 10, 32); err == nil {
		if num < 1 || num > 65535 {
			log.Println("authentication/postgres: port override value is outside valid range, port number unchanged")
			return config
		}

		config.port = port
	} else {
		log.Println("authentication/postgres: invalid port override, port unchanged")
	}
	return config
}

/*
	Host changes the config user to the supplied string.
*/
func (config PGConfig) User(user string) PGConfig {
	if user != "" {
		config.user = user
	}
	return config
}

/*
	Host changes the config database name to the supplied string.
*/
func (config PGConfig) DBName(dbname string) PGConfig {
	if dbname != "" {
		config.dbname = dbname
	}
	return config
}

/*
	Host changes the config password to the supplied string.
*/
func (config PGConfig) Password(password string) PGConfig {
	config.password = password
	return config
}
