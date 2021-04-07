package routes

import (
	"fmt"
	"net/http"

	"github.com/markstanden/authentication"
)

/*
	** SignUp **
	SignUp returns the http.Handler for the signup route.
	It is intended as the primary way to create a new user account.
	For convenience it should log in the user and redirect to the users secure area.
*/
func SignUp(us authentication.UserService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "GET" {

			w.Header().Set("type", "html")
			fmt.Fprintln(w, getHTMLFORM("Sign Up here for a new account!"))
		}
		if r.Method == "POST" {
			// Parse the form data in the response
			err := r.ParseForm()
			// Check for errors in the parsing
			if err != nil {
				fmt.Println(w, "Failed to Parse Form: ", err)
			}

			// Check the form data
			if r.PostForm.Get("name") == "" || r.PostForm.Get("email") == "" || r.PostForm.Get("password") == "" || r.PostForm.Get("confirmpassword") == "" {
				fmt.Fprintln(w, getHTMLFORM("Empty field(s) in form"))
				return
			}
			// Check that the passwords match
			if r.PostForm.Get("password") != r.PostForm.Get("confirmpassword") {
				fmt.Fprintln(w, getHTMLFORM("Passwords do not match, Cannot create account"))
				return
			}

			/*
				Create the user and add to the datastore
			*/
			u, err := us.NewUser(r.PostForm.Get("name"), r.PostForm.Get("email"), r.PostForm.Get("password"))
			if err != nil {
				fmt.Fprintln(w, getHTMLFORM("Failed to create user, please try again"))
				return

				/* email verification goes here? */

			}
			fmt.Fprintf(w, `
			Welcome %v,
			DB ID:
			%v
			Email:
			%v
			TokenUserID
			%v
			Hashed Password:
			%v
				`, u.Name, u.UniqueID, u.Email, u.TokenUserID, u.HashedPassword)
		}
	})
}

/*
	** getHTMLForm **
	Internal, temporary function to create and display the user input form
	The subtitle parameter is being used as a status display line, so the form
	can be reissued in event of an error in invalid input.
*/
func getHTMLFORM(subtitle string) (html string) {
	html = fmt.Sprintf(`
		<div style="padding-left:10rem">
			<h1> Sign Up </h1>
			<h3> %v </h3>
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
		</div>
	`, subtitle)
	return
}
