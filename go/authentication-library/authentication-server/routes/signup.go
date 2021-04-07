package routes

import (
	"fmt"
	"net/http"

	"github.com/markstanden/authentication"
)

// SignUp produces the signup route
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

			u, err := us.NewUser(r.PostForm.Get("name"), r.PostForm.Get("email"), r.PostForm.Get("password"))
			if err != nil {
				fmt.Fprintln(w, getHTMLFORM("Failed to create user, please try again"))
				return
			}
			fmt.Fprintf(w, "Welcome:\t\t%v\n\nDB ID:\t\t%v\nEmail:\t\t%vTokenUserID\t\t%v\nHashed Password:\t%v", u.Name, u.UniqueID, u.Email, u.TokenUserID, u.HashedPassword)
		}
	})
}

func getHTMLFORM(subtitle string) (html string) {
	html = fmt.Sprintf(`
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
	`, subtitle)
	return
}
