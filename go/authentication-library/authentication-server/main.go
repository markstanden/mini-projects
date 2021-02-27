package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", signin)
	http.ListenAndServe(":8080", nil)
}

func signin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("type", "html")
	fmt.Fprintln(w, `
		<h1> Homepage </h1>
		<form method="post">
			<label for="name">Name:</label>
			<input id="name" name="name" type="text" />
			<label for="email">Email</label>
			<input id="email" name="email" type="email" />
			<input value="Submit Info" type="submit" />
		</form>
	`)

}

func hashPassword(ptPassword string) []byte {
	scrypt.	
}


