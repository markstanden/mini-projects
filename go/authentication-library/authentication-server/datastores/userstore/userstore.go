/*
	userstore accesses the connected PSQL instance and
	manipulates the users table
*/
package userstore

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/markstanden/authentication"
	"github.com/markstanden/authentication/datastores/postgres"
)

var (
	ErrEmailAddressAlreadyInUse = errors.New("email address already in use")
	ErrTokenUserIDInUse         = errors.New("TokenUserID already in use")
	ErrInvalidInput             = errors.New("invalid input")
)

// Userstore is a struct providing a psql implementation of authentication.UserDataStore
type Userstore struct {
	DB postgres.DataStore
}

/*
	**  New (Userstore)  **
	New returns a new Userstore object with the
	supplied datastore.
*/
func New(db postgres.DataStore) (us Userstore) {
	return Userstore{
		DB: db,
	}
}

/*
	**  Add  **
	adds the user to the Database
*/
func (us Userstore) Add(u *authentication.User) (err error) {
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
func (us Userstore) Delete(u *authentication.User) (err error) {
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
func (us Userstore) Find(key, value string) (u *authentication.User, err error) {
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
*/
func (us Userstore) Update(user *authentication.User, upd authentication.User) (err error) {

	if !valid(upd) {
		return ErrInvalidInput
	}

	_, err = us.DB.Exec("UPDATE users SET (name, email, hashedpassword, tokenuserid) = ($1, $2, $3, $4) WHERE uniqueid = $5", upd.Name, upd.Email, upd.HashedPassword, upd.TokenUserID, upd.UniqueID)
	return err
}

/*
	** UpdateRefresh **
	UpdateRefresh is a private function that updates the users Refresh token
	within the datastore.
*/
func (us Userstore) UpdateRefreshToken(u *authentication.User, refresh string) (err error) {
	if refresh == "" {
		return ErrInvalidInput
	}
	_, err = us.DB.Exec("UPDATE users SET currentrefresh = $1 WHERE uniqueid = $2", refresh, u.UniqueID)
	return err
}

/*
	** valid **
	valid takes a user as an argument and returns true if the fields
	pass the field specific validation checks
*/
func valid(user authentication.User) bool {
	if !validName(user.Name) ||
		!validEmail(user.Email) ||
		!validPW(user.HashedPassword) ||
		!validTokenID(user.TokenUserID) {
		return false
	}
	return true
}

/*
	** validName **
	valid name is a private method that verifies the supplied string as
*/
func validName(name string) bool {
	length := len(name)
	if length < 2 {
		return false
	}
	if length > 100 {
		return false
	}
	return true
}

/*
	*** validEmail ***
	validEmail checks the supplied string is a valid email address
	and returns true if valid.
*/
func validEmail(input string) bool {
	parts := strings.Split(input, "@")
	if len(parts) != 2 {
		return false
	}
	domain := strings.Split(parts[1], ".")
	if len(domain) != 2 {
		return false
	}
	return true
}

/*
	** validPW **
	validPW is a private method that verifies the supplied string as an appropriate password
*/
func validPW(name string) bool {
	length := len(name)
	if length < 2 {
		return false
	}
	if length > 255 {
		return false
	}
	return true
}

func validTokenID(tokenID string) bool {
	length := len(tokenID)
	if length < 32 {
		return false
	}
	if length > 255 {
		return false
	}
	return true
}

/*
	** updateAccess **
	updateAccess is a private function that updates the ID string that identifies the current access token
	within the datastore.
*/
func (us Userstore) updateAccess(uniqueID int, accessID string) (err error) {
	_, err = us.DB.Exec("UPDATE users SET currentaccess = $1 WHERE uniqueid = $2", accessID, uniqueID)
	return err
}

/*
	**  FullReset  **
	drops and re-Creates the user table
*/
func (us Userstore) FullReset() (err error) {
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
