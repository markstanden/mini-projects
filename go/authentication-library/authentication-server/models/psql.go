package models

import (
	"database/sql"
	"fmt"
	"log"

	// This is the required postgres driver for the database/sql package
	_ "github.com/lib/pq"
)

// Psql is a more generic label for the store connection
type Psql struct{
  DB *sql.DB
}

// NewConnection returns a new Postgres DB instance
func NewConnection(host, username, password, databaseName string, port int) (p Psql) {
  // Create a connection string with password argument, incase a password is added at a later date
  //psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+ "password=%s databaseName=%s sslmode=disable", host, port, user, password, databaseName)

  // Create a connection string without a password argument
  connectionString := fmt.Sprintf("host=%s port=%d user=%s "+
    "dbname=%s sslmode=disable",
    host, port, username, databaseName)

  // Connect to postgres using the connection string
  var err error
  p.DB, err = sql.Open("postgres", connectionString)
  if err != nil {
    panic(err)
  }
  return p
}

// Close closes the connection to the database
func (p Psql) Close() {
	p.DB.Close()
}

// Create a new database if required
func (p Psql) Create() {
  p.DB.Exec(`CREATE TABLE IF NOT EXIST (
    id varchar(255) NOT NULL,
    name varchar(255) NOT NULL,
    email varchar(255) NOT NULL,
    hashedpassword varchar(255) NOT NULL,
    token varchar(255) NOT NULL);`)
}

// FindByID returns the first matching user (IDs should be unique) and returns a User object
func (p Psql) FindByID(id string) (u *User, err error) {
  fmt.Println("Got here", id)
	rows := p.DB.QueryRow("SELECT id, name, email, hashedpassword, token FROM users WHERE id = $1", id)
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

  return &User{
    UniqueID: uid,
    Name: name,
    Email: email,
    HashedPassword: hashedPassword,
    Token: token,
  }, nil
}

// FindByEmail returns the first matching user (Emails should be unique) and returns a User object
func (p Psql) FindByEmail(em string) (u *User, err error) {
  fmt.Println("Got here", em)
	rows := p.DB.QueryRow("SELECT id, name, email, hashedpassword, token FROM users WHERE email = $1", em)
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

  return &User{
    UniqueID: uid,
    Name: name,
    Email: email,
    HashedPassword: hashedPassword,
    Token: token,
  }, nil
}

// FindByToken returns the first matching user (Tokens should be unique) and returns a User object
func (p Psql) FindByToken(t string) (u *User, err error) {
  	rows := p.DB.QueryRow("SELECT id, name, email, hashedpassword, token FROM users WHERE token = $1", t)
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

  return &User{
    UniqueID: uid,
    Name: name,
    Email: email,
    HashedPassword: hashedPassword,
    Token: token,
  }, nil
}




