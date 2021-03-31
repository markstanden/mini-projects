package postgres_test

import (
	"testing"

	"github.com/markstanden/authentication/deployment/googlecloud"
	"github.com/markstanden/authentication/userstore/postgres"
)

func GetTestConfig(dbname string) postgres.PGConfig {
	/*
		attempt to to connect to the google secret store (if possible) to retreive secret for production tests
	*/
	pw := googlecloud.NewSecretHandler().GetSecret("PGPASSWORD")("latest")

	return postgres.NewConfig().DBName(dbname).Password(pw)
}

func TestNewConfig(t *testing.T) {
	testCases := []struct {
		desc string
	}{
		{desc: "Test host from env"},
		{desc: "Test port from env"},
		{desc: "Test user from env"},
		{desc: "Test database from env"},
		{desc: "Test password from env"},
		{desc: "Test password from callback"},
		{desc: "Test empty host from env"},
		{desc: "Test empty port from env"},
		{desc: "Test empty user from env"},
		{desc: "Test empty database from env"},
		{desc: "Test empty password from env"},
		{desc: "Test empty password from callback"},
		{
			desc: "",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

		})
	}
}

/*
TODO
Test
	(config PGConfig) Connect() (ds DataStore, err error)
	(config PGConfig) FromEnv() PGConfig
	(config PGConfig) Host(host string) PGConfig
	(config PGConfig) Port(port string) PGConfig
	(config PGConfig) User(user string) PGConfig
	(config PGConfig) DBName(dbname string) PGConfig
	(config PGConfig) Password(password string) PGConfig

*/
