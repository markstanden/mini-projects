package securerandom

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// modified version of example here:
// https://golang.org/pkg/crypto/rand/#example_Read

// ByteSlice produces a slice of random bytes, of a provided length
// used as the base for other outputs
func ByteSlice(length int) (bs []byte, err error) {
	bs = make([]byte, length)
	r, err := rand.Read(bs)
	if err != nil {
		return nil, fmt.Errorf("securerandom/ByteSlice: failed to create secure slice of bytes \n%v", err)
	}
	if r != length {
		return nil, fmt.Errorf("securerandom/ByteSlice: failed to create %v byte string", length)
	}
	return bs, nil
}

// String returns a secure URL compliant string of a provided length
func String(length int) (message string, err error) {
	bs, err := ByteSlice(length)
	if err != nil {
		return "", fmt.Errorf("securerandom/String: failed to create secure slice of bytes \n%v", err)
	}
	// encode to URL encoded base 64 string
	// so it can be transmitted in a URL if required
	message = base64.URLEncoding.EncodeToString(bs)
	return message, nil
}
