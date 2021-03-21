package routes

import (
	"fmt"
	"log"
	"net/http"

	"github.com/markstanden/argonhasher"
	"github.com/markstanden/authentication"
	token "github.com/markstanden/authentication/tokenhandler"
)

// SignUp produces the signup route
func SignUp(us authentication.UserService) http.Handler {
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
			log.Println("User Account Created OK, Looking up user")
			userCheck, err := us.Find("email", r.PostForm.Get("email"))
			if err != nil {
				log.Println("failed to lookup created user account :\n", err)
			}
			/* fmt.Fprintf(w, `
			User Account Created OK:
			ID: %v
			Name: %v
			Email: %v
			Hash: %v
			Token: %v
			Error: %v
			`, userCheck.UniqueID, userCheck.Name, userCheck.Email, userCheck.HashedPassword, userCheck.Token, err)
			*/
			t, err := token.Create(userCheck, "secretcode")
			if err != nil {
				fmt.Fprint(w, err.Error())
			}
			dt, err := token.Decode(t, "secretcode")
			if err != nil {
				fmt.Fprint(w, err.Error())
			}
			fmt.Fprint(w, t, "\n", dt)
			//http.RedirectHandler("/", http.StatusSeeOther)
			// Create a unique token to represent the user
			//t := jwt.NewToken()
			// Create a JWT
			//token := jwt.New()
			//token.Payload.JTI = hash
			//token.Encode()

			//fmt.Println(token.Decode())
			//
		}
	})
}
