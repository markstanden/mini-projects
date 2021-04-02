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
	query := "INSERT INTO users (name, email, hashedpassword, tokenid) VALUES ($1, $2, $3, $4) RETURNING uniqueid"
	err = us.DB.QueryRow(query, u.Name, u.Email, u.HashedPassword, u.TokenID).Scan(&id)
	if err != nil {
		return err
	}

	// The current user doesn't have an id set yet, so set it now.
	u.UniqueID = id

	// Log addition to database.
	//log.Printf("authentication/postgres: user (%d) added to db", id)

	//return the ID of the created user
	return nil
}

/*
	**  Delete  **
	deletes a user from the Database
*/
func (us UserService) Delete(u *authentication.User) (err error) {
	_, err = us.DB.Exec("DELETE FROM users WHERE email = $1", u.Email)
	return err
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
		row = us.DB.QueryRow("SELECT uniqueid, name, email, hashedpassword, tokenid FROM users WHERE email = $1", value)
	case "tokenid":
		row = us.DB.QueryRow("SELECT uniqueid, name, email, hashedpassword, tokenid FROM users WHERE tokenid = $1", value)
	default:
		return nil, errors.New("user not found")
	}

	uniqueID := 0
	name := ""
	email := ""
	hashedPassword := ""
	tokenID := ""
	err = row.Scan(&uniqueID, &name, &email, &hashedPassword, &tokenID)

	switch err {
	case sql.ErrNoRows:
		return nil, fmt.Errorf("user not found")
	case nil:
		return &authentication.User{
			UniqueID:       uniqueID,
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
	This definitely needs refactoring!
*/
func (us UserService) Update(u *authentication.User, updatedFields authentication.User) (err error) {

	if updatedFields.Name != "" && updatedFields.Name != u.Name {
		us.updateName(u.UniqueID, updatedFields.Name)
	}
	if updatedFields.Email != "" && updatedFields.Email != u.Email {
		us.updateEmail(u.UniqueID, updatedFields.Email)
	}
	if updatedFields.HashedPassword != "" && updatedFields.HashedPassword != u.HashedPassword {
		us.updateHashedPW(u.UniqueID, updatedFields.HashedPassword)
	}
	if updatedFields.TokenID != "" && updatedFields.TokenID != u.TokenID {
		us.updateTokenID(u.UniqueID, updatedFields.TokenID)
	}

	return err
}

func (us UserService) updateName(uniqueID int, name string) (err error) {
	_, err = us.DB.Exec("UPDATE users SET name = $1 WHERE uniqueid = $2", name, uniqueID)
	return err
}
func (us UserService) updateEmail(uniqueID int, email string) (err error) {
	_, err = us.DB.Exec("UPDATE users SET email = $1 WHERE uniqueid = $2", email, uniqueID)
	return err
}
func (us UserService) updateHashedPW(uniqueID int, hashedPW string) (err error) {
	_, err = us.DB.Exec("UPDATE users SET hashedpassword = $1 WHERE uniqueid = $2", hashedPW, uniqueID)
	return err
}
func (us UserService) updateTokenID(uniqueID int, tokenID string) (err error) {
	_, err = us.DB.Exec("UPDATE users SET tokenid = $1 WHERE uniqueid = $2", tokenID, uniqueID)
	return err
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
    uniqueid SERIAL PRIMARY KEY,
    name varchar(255) NOT NULL,
    email varchar(255) UNIQUE NOT NULL,
    hashedpassword varchar(160) NOT NULL,
    tokenid varchar(160) UNIQUE NOT NULL);`)
	if err != nil {
		return fmt.Errorf("authentication/postgres: Failed to create users table:\n%v", err)
	}

	//log.Println("authentication/postgres: users table dropped and created ok")
	return nil
}
