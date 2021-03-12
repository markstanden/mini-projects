package googlecloud

import (
	"fmt"
)

// SecretStore is the base struct of our authentication.SecretStore interface implemetation
// Basically a wrapper for the google cloud API
type SecretStore struct{
	// Wraps google API
}

// GetSecrets looks up the keys in the provided slice and returns a map[string]string
// key value store.
/* func GetSecrets(keys []string) (secrets authentication.SecretStore, err error) {
	return nil, fmt.Errorf("unable to retrieve secrets: ")
} */

// UpdatedSecret looks up the secret from the store and updates its value in the map
func (secrets SecretStore) GetSecret(key string) (secret string, err error) {
	return "", fmt.Errorf("unable to retrieve secret: ")
}