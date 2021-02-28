package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/markstanden/authentication/argon"
)

func main() {
	http.HandleFunc("/", signin)
	http.ListenAndServe(":8080", nil)
}

func signin(w http.ResponseWriter, r *http.Request) {
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
		err := r.ParseForm()
		if err != nil {
			log.Panic("Failed to parseform")
		}

		fmt.Printf(`
		Password: %s,
		Hash: %s,
		`, r.PostForm.Get("password"), argon.Encode(r.PostForm.Get("password")))
	}

}
