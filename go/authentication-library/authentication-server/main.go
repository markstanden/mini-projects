package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", home)
	http.ListenAndServe(":8080", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("type", "html")
	fmt.Fprintln(w, `
		<h1> Homepage </h1>
		<form method="post">
			<input name="name" type="text" />
			<input name="email" type="email" />
			<input value="Submit Info" type="submit" />
		</form>
	`)
}
