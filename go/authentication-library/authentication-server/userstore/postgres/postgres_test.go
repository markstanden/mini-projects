package postgres

import (
	"testing"

	"github.com/markstanden/authentication/deployment/googlecloud"
)

func GetTestConfig() PGConfig {
	/*
		create a new default config but with the dbname "test"
	*/
	config := NewConfig().DBName("test")

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

func TestHost(t *testing.T) {
	testCases := []struct {
		desc string
		host string
		valid bool
	}{
		{desc: "Test host from string", host: "testhost", valid: true},
		{desc: "Test empty host from string", host: "", valid: false},
		{desc: "Test ip from string", host: "127.0.0.1", valid: true},
	}
	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			config := GetTestConfig().Host(test.host)
			if test.valid && config.host != test.host {
				t.Errorf("Failed to set host using builder function")
			}
			if !test.valid && config.host == test.host {
				t.Errorf("Override of known invalid host made")
			}
		})
	}
}

func TestPort(t *testing.T) {
	testCases := []struct {
		desc string
		port string
		valid bool
	}{
		{desc: "Test port 0", port: "0", valid: false},
		{desc: "Test port 1 (first possible)", port: "1", valid: true},
		{desc: "Test port 65535 (max)", port: "65535", valid: true},
		{desc: "Test port 65536 (1 too high)", port: "65536", valid: false},
		{desc: "Test empty port", port: "", valid: false},
		{desc: "Test non numeric port", port: "127.0.0.1", valid: false},
	}
	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			config := GetTestConfig().Port(test.port)
			if test.valid && config.port != test.port {
				t.Errorf("Failed to set port using builder function")
			}
			if !test.valid && config.port == test.port {
				t.Errorf("Override of known invalid port made")
			}
		})
	}
}

func TestUser(t *testing.T) {
	testCases := []struct {
		desc string
		user string
		valid bool
	}{
		{desc: "Test user from string", user: "testuser", valid: true},
		{desc: "Test empty user from string", user: "", valid: false},
	}
	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			config := GetTestConfig().User(test.user)
			if test.valid && config.user != test.user {
				t.Errorf("Failed to set user using builder function")
			}
			if !test.valid && config.user == test.user {
				t.Errorf("Override of known invalid user made")
			}
		})
	}
}

func TestDBName(t *testing.T) {
	testCases := []struct {
		desc string
		dbname string
		valid bool
	}{
		{desc: "Test dbname from string", dbname: "testdbname", valid: true},
		/* empty dbname should definitely fail */
		{desc: "Test empty dbname from string", dbname: "", valid: false},
	}
	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			config := GetTestConfig().DBName(test.dbname)
			if test.valid && config.dbname != test.dbname {
				t.Errorf("Failed to set dbname using builder function")
			}
			if !test.valid && config.dbname == test.dbname {
				t.Errorf("Override of known invalid dbname made")
			}
		})
	}
}

func TestPassword(t *testing.T) {
	testCases := []struct {
		desc string
		password string
		valid bool
	}{
		{desc: "Test password from string", password: "testpassword", valid: true},
		/*
			Might be required to override with a blank password, so valid.
			Can't imagine it would be a good idea though
		*/
		{desc: "Test empty password from string", password: "", valid: true},
	}
	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			config := GetTestConfig().Password(test.password)
			if test.valid && config.password != test.password {
				t.Errorf("Failed to set password using builder function")
			}
			if !test.valid && config.password == test.password {
				t.Errorf("Override of known invalid password made")
			}
		})
	}
}

func TestFromEnv(t *testing.T) {
	testCases := []struct {
		desc string
		host string
		hostshouldupdate bool
		port string
		portshouldupdate bool
		user string
		usershouldupdate bool
		dbname string
		dbnameshouldupdate bool
		password string
		passwordshouldupdate bool
	}{
		{
			desc: "Test full config from env",
			host: "hostfrom env",
			hostshouldupdate: true,
			port:"10001",
			portshouldupdate: true,
			user: "fromenv",
			usershouldupdate: true,
			dbname: "dbfromenv",
			dbnameshouldupdate: true,
			password: "passwordfromenv",
			passwordshouldupdate: true
		},
		{desc: "Test config from env", host: "hostfrom env", port:"10001", user: "fromenv", dbname: "dbfromenv", password: "passwordfromenv", valid: true},
		/*
			Might be required to override with a blank password, so valid.
			Can't imagine it would be a good idea though
		*/
		{desc: "Test empty password from string", password: "", valid: true},
	}
	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			config := GetTestConfig().Password(test.password)
			if test.valid && config.password != test.password {
				t.Errorf("Failed to set password using builder function")
			}
			if !test.valid && config.password == test.password {
				t.Errorf("Override of known invalid password made")
			}
		})
	}
}



/*
TODO
Test
	(config PGConfig) Connect() (ds DataStore, err error)
	(config PGConfig) FromEnv() PGConfig

*/
