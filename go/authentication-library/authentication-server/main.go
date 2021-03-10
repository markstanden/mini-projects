package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/markstanden/argonhasher"
	"github.com/markstanden/authentication/cache"
	"github.com/markstanden/authentication/config"

	//"github.com/markstanden/authentication/http"
	"github.com/markstanden/authentication/postgres"
)

var c *cache.UserCache

func main() {
	if err := run(os.Args, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(args []string, stdout io.Writer) error {

	// open a connection to the database
	db := postgres.NewConnection(
		config.PGhost,
		config.PGusername,
		config.PGpassword,
		config.PGdatabaseName,
		config.PGport,
	)

	// check the database connection is up and running
	err := db.DB.Ping()
	if err != nil {
		fmt.Println("Connection Failure", err)
	}
	// Close the database when the server ends
	defer db.DB.Close()

	// Create a user cache
	c = cache.NewUserCache(db)

	// Create a handler for our routes
	http.HandleFunc("/signin", signin)
	http.ListenAndServe(":8080", nil)
	return nil
}

// SignIn produces the signin route
func signin(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {

		w.Header().Set("type", "html")
		fmt.Fprintln(w, `
		<h1> Homepage </h1>
		<form method="post">
			<label for="email">Email:</label>
			<input id="email" name="email" type="email" /><br>
			<label for="password">Password</label>
			<input id="password" name="password" type="password" /><br>
			<input value="Submit Info" type="submit" />
		</form>
	`)
	}
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			log.Println("Failed to parseform: ", err)

		}

		compareOK := false

		hash, err := argonhasher.Encode(r.PostForm.Get("password"))
		if err != nil {
			log.Println("failed to create hash: ", err)
		}

		err = argonhasher.Compare(r.PostForm.Get("password"), hash)
		if err != nil {
			log.Println("failed to make comparison: ", err)
		} else {
			compareOK = true
		}

		// Output the password, and the hash, and the result of the comparison.
		fmt.Printf(`
		Password: %s,
		Hash: %s,
		compare ok?: %v,
		`, r.PostForm.Get("password"), hash, compareOK)

		user, err := c.FindByEmail(r.PostForm.Get("email"))
		if err != nil {
			log.Println("User Not Found")
		}
		fmt.Fprintln(w, user)
	}

	// Create a JWT
	//token := jwt.New()
	//token.Payload.JTI = hash
	//token.Encode()

	//fmt.Println(token.Decode())
	//
}
