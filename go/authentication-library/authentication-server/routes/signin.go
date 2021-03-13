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
			// Parse the form data in the response
			err := r.ParseForm()
			// Check for errors in the parsing
			if err != nil {
				log.Println("Failed to Parse Form: ", err)
			}

			user, err := us.FindByEmail(r.PostForm.Get("email"))
			if err != nil {
				fmt.Fprintf(w, "failed to lookup user account, invalid UserName :/n%v", err)
			}

			// Initialise a boolean variable that hold whether the password matches the stored, hashed password.
			compareOK := false

			err = argonhasher.Compare(r.PostForm.Get("password"), user.HashedPassword)
			if err != nil {
				fmt.Fprintf(w, "Your password is incorrect!\n%v", err)
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
			Token: %v
			Error: %v
			`, user.UniqueID, user.Name, user.Email, user.Token, err)
			}

			// Create a JWT
			//token := jwt.New()
			//token.Payload.JTI = hash
			//token.Encode()

			//fmt.Println(token.Decode())
			//
		}
	})
}
