package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/markstanden/authentication"
)

type SecretService struct {
	DB DataStore
}

func (ss SecretService) AddSecret(s authentication.Secret) (err error) {
	var id int
	now := time.Now().UTC().Unix()
	sql := "INSERT INTO keys (keyname, keyid, value, created) VALUES ($1, $2, $3, $4) RETURNING id"
	err = ss.DB.QueryRow(sql, s.KeyName, s.KeyID, s.Value, now).Scan(&id)
	if err != nil {
		return err
	}

	// Log addition to database.
	log.Printf("authentication/postgres: id (%d) key (%v) added to db", id, s.KeyName)

	//return the ID of the created user
	return nil
}

func (ss SecretService) GetSecret(name string) func(version string) (secret string) {

	switch name {
	case "SecretKey":
		return func(keyID string) (secret string) {
			var (
				row     *sql.Row
				value   string
				created int64
			)

			row = ss.DB.QueryRow("SELECT value, created FROM keys WHERE keyid = $1", keyID)
			err := row.Scan(&value, &created)
			switch err {
			case sql.ErrNoRows:
				return
			}
			fmt.Println("secretservice/GetSecret KeyID: ", keyID)
			fmt.Println("secretservice/GetSecret Value: ", value)
			return value
		}
	}
	return func(version string) string {
		return ""
	}
}

func (ss SecretService) FullReset() (err error) {
	// If the table already exists, drop it
	_, err = ss.DB.Exec(`DROP TABLE IF EXISTS keys;`)
	if err != nil {
		return fmt.Errorf("authentication/postgres: Failed to drop keys table:\n%v", err)
	}

	// Create the new key table
	_, err = ss.DB.Exec(`CREATE TABLE keys (
    id SERIAL PRIMARY KEY,
	keyname varchar(64) NOT NULL,
    keyid varchar(64) UNIQUE NOT NULL,
    value varchar(255) NOT NULL,
    created integer UNIQUE NOT NULL);`)
	if err != nil {
		return fmt.Errorf("authentication/postgres: Failed to create keys table:\n%v", err)
	}

	log.Println("authentication/postgres: keys table dropped and created ok")
	return nil
}
