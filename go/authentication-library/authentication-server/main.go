package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/markstanden/argonhasher"
	"github.com/markstanden/jwt"
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
			log.Println("Failed to parseform: ", err)

		}

		compareOK := false

		hash, err := argonhasher.Encode(r.PostForm.Get("password"))
		if err != nil {
			log.Println("failed to create hash: ", err)
		}

		err = argonhasher.Compare(r.PostForm.Get("password"), hash)
		if err != nil {
			log.Println("failed to make comparison: ", err)
		} else {
			compareOK = true
		}

		// Output the password, and the hash, and the result of the comparison.
		fmt.Printf(`
		Password: %s,
		Hash: %s,
		compare ok?: %v,
		`, r.PostForm.Get("password"), hash, compareOK)
		
		// Create a JWT
		token := jwt.New()
		token.Payload.JTI = hash
		token.Encode()

		fmt.Println(token.Decode())
		// 
	}

}
