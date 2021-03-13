package routes

import (
	"fmt"
	"net/http"

	"github.com/markstanden/authentication"
)

// SignIn produces the signin route
func CreateUsersTable(us authentication.UserService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "GET" {

			w.Header().Set("type", "html")
			sqlResult, err := us.Create()
			fmt.Fprintf(w, `
		<h1> Creating Users Table: /nError : %v/nSQL Result: %v</h1>
	`, err, sqlResult)
		}

	})
}
