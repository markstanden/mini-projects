package routes

import (
	"fmt"
	"log"
	"net/http"

	"github.com/markstanden/authentication"
)

// SignIn produces the signin route
func SignIn(us authentication.UserService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "GET" {

			w.Header().Set("type", "html")
			fmt.Fprintln(w, getHTML("Sign In", ""))
		}
		if r.Method == "POST" {
			// Parse the form data in the response
			err := r.ParseForm()
			// Check for errors in the parsing
			if err != nil {
				log.Println("Failed to Parse Form: ", err)
				return
			}

			u, err := us.Login(r.PostForm.Get("name"), r.PostForm.Get("password"))
			if err != nil {
				fmt.Fprintln(w, getHTML("Invalid Username/Password - Have another go...", ""))
			}

			fmt.Fprintf(w, `
				User Account Details:
				ID: %v
				Name: %v
				Email: %v
				TokenUserID: %v
				Error: %v
				`, u.UniqueID, u.Name, u.Email, u.TokenUserID, err)
		}
	})
}

func getHTML(title, email string) (html string) {
	return fmt.Sprintf(`
		<h1> %v </h1>
		<form action="/signin" method="POST">
			<label for="email">Email:</label>
			<input id="email" name="email" type="email" maxlength="255" placeholder="%v"/><br>
			<label for="password">Password</label>
			<input id="password" name="password" type="password" maxlength="255"/><br>
			<input value="Submit Info" type="submit" />
		</form>`, title, email)
}
