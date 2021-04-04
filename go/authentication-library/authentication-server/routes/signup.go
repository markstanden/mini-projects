package routes

import (
	"fmt"
	"log"
	"net/http"

	"github.com/markstanden/argonhasher"
	"github.com/markstanden/authentication"
	"github.com/markstanden/securerandom"
)

// SignUp produces the signup route
func SignUp(us authentication.UserDataStore, ts authentication.AccessTokenService) http.Handler {
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
			if r.PostForm.Get("name") == "" || r.PostForm.Get("email") == "" || r.PostForm.Get("password") == "" || r.PostForm.Get("confirmpassword") == "" {
				log.Println("empty field(s) in form")
				http.Redirect(w, r, "/signup", http.StatusSeeOther)
				return
			}
			// Check that the passwords match
			if r.PostForm.Get("password") != r.PostForm.Get("confirmpassword") {
				log.Println("Passwords do not match, Cannot create account")
				http.Redirect(w, r, "/signup", http.StatusSeeOther)
				return
			}

			// hash the password
			passwordHash := argonhasher.Encode(r.PostForm.Get("password"), 0)
			if passwordHash == "" {
				log.Println("failed to create hash: ", err)
			}

			err = us.Add(&authentication.User{
				Name:           r.PostForm.Get("name"),
				Email:          r.PostForm.Get("email"),
				HashedPassword: passwordHash,
				TokenUserID:    securerandom.String(32),
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

			jwt, jwtid, err := ts.Create(userCheck.TokenUserID)
			if err != nil {
				fmt.Fprintf(w, "/routes/signup: error creating jwt\n%v\nError:\n%v", jwt, err.Error())
			}
			log.Println("/routes/signup:\n\tEncoded jwt: \n", jwt, "\n\tjwtid (From Create):\n\t", jwtid)

			userCheck.TokenUserID = jwtid
			// update user table
			//.....

			uid, jwtid, err := ts.Decode(jwt)
			if err != nil {
				log.Printf("/routes/signup: error decoding jwt\n%v\nError:\n%v", jwt, err.Error())
				http.Redirect(w, r, "/", http.StatusSeeOther)
			}

			log.Println("/routes/signup: Decoded JWT OK...")
			log.Println("\n/routes/signup:\nUserID (Decoded from JWT):\n", uid)
			log.Println("\n/routes/signup:\njwtid (Decoded from JWT):\n", jwtid)

			u, err := us.Find("tokenid", uid)
			if err != nil {
				log.Printf("/routes/signup: error looking up user\n%v\nError:\n%v", jwt, err.Error())
				return
			}
			log.Printf("/routes/signup: created and decoded jwt.\nJWT String:\n%v\nUserData:\n%v", jwt, u)
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}

	})
}
