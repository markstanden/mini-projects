package routes

import (
	"fmt"
	"log"
	"net/http"

	"github.com/markstanden/argonhasher"
	"github.com/markstanden/authentication"
)

// SignUp produces the signup route
func SignUp(us authentication.UserService, ts authentication.TokenService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "GET" {

			w.Header().Set("type", "html")
			fmt.Fprintln(w, `
		<h1> Sign Up for a new StandenSoft Account </h1>
		<form action="/signup" method="POST">
			<label for="name">Name:</label>
			<input id="name" name="name" type="text" maxlength="255"/><br>
			<label for="email">Email:</label>
			<input id="email" name="email" type="email" maxlength="255"/><br>
			<label for="password">Password</label>
			<input id="password" name="password" type="password" maxlength="255"/><br>
			<label for="confirmpassword">Confirm Password</label>
			<input id="confirmpassword" name="confirmpassword" type="password" /><br>
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

			// Check the form data

			// Check that the passwords match
			if r.PostForm.Get("password") != r.PostForm.Get("confirmpassword") {
				log.Println("Passwords do not match, Cannot create account")
				http.Redirect(w, r, "/signup", http.StatusSeeOther)
				return
			}

			// hash the password
			passwordHash, err := argonhasher.Encode(r.PostForm.Get("password"))
			if err != nil {
				log.Println("failed to create hash: ", err)
			}

			idkey := r.PostForm.Get("name") + r.PostForm.Get("email")
			idHash, err := argonhasher.Encode(idkey)
			if err != nil {
				log.Println("failed to create token hash: ", err)
			}

			err = us.Add(&authentication.User{
				Name:           r.PostForm.Get("name"),
				Email:          r.PostForm.Get("email"),
				HashedPassword: passwordHash,
				Token:          idHash,
			})
			if err != nil {
				log.Println("failed to create user account :\n", err)
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
			log.Println("User Account Created OK")
			userCheck, err := us.Find("email", r.PostForm.Get("email"))
			if err != nil {
				log.Println("failed to lookup created user account :\n", err)
			}
			log.Println("authentication/signup: created user read from store: \n", userCheck)

			t, err := ts.Create(userCheck)
			if err != nil {
				fmt.Fprintf(w, "authentication/signup: error creating jwt\n%v\nError:\n%v", t, err.Error())
			}
			log.Println("authentication/signup: created jwt: \n", t)

			dt, err := ts.Decode(t)
			if err != nil {
				fmt.Fprintf(w, "authentication/signup: error decoding jwt\n%v\nError:\n%v", t, err.Error())
			}
			log.Println("authentication/signup: decoded jwt: \n", dt)
			fmt.Fprintf(w, "authentication/signup: created and decoded jwt\n%v\nData:\n%v", t, dt)

		}
	})
}
