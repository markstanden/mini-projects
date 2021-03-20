package main

import (
	"fmt"
	"io"
	"os"

	"github.com/markstanden/token"
)

//var jwtDotIO = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.5YLJzijebmz3VWxPTZugJh1JA-q60WOmsEeSl8Ra55cqAdY7wq0mahPDS3U_912j0LRGL7LEEeO1c57muxgTzg"
var secret = "supersecret"

func main() {
	if err := run(os.Args, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(args []string, stdout io.Writer) error {
	jwt, err := token.NewToken(secret, "My Server", "User Idenfier", "My Website", "Unique ID for the JWT", "Key ID (to identify secret version used to encrypt signature")
	if err != nil {
		fmt.Println("Error: ", err)
	}
	fmt.Println("Created Token: ", jwt)

	data, err := token.Decode(jwt, secret)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	fmt.Println("Raw Data: ", data)
	return nil
}
