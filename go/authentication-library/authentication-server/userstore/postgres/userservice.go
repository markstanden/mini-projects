package postgres

import (
	"database/sql"
	"errors"
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

/*
	NewUserService returns a new UserService object with the
	supplied datastore.
*/
func NewUserService(db DataStore) (us UserService) {
	return UserService{
		DB: db,
	}
}

/*
	**  Add  **
	adds the user to the Database
*/
func (us UserService) Add(u *authentication.User) (err error) {
	if u.Name == "" || u.Email == "" || u.HashedPassword == "" || u.TokenID == "" {
		return errors.New("missing user data")
	}
	var id int
	query := "INSERT INTO users (name, email, hashedpassword, tokenid) VALUES ($1, $2, $3, $4) RETURNING id"
	err = us.DB.QueryRow(query, u.Name, u.Email, u.HashedPassword, u.TokenID).Scan(&id)
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

/*
	**  Delete  **
	deletes a user from the Database
*/
func (us UserService) Delete(u *authentication.User) (err error) {
	return nil
}

/*
	**  Find  **
	finds the first instance of the key value pair in the database.
	it is intended to search unique keys only.
	valid options for key:
		email 	- The user's entered email address
		tokenid	- The
*/
func (us UserService) Find(key, value string) (u *authentication.User, err error) {
	var row *sql.Row

	switch key {
	case "email":
		row = us.DB.QueryRow("SELECT id, name, email, hashedpassword, tokenid FROM users WHERE email = $1", value)
	case "tokenid":
		row = us.DB.QueryRow("SELECT id, name, email, hashedpassword, tokenid FROM users WHERE tokenid = $1", value)
	default:
		return nil, errors.New("user not found")
	}

	uid := 0
	name := ""
	email := ""
	hashedPassword := ""
	tokenID := ""
	err = row.Scan(&uid, &name, &email, &hashedPassword, &tokenID)

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
			TokenID:        tokenID,
		}, nil
	default:
		log.Println("authentication/sql: user lookup error")
		return nil, err
	}

}

/*
	**  Update  **
	updates a user in the Database
*/
func (us UserService) Update(u *authentication.User) (err error) {
	return nil
}

/*
************** DEVELOPMENT USE ONLY!!! ***************
 */
/*
	**  FullReset  **
	drops and re-Creates the user table
*/
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
    tokenid varchar(160) UNIQUE NOT NULL);`)
	if err != nil {
		return fmt.Errorf("authentication/postgres: Failed to create users table:\n%v", err)
	}

	log.Println("authentication/postgres: users & keys table dropped and created ok")
	return nil
}
