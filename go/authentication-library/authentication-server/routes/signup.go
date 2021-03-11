package routes

import (
	"fmt"
	"log"
	"net/http"

	"github.com/markstanden/argonhasher"
	"github.com/markstanden/authentication"
)

// SignUp produces the signup route
func SignUp(us authentication.UserService) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request){

	if r.Method == "GET" {

		w.Header().Set("type", "html")
		fmt.Fprintln(w, `
		<h1> Sign Up for a new StandenSoft Account </h1>
		<form method="post">
			<label for="name">Email:</label>
			<input id="name" name="name" type="text" /><br>
			<label for="email">Email:</label>
			<input id="email" name="email" type="email" /><br>
			<label for="password">Password</label>
			<input id="password" name="password" type="password" /><br>
			<label for="confirmpassword">Password</label>
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
			http.RedirectHandler("/signin", http.StatusSeeOther)
		}

		// hash the password
		hash, err := argonhasher.Encode(r.PostForm.Get("password"))
		if err != nil {
			log.Println("failed to create hash: ", err)
		}

		fmt.Println(hash)
		// Create a unique token to represent the user
		//t := jwt.NewToken()
	// Create a JWT
	//token := jwt.New()
	//token.Payload.JTI = hash
	//token.Encode()

	//fmt.Println(token.Decode())
	//
	}})
} 