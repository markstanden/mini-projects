package models

import (
	"database/sql"
	"fmt"

	// This is the required postgres driver for the database/sql package
	_ "github.com/lib/pq"
)

// Store is a more generic label for the store connection
type Store struct{
  db *sql.DB
}

// New returns a new Postgres DB instance
func New(host, user, password, databaseName string, port int) (psqlDB Store){
  // Create a connection string with password argument, incase a password is added at a later date
  //psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+ "password=%s databaseName=%s sslmode=disable", host, port, user, password, databaseName)

  // Create a connection string without a password argument
  connectionString := fmt.Sprintf("host=%s port=%d user=%s "+
    "dbname=%s sslmode=disable",
    host, port, user, databaseName)

  // Connect to postgres using the connection string
  var err error
  psqlDB.db, err = sql.Open("postgres", connectionString)
  if err != nil {
    panic(err)
  }
  
  return psqlDB
}

// Close closes the connection to the database
func (s Store) Close() {
	s.db.Close()
}





