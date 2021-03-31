package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/markstanden/authentication"
	"github.com/markstanden/securerandom"
)

type SecretService struct {
	DB DataStore

	Lifespan int64
}

func NewSecretService(db DataStore, lifespan int64) (ss SecretService) {
	return SecretService{
		DB:       db,
		Lifespan: lifespan,
	}
}

/*
	GetKeyID returns the latest KeyID provided there is one within the validity window.
	Valid Key present:
		- Returns the KeyID with the largest creation date within the validy window
	Valid Key not found:
		- Triggers the creation of a new key, and returns the created Key ID
*/
func (ss SecretService) GetKeyID(keyName string) (keyID string) {
	now := time.Now().UTC().Unix()
	earliestValid := now - ss.Lifespan
	fmt.Println("SecretService/GetKeyID: earliestValid", earliestValid)
	query := `SELECT keyid FROM keys WHERE keyname = $1 AND created <= $2 GROUP BY keyid HAVING MAX(created) > $3`
	row := ss.DB.QueryRow(query, keyName, now, earliestValid)
	err := row.Scan(&keyID)
	log.Println("SecretService/GetKeyID:\n\tkeyID:\n\t", keyID)
	switch err {
	case sql.ErrNoRows:
		log.Println("SecretService/GetKeyID:\n\tErrNoRows Reached")
		s := authentication.Secret{
			KeyName: keyName,
			KeyID:   securerandom.String(16),
			Value:   securerandom.String(128),
			Created: now}
		err := ss.AddSecret(s)
		if err != nil {
			return ""
		}
		return s.KeyID
	default:
		return keyID
	}
}

func (ss SecretService) AddSecret(s authentication.Secret) (err error) {
	fmt.Println("SecretService/AddSecret:\n\tRequest to add secret made.\n\tSecret:\n\t", s)
	query := "INSERT INTO keys (keyname, keyid, value, created) VALUES ($1, $2, $3, $4)"
	_, err = ss.DB.Exec(query, s.KeyName, s.KeyID, s.Value, s.Created)
	if err != nil {
		log.Println("SecretService/AddSecret:\n\terr:\n\t", err)
		return err
	}

	// Log addition to database.
	log.Printf("SecretService/AddSecret:\n\tKeyName %v\n\tKeyID %v\nSuccessfully added to db", s.KeyName, s.KeyID)

	return nil
}

func (ss SecretService) GetSecret(keyName string) func(keyID string) (value string) {

	switch keyName {
	case "JWT":
		return func(keyID string) (value string) {
			row := ss.DB.QueryRow("SELECT value FROM keys WHERE keyid = $1", keyID)
			err := row.Scan(&value)
			switch err {
			case sql.ErrNoRows:
				return
			}
			log.Println("secretservice/GetSecret Secret Request Made:\nKeyID: ", keyID)
			log.Println("secretservice/GetSecret Value: ", value)
			return value
		}
	}
	return func(keyID string) string {
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
