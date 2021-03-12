package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	// This is the required postgres driver for the database/sql package
	_ "github.com/lib/pq"
	"github.com/markstanden/authentication"
)

// UserService is a struct providing a psql implementation of authentication.UserService
type UserService struct{
  DB *sql.DB
  secret *authentication.SecretStore
}

// PGConfig is the Postgres configuration options struct
type pgconfig struct{
  host string
  port string
  user string
  password string
  dbname string
}

// GetPostgresEnvConfig returns a populated struct with the postgres config options
func getPostgresEnvConfig() (config pgconfig) {
  
  // Host is the name of host to connect to
  if config.host = os.Getenv("PGHOST"); config.host == "" {
    log.Println("PGHOST environment variable not set, using default instead")
    config.host = "localhost"
  } 
  
  // Port is the port number to connect to at the server host, or socket file name extension for Unix-domain connections.
  if config.port = os.Getenv("PGPORT"); config.port == "" {
    log.Println("PGPORT environment variable not set, using default instead")
    config.port = "5432"
  }

  // User is the PostgreSQL user name to connect as. Defaults to be the same as the operating system name of the user running the application.
  if config.user = os.Getenv("PGUSER"); config.user == "" {
    log.Println("PGUSER environment variable not set, using default instead")
    config.user = "postgres"
  }

  // Password is the password to be used if the server demands password authentication.
  if config.password = os.Getenv("PGPASSWORD"); config.password == "" {
    log.Println("PGPASSWORD environment variable not set, using default instead")
    config.password = ""
  }

  // DBName is the database name. Defaults to be the same as the user name.
  if config.dbname = os.Getenv("PGDATABASE"); config.dbname == "" {
    log.Println("PGDATABASE environment variable not set, using default instead")
    config.dbname = config.user
  }
  return config
}

// NewConnection returns a new Postgres DB instance
func NewConnection() (us UserService) {
  
  config := getPostgresEnvConfig()

  // Create a connection string with password argument, incase a password is added at a later date
  //psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+ "password=%s databaseName=%s sslmode=disable", host, port, user, password, databaseName)

  // Create a connection string without a password argument
  connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
    config.host, config.port, config.user, config.password, config.dbname)

  // Connect to postgres using the connection string
  var err error
  us.DB, err = sql.Open("postgres", connectionString)
  if err != nil {
    panic(err)
  }
  return us
}

// Create a new database if required
func (us UserService) Create() {
  us.DB.Exec(`CREATE TABLE IF NOT EXIST (
    id varchar(255) NOT NULL,
    name varchar(255) NOT NULL,
    email varchar(255) NOT NULL,
    hashedpassword varchar(255) NOT NULL,
    token varchar(255) NOT NULL);`)
}

// FindByID returns the first matching user (IDs should be unique) and returns a User object
func (us UserService) FindByID(id string) (u *authentication.User, err error) {
  fmt.Println("Got here", id)
	rows := us.DB.QueryRow("SELECT id, name, email, hashedpassword, token FROM users WHERE id = $1", id)
  fmt.Println(rows)
	
  uid := ""
  name := "" 
  email := ""
  hashedPassword := "";
  token := "";

  err = rows.Scan(&uid, &name, &email, &hashedPassword, &token)
  if err != nil {
    log.Println(err)
		return nil, err
	}

  return &authentication.User{
    UniqueID: uid,
    Name: name,
    Email: email,
    HashedPassword: hashedPassword,
    Token: token,
  }, nil
}

// FindByEmail returns the first matching user (Emails should be unique) and returns a User object
func (us UserService) FindByEmail(em string) (u *authentication.User, err error) {
  fmt.Println("Got here", em)
	rows := us.DB.QueryRow("SELECT id, name, email, hashedpassword, token FROM users WHERE email = $1", em)
  fmt.Println(rows)
	
  uid := ""
  name := "" 
  email := ""
  hashedPassword := "";
  token := "";

  err = rows.Scan(&uid, &name, &email, &hashedPassword, &token)
  if err != nil {
    log.Println(err)
		return nil, err
	}

  return &authentication.User{
    UniqueID: uid,
    Name: name,
    Email: email,
    HashedPassword: hashedPassword,
    Token: token,
  }, nil
}

// FindByToken returns the first matching user (Tokens should be unique) and returns a User object
func (us UserService) FindByToken(t string) (u *authentication.User, err error) {
  	rows := us.DB.QueryRow("SELECT id, name, email, hashedpassword, token FROM users WHERE token = $1", t)
  fmt.Println(rows)
	
  uid := ""
  name := "" 
  email := ""
  hashedPassword := "";
  token := "";

  err = rows.Scan(&uid, &name, &email, &hashedPassword, &token)
  if err != nil {
    log.Println(err)
		return nil, err
	}

  return &authentication.User{
    UniqueID: uid,
    Name: name,
    Email: email,
    HashedPassword: hashedPassword,
    Token: token,
  }, nil
}




