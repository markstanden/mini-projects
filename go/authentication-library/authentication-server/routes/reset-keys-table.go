package routes

import (
	"fmt"
	"log"
	"net/http"

	"github.com/markstanden/authentication"
)

// CreateUsersTable drops the existing user table if it exists and creates a
// new fresh table. Obviously this is for use in development only, to allow quick changes to
// the database structure / authentication.User struct object.
func ResetKeysTable(ss authentication.SecretDataStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "GET" {

			w.Header().Set("type", "html")
			fmt.Fprintln(w, `
				<h1> Game Over </h1>
				<form action="/reset-keys-table" method="POST">
					<input value="Reset the Keys table and start again?" type="submit" />
				</form>
			`)
		}
		if r.Method == "POST" {
			log.Println("authentication/routes: Request to drop table received")
			w.Header().Set("type", "html")
			if err := ss.FullReset(); err != nil {
				log.Println("authentication/routes: Failed to Reset keys table database:\n\t", err)
			}

			// Table has been dropped and recreated, so log and redirect to root route.
			log.Println("authentication/routes: Keys table dropped and created ok")
			http.Redirect(w, r, "/", http.StatusSeeOther)

		}

	})
}
