package googlecloud

import (
	"bytes"
	"context"
	"fmt"
	"io"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

// SecretStore is the base struct of our authentication.SecretStore interface implemetation
// Basically a wrapper for the google cloud API
type SecretStore struct{
	// Wraps google API
	// May add fields to this later
}

// GetSecrets looks up the keys in the provided slice and returns a map[string]string
// key value store.
/* func GetSecrets(keys []string) (secrets authentication.SecretStore, err error) {
	return nil, fmt.Errorf("unable to retrieve secrets: ")
} */

// GetSecret looks up the secret from the store and updates its value in the map
func (secrets SecretStore) GetSecret(project, key, version string) (secret string, err error) {
	
	// request string needs to look like this
	// name := "projects/my-project/secrets/my-secret/versions/5"
	// name := "projects/my-project/secrets/my-secret/versions/latest"
	requestString := fmt.Sprintf("projects/%v/secrets/%v/versions/%v", project, key, version)

	// create a new buffer to receive the secret
	buf := new(bytes.Buffer)
	
	err = accessSecretVersion(buf, requestString)
	if err != nil {
		return
	}
	return buf.String(), nil
}

// accessSecretVersion accesses the payload for the given secret version if one
// exists. The version can be a version number as a string (e.g. "5") or an
// alias (e.g. "latest").
func accessSecretVersion(w io.Writer, name string) error {
        // name := "projects/my-project/secrets/my-secret/versions/5"
        // name := "projects/my-project/secrets/my-secret/versions/latest"
		
        // Create the client.
        ctx := context.Background()
        client, err := secretmanager.NewClient(ctx)
        if err != nil {
                return fmt.Errorf("failed to create secretmanager client: %v", err)
        }

        // Build the request.
        req := &secretmanagerpb.AccessSecretVersionRequest{
                Name: name,
        }

        // Call the API.
        result, err := client.AccessSecretVersion(ctx, req)
        if err != nil {
                return fmt.Errorf("failed to access secret version: %v", err)
        }

        // WARNING: Do not print the secret in a production environment - this snippet
        // is showing how to access the secret material.
        fmt.Fprint(w, string(result.Payload.Data))
        return nil
}