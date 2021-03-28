package securerandom

import (
	"crypto/hmac"
	"crypto/sha512"
	"fmt"
	"log"
	"testing"
)

func TestByteSlice(t *testing.T) {

	tests := []uint{512, 256, 128, 64, 32, 16, 8, 4, 2, 1, 0}

	for _, l := range tests {
		bs := ByteSlice(l)
		if bs == nil {
			t.Errorf("failed to create secure byteslice")
		}
		if len(bs) != int(l) {
			t.Errorf("incorrect length: wanted %d, got %d", l, len(bs))
		}
	}
}

func TestString(t *testing.T) {

	tests := []uint{512, 256, 128, 64, 32, 16, 8, 4, 31, 23, 17, 13, 11, 7, 5, 3, 2, 1}

	for _, l := range tests {
		s := String(l)
		if s == "" {
			t.Errorf("failed to create secure string")
		}
		if len(s) != int(l) {
			t.Errorf("incorrect length: wanted %d, got %d - %v", l, len(s), s)
		}
	}
}

func ExampleByteSlice() {
	// create a secure slice of bytes for use as a cryto key in a hash

	keylen := 64 //bytes
	pw := []byte("plain text password")

	// ByteSlice requires uint value for length
	key := ByteSlice(uint(keylen))

	// Use the created key to hash the plain text password
	hmac := hmac.New(sha512.New, key)
	hmac.Write(pw)
	hashedPW := hmac.Sum(nil)

	// Clear the plain text password
	pw = nil

	// Use the key and hashedPW
	log.Printf("Password: \n%v\nSecure Key: \n%v\nHashed Password:\n%v", pw, key, hashedPW)
}

func ExampleString() {
	// create a secure base64 encoded string for use as a unique identifer for a user

	type user struct {
		id    string
		name  string
		email string
	}

	// The size of our identifier
	var idLen uint = 128 //characters

	// String requires uint value for length
	key := String(idLen)

	// We now have a 128 character URL safe string to add to our user
	u := user{
		id:    key,
		name:  "Testy McTestface",
		email: "testy@mctestface.com",
	}

	// Use the new user
	fmt.Printf("User ID: \n%v\nName: \n%v\nEmail:\n%v", u.id, u.name, u.email)
}