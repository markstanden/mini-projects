package postgres

import (
	"os"
	"testing"
)

func TestNewConfig(t *testing.T) {
	config := GetTestConfig()
	if config.host == "" ||
		config.port == "" ||
		config.dbname == "" ||
		config.user == "" ||
		config.password == "" {
		t.Error("failed to set config to defaults")
	}
}

func TestHost(t *testing.T) {
	testCases := []struct {
		desc  string
		host  string
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
		desc  string
		port  string
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
		desc  string
		user  string
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
		desc   string
		dbname string
		valid  bool
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
		desc     string
		password string
		valid    bool
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
		desc                 string
		host                 string
		hostshouldupdate     bool
		port                 string
		portshouldupdate     bool
		user                 string
		usershouldupdate     bool
		dbname               string
		dbnameshouldupdate   bool
		password             string
		passwordshouldupdate bool
	}{
		{
			desc:                 "Test full config from env",
			host:                 "hostfrom env",
			hostshouldupdate:     true,
			port:                 "10001",
			portshouldupdate:     true,
			user:                 "fromenv",
			usershouldupdate:     true,
			dbname:               "dbfromenv",
			dbnameshouldupdate:   true,
			password:             "passwordfromenv",
			passwordshouldupdate: true,
		},
		{
			desc:                 "Just Host",
			host:                 "hostfrom env",
			hostshouldupdate:     true,
			port:                 "",
			portshouldupdate:     false,
			user:                 "",
			usershouldupdate:     false,
			dbname:               "",
			dbnameshouldupdate:   false,
			password:             "",
			passwordshouldupdate: false,
		},
		{
			desc:                 "Just Port",
			host:                 "",
			hostshouldupdate:     false,
			port:                 "10001",
			portshouldupdate:     true,
			user:                 "",
			usershouldupdate:     false,
			dbname:               "",
			dbnameshouldupdate:   false,
			password:             "",
			passwordshouldupdate: false,
		},
		{
			desc:                 "Just User",
			host:                 "",
			hostshouldupdate:     false,
			port:                 "",
			portshouldupdate:     false,
			user:                 "fromenv",
			usershouldupdate:     true,
			dbname:               "",
			dbnameshouldupdate:   false,
			password:             "",
			passwordshouldupdate: false,
		},
		{
			desc:                 "Just dbname",
			host:                 "",
			hostshouldupdate:     false,
			port:                 "",
			portshouldupdate:     false,
			user:                 "",
			usershouldupdate:     false,
			dbname:               "dbfromenv",
			dbnameshouldupdate:   true,
			password:             "",
			passwordshouldupdate: false,
		},
		{
			desc:                 "Just password",
			host:                 "",
			hostshouldupdate:     false,
			port:                 "",
			portshouldupdate:     false,
			user:                 "",
			usershouldupdate:     false,
			dbname:               "",
			dbnameshouldupdate:   false,
			password:             "passwordfromenv",
			passwordshouldupdate: true,
		},
	}
	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			/*
				Take the current values of the env variables and reassign them once the test is complete, as they may be needed.
			*/
			envlist := []string{"PGHOST", "PGPORT", "PGUSER", "PGDATABASE", "PGPASSWORD"}
			for _, name := range envlist {
				value, ok := os.LookupEnv(name)
				if ok {
					defer os.Setenv(name, value)
				}
			}

			if test.host != "" {
				os.Setenv("PGHOST", test.host)
			}
			if test.port != "" {
				os.Setenv("PGPORT", test.port)
			}
			if test.user != "" {
				os.Setenv("PGUSER", test.user)
			}
			if test.dbname != "" {
				os.Setenv("PGDATABASE", test.dbname)
			}
			if test.password != "" {
				os.Setenv("PGPASSWORD", test.password)
			}

			config := GetTestConfig().FromEnv()

			t.Run("Host Test", func(t *testing.T) {
				if test.hostshouldupdate && config.host != test.host {
					t.Errorf("Failed to set host using builder function")
				}
				if !test.hostshouldupdate && config.host == test.host {
					t.Errorf("unexpected override host made")
				}
			})
			t.Run("Port Test", func(t *testing.T) {
				if test.portshouldupdate && config.port != test.port {
					t.Errorf("Failed to set port using builder function")
				}
				if !test.portshouldupdate && config.port == test.port {
					t.Errorf("unexpected override port made")
				}
			})
			t.Run("User Test", func(t *testing.T) {
				if test.usershouldupdate && config.user != test.user {
					t.Errorf("Failed to set user using builder function")
				}
				if !test.usershouldupdate && config.user == test.user {
					t.Errorf("unexpected override user made")
				}
			})
			t.Run("DBName Test", func(t *testing.T) {
				if test.dbnameshouldupdate && config.dbname != test.dbname {
					t.Errorf("Failed to set dbname using builder function")
				}
				if !test.dbnameshouldupdate && config.dbname == test.dbname {
					t.Errorf("unexpected override dbname made")
				}
			})
			t.Run("Password Test", func(t *testing.T) {
				if test.passwordshouldupdate && config.password != test.password {
					t.Errorf("Failed to set password using builder function")
				}
				if !test.passwordshouldupdate && config.password == test.password {
					t.Errorf("unexpected override password made")
				}
			})
		})
	}
}

func TestConnect(t *testing.T) {
	db, err := GetTestConfig().Connect()
	if err != nil {
		t.Errorf("Failed to connect to test database.")
	}

	// err = db.Ping()
	// if err != nil {
	// 	t.Errorf("Test DB Failed Ping test\n%v", err)
	// }
}
