package googlecloud

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

// SecretStore is the base struct of our authentication.SecretStore interface implemetation
// Basically a wrapper for the google cloud API
type DeploymentService struct {
	// Wraps google API
	Project string
}

// GetSecret looks up the secret from the store and updates its value in the map
// project is the project id (ours is a 12 digit int)
// key is the name of the stored data
// version is the fixed version number i.e. "5" or the current value "latest"
func (ds DeploymentService) GetSecret(key string) func(version string) (secret string, err error) {

	// request string needs to look like this
	// name := "projects/my-project/secrets/my-secret/versions/5"
	// name := "projects/my-project/secrets/my-secret/versions/latest"

	// create a new buffer to receive the secret
	buf := new(bytes.Buffer)

	// return the function that can choose the required version
	return func(version string) (secret string, err error) {

		log.Printf("GetSecret ran looking for Version: %v\n", version)
		requestString := fmt.Sprintf("projects/%v/secrets/%v/versions/%v", ds.Project, key, version)

		log.Printf("the request string was: %v\n", requestString)
		err = accessSecretVersion(buf, requestString)
		if err != nil {
			return
		}

		log.Printf("and it returned %v\n", buf.String())
		return buf.String(), nil
	}
}

// accessSecretVersion accesses the payload for the given secret version if one
// exists. The version can be a version number as a string (e.g. "5") or an
// alias (e.g. "latest").
func accessSecretVersion(w io.Writer, name string) error {
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

	fmt.Fprint(w, string(result.Payload.Data))
	return nil
}
