package main

import (
	"fmt"
	"log"
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
	hash, err := generateFromPassword("password123",
		struct {
			memory      uint32
			iterations  uint32
			parallelism uint8
			saltLength  uint32
			keyLength   uint32
		}{
			memory:      64 * 1024,
			iterations:  3,
			parallelism: 2,
			saltLength:  16,
			keyLength:   32,
		})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(hash)
}
