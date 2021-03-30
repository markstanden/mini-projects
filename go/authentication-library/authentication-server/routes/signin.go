package routes

import (
	"fmt"
	"log"
	"net/http"

	"github.com/markstanden/argonhasher"
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
			}

			user, err := us.Find("email", r.PostForm.Get("email"))
			if err != nil {
				fmt.Fprintln(w, getHTML("Sign In - Invalid UserName", r.PostForm.Get("email")))
			}

			// Initialise a boolean variable that hold whether the password matches the stored, hashed password.
			compareOK := false
			valid := argonhasher.Confirm(r.PostForm.Get("password"), user.HashedPassword)
			if !valid {
				fmt.Fprintln(w, getHTML("Sign In - Invalid Password", r.PostForm.Get("email")))
			} else {
				compareOK = true
			}

			if compareOK {
				log.Println("User Account Logged In OK")

				fmt.Fprintf(w, `
			User Account Details:
			ID: %v
			Name: %v
			Email: %v
			TokenID: %v
			Error: %v
			`, user.UniqueID, user.Name, user.Email, user.TokenID, err)
			}
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
