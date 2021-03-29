package postgres

import (
	"database/sql"
	"fmt"
	"log"

	// This is the required postgres driver for the database/sql package
	_ "github.com/lib/pq"
	"github.com/markstanden/authentication"
)

// UserService is a struct providing a psql implementation of authentication.UserService
type UserService struct {
	DB DataStore
}

// Drop drops the user table
func (us UserService) FullReset() (err error) {
	// If the table already exists, drop it
	_, err = us.DB.Exec(`DROP TABLE IF EXISTS users;`)
	if err != nil {
		return fmt.Errorf("authentication/postgres: Failed to drop users table:\n%v", err)
	}

	// Create the new user table
	_, err = us.DB.Exec(`CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name varchar(255) NOT NULL,
    email varchar(255) UNIQUE NOT NULL,
    hashedpassword varchar(160) NOT NULL,
    token varchar(160) UNIQUE NOT NULL);`)
	if err != nil {
		return fmt.Errorf("authentication/postgres: Failed to create users table:\n%v", err)
	}

	log.Println("authentication/postgres: users & keys table dropped and created ok")
	return nil
}

// Find returns the first instance of the key value pair in the database.
// it is intended to search unique keys only (id, email, token)
func (us UserService) Find(key, value string) (u *authentication.User, err error) {
	var row *sql.Row

	switch key {
	case "email":
		row = us.DB.QueryRow("SELECT id, name, email, hashedpassword, token FROM users WHERE email = $1", value)
	case "token":
		row = us.DB.QueryRow("SELECT id, name, email, hashedpassword, token FROM users WHERE token = $1", value)
	}

	uid := 0
	name := ""
	email := ""
	hashedPassword := ""
	token := ""
	err = row.Scan(&uid, &name, &email, &hashedPassword, &token)

	switch err {
	case sql.ErrNoRows:
		log.Println("authentication/sql: user not found")
		return nil, fmt.Errorf("user not found")
	case nil:
		return &authentication.User{
			UniqueID:       uid,
			Name:           name,
			Email:          email,
			HashedPassword: hashedPassword,
			Token:          token,
		}, nil
	default:
		log.Println("authentication/sql: user lookup error")
		return nil, err
	}

}

// Add adds the user to the database
func (us UserService) Add(u *authentication.User) (err error) {
	var id int
	sql := "INSERT INTO users (name, email, hashedpassword, token) VALUES ($1, $2, $3, $4) RETURNING id"
	err = us.DB.QueryRow(sql, u.Name, u.Email, u.HashedPassword, u.Token).Scan(&id)
	if err != nil {
		return err
	}

	// The current user doesn't have an id set yet, so set it now.
	u.UniqueID = id

	// Log addition to database.
	log.Printf("authentication/postgres: user (%d) added to db", id)

	//return the ID of the created user
	return nil
}
