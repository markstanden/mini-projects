package googlecloudsecrets

import (
	"fmt"

	"github.com/markstanden/authentication"
)

// SecretStore is the base struct of our authentication.SecretStore interface implemetation
type SecretStore struct{
	Secrets map[string] string
}

// GetSecrets looks up the keys in the provided slice and returns a map[string]string
// key value store.
func GetSecrets(keys []string) (secrets authentication.SecretStore, err error) {
	return nil, fmt.Errorf("Unable to retrieve secrets: ")
}

func (secrets SecretStore) updatedSecret(key string) (secret string, err error) {
	return "", fmt.Errorf("Unable to retrieve secret: ")
}