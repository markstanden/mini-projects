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
type DeploymentService struct {
	/*
		Project is the unique project identifier number within Google Cloud Platform.
		The string is used as part of the request string that will be sent to the GCP secret store.
	*/
	ProjectID string
}

func NewSecretHandler() (ds *DeploymentService) {
	return &DeploymentService{
		ProjectID: "145660875199",
	}
}

/*
	GetSecret looks up the secret from the store and updates its value in the map
	project is the project id (ours is a 12 digit int)
	key is the name of the stored data
	version is the fixed version number i.e. "5" or the current value "latest"
*/
func (ds DeploymentService) GetSecret(SecretName string) func(version string) (secret string) {

	// create a new buffer to receive the secret
	buf := new(bytes.Buffer)

	/*
		return a function that can choose the required version
	*/
	return func(version string) (secret string) {

		/*
			a requestString is a path to a secret within the Google Cloud Provider.
			It allows the secret manager to request a particular version of a secret i.e. :
			requestString := "projects/<DeploymentService.Project>/secrets/<SecretName>/versions/5"
			or the latest version of the secret, i.e. :
			requestString := "projects/my-project/secrets/my-secret/versions/latest"
		*/
		requestString := fmt.Sprintf("projects/%v/secrets/%v/versions/%v", ds.ProjectID, SecretName, version)

		if err := accessSecretVersion(buf, requestString); err != nil {
			return ""
		}

		/*
			we need to reset the buffer once done,
			otherwise the next time the function is called the buffer is filled again!
		*/
		defer buf.Reset()

		return buf.String()
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
