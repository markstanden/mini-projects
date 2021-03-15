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
func CreateUsersTable(us authentication.UserService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "GET" {

			w.Header().Set("type", "html")
			fmt.Fprintln(w, `
		<h1> Game Over </h1>
		<form action="/create-users-table" method="POST">
			<input value="Reset the Users table and start again?" type="submit" />
		</form>
	`)
		}

		if r.Method == "POST" {
			w.Header().Set("type", "html")
			if err := us.FullReset(); err != nil {
				log.Println("authentication/routes: Failed to create Users table database:\n\t", err)
			}

			// Table has been dropped and recreated, so log and redirect to root route.
			log.Println("authentication/routes: Table dropped and created ok")
			http.Redirect(w, r, "/", http.StatusSeeOther)

		}

	})
}
