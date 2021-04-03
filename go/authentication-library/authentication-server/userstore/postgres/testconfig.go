package postgres

import (
	"os"

	"github.com/markstanden/authentication/deployment/googlecloud"
)

func GetTestConfig() PGConfig {
	/*
		create a new default config but with the dbname "test"
	*/
	config := NewConfig().DBName(os.Getenv("PGTESTDB")).Host(os.Getenv("PGTESTHOST")).Port(os.Getenv("PGTESTPORT"))

	/*
		attempt to to connect to the google secret store (if possible) to retreive secret for production tests
	*/
	pw := googlecloud.NewSecretHandler().GetSecret("PGTESTPASSWORD")("latest")
	/*
		GetSecret returns an empty string on failure, and will fail if in development env
	*/
	if pw == "" {
		/* if the GCP password fails, or is empty just use defaults */
		return config
	}
	/* override default password with gcp password */
	return config.Password(pw)
}
