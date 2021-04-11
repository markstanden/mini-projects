package pguserdatastore

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/markstanden/authentication"
	"github.com/markstanden/authentication/datastores/postgres"
)

var (
	ErrEmailAddressAlreadyInUse = errors.New("email address already in use")
	ErrTokenUserIDInUse         = errors.New("TokenUserID already in use")
	ErrInvalidInput             = errors.New("invalid input")
)

// UserService is a struct providing a psql implementation of authentication.UserDataStore
type PGUserDataStore struct {
	DB postgres.DataStore
}

/*
	**  NewUserService  **
	NewUserService returns a new UserService object with the
	supplied datastore.
*/
func NewUserService(db postgres.DataStore) (us PGUserDataStore) {
	return PGUserDataStore{
		DB: db,
	}
}

/*
	**  Add  **
	adds the user to the Database
*/
func (us PGUserDataStore) Add(u *authentication.User) (err error) {
	if u.Name == "" || u.Email == "" || u.HashedPassword == "" || u.TokenUserID == "" {
		return errors.New("missing user data")
	}
	/*
		check the record does not already exist
	*/
	if _, err := us.Find("email", u.Email); err == nil {
		return ErrEmailAddressAlreadyInUse
	}

	/*
		check for duplicated tokenUserID
	*/
	if _, err := us.Find("tokenuserid", u.TokenUserID); err == nil {
		return ErrTokenUserIDInUse
	}

	var id int
	query := "INSERT INTO users (name, email, hashedpassword, tokenuserid, currentrefresh) VALUES ($1, $2, $3, $4, $5) RETURNING uniqueid"
	err = us.DB.QueryRow(query, u.Name, u.Email, u.HashedPassword, u.TokenUserID, u.CurrentRefreshToken).Scan(&id)
	if err != nil {
		return err
	}

	// The current user doesn't have an id set yet, so set it now.
	u.UniqueID = id

	return nil
}

/*
	**  Delete  **
	deletes a user from the Database
*/
func (us PGUserDataStore) Delete(u *authentication.User) (err error) {
	_, err = us.DB.Exec("DELETE FROM users WHERE email = $1", u.Email)
	return err
}

/*
	**  Find  **
	finds the first instance of the key value pair in the database.
	it is intended to search unique keys only.
	valid options for key:
		email 	- The user's entered email address
		tokenuserid	- The
*/
func (us PGUserDataStore) Find(key, value string) (u *authentication.User, err error) {
	var row *sql.Row

	switch key {
	case "email":
		row = us.DB.QueryRow("SELECT uniqueid, name, email, hashedpassword, tokenuserid, currentrefresh FROM users WHERE email = $1", value)
	case "tokenuserid":
		row = us.DB.QueryRow("SELECT uniqueid, name, email, hashedpassword, tokenuserid, currentrefresh FROM users WHERE tokenuserid = $1", value)
	default:
		return nil, errors.New("user not found")
	}

	uniqueID := 0
	name := ""
	email := ""
	hashedPassword := ""
	tokenUserID := ""
	currentRefresh := ""
	err = row.Scan(&uniqueID, &name, &email, &hashedPassword, &tokenUserID, &currentRefresh)

	switch err {
	case sql.ErrNoRows:
		return nil, fmt.Errorf("user not found")
	case nil:
		return &authentication.User{
			UniqueID:            uniqueID,
			Name:                name,
			Email:               email,
			HashedPassword:      hashedPassword,
			TokenUserID:         tokenUserID,
			CurrentRefreshToken: currentRefresh,
		}, nil
	default:
		log.Println("authentication/sql: user lookup error")
		return nil, err
	}

}

/*
	**  Update  **
	updates a user in the Database
	This definitely needs refactoring!
*/
func (us PGUserDataStore) Update(u *authentication.User, updatedFields authentication.User) (err error) {

	fmt.Println(updatedFields)

	if updatedFields.Name == "" || updatedFields.Name == u.Name {
		return ErrInvalidInput
	}
	us.updateName(u.UniqueID, updatedFields.Name)

	if updatedFields.Email == "" || updatedFields.Email == u.Email {
		return ErrInvalidInput
	}
	us.updateEmail(u.UniqueID, updatedFields.Email)

	if updatedFields.HashedPassword == "" || updatedFields.HashedPassword == u.HashedPassword {
		return ErrInvalidInput
	}
	us.updateHashedPW(u.UniqueID, updatedFields.HashedPassword)

	if updatedFields.CurrentRefreshToken == "" || updatedFields.CurrentRefreshToken == u.CurrentRefreshToken {
		return ErrInvalidInput
	}
	us.updateRefresh(u.UniqueID, updatedFields.CurrentRefreshToken)

	return nil
}

/*
	*** updateName ***
	updateName is a private function that updates the users name
	within the datastore.
*/
func (us PGUserDataStore) updateName(uniqueID int, name string) (err error) {
	_, err = us.DB.Exec("UPDATE users SET name = $1 WHERE uniqueid = $2", name, uniqueID)
	return err
}

/*
	*** updateEmail ***
	updateEmail is a private function that updates the users email
	within the datastore.
*/
func (us PGUserDataStore) updateEmail(uniqueID int, email string) (err error) {
	_, err = us.DB.Exec("UPDATE users SET email = $1 WHERE uniqueid = $2", email, uniqueID)
	return err
}

/*
	*** updateHashedPW ***
	updateHashedPW is a private function that updates the users hashed password
	within the datastore.
*/
func (us PGUserDataStore) updateHashedPW(uniqueID int, hashedPW string) (err error) {
	_, err = us.DB.Exec("UPDATE users SET hashedpassword = $1 WHERE uniqueid = $2", hashedPW, uniqueID)
	return err
}

/*
	*** updateRefresh ***
	updateRefresh is a private function that updates the users Refresh token
	within the datastore.
*/
func (us PGUserDataStore) updateRefresh(uniqueID int, refresh string) (err error) {
	_, err = us.DB.Exec("UPDATE users SET currentrefresh = $1 WHERE uniqueid = $2", refresh, uniqueID)
	return err
}

/*
	*** updateAccess ***
	updateAccess is a private function that updates the ID string that identifies the current access token
	within the datastore.
*/
func (us PGUserDataStore) updateAccess(uniqueID int, accessID string) (err error) {
	_, err = us.DB.Exec("UPDATE users SET currentaccess = $1 WHERE uniqueid = $2", accessID, uniqueID)
	return err
}

/*
	**  FullReset  **
	drops and re-Creates the user table
*/
func (us PGUserDataStore) FullReset() (err error) {
	// If the table already exists, drop it
	_, err = us.DB.Exec(`DROP TABLE IF EXISTS users;`)
	if err != nil {
		return fmt.Errorf("authentication/postgres: Failed to drop users table:\n%v", err)
	}

	// Create the new user table
	_, err = us.DB.Exec(`CREATE TABLE users (
    uniqueid SERIAL PRIMARY KEY,
    name varchar(255) NOT NULL,
    email varchar(255) UNIQUE NOT NULL,
    hashedpassword varchar(160) NOT NULL,
    currentrefresh varchar(255) UNIQUE,
    tokenuserid varchar(160) UNIQUE NOT NULL);`)
	if err != nil {
		return fmt.Errorf("authentication/postgres: Failed to create users table:\n%v", err)
	}

	//log.Println("authentication/postgres: users table dropped and created ok")
	return nil
}
