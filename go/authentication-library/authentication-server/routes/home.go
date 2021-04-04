package routes

import (
	"fmt"
	"net/http"

	"github.com/markstanden/authentication"
)

// SignIn produces the signin route
func Home(us authentication.UserDataStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "GET" {

			w.Header().Set("type", "html")
			fmt.Fprintln(w, `
		<h1> Homepage </h1>
		<a href="/signin">Sign In</a><br>
		<a href="/signup">Sign Up</a>
	`)
		}

	})
}
