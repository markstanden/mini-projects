package hash

import (
	"crypto/hmac"
	"crypto/sha512"
)

// hash uses HMAC Sha512 to hash the provided message using
// the provided secret.
// the hash is returned as a URL encoded base64 string
func HS512(message, secret string) (hash []byte) {
	hmac := hmac.New(sha512.New, []byte(secret))
	hmac.Write([]byte(message))
	bs := hmac.Sum(nil)
	return bs
}
