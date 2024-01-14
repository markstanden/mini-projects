package securerandom

import (
	"crypto/rand"
	"encoding/base64"
	"math"
)

// modified version of example here:
// https://golang.org/pkg/crypto/rand/#example_Read

// ByteSlice produces a slice of random bytes, of a provided length
// also used as the base for other outputs
// returns nil if it fails to make the []byte,
// or if the created slice is the wrong size
func ByteSlice(length uint) (bs []byte) {

	bs = make([]byte, length)
	r, err := rand.Read(bs)
	if err != nil {
		return nil
	}
	if r != int(length) {
		return nil
	}
	return bs
}

// String returns a secure URL compliant string of a provided length
// returns a zero value string if it fails to create the message
func String(length uint) (message string) {

	// base64 uses 4 characters to represent 3 bytes so convert
	b := ((3 * float64(length)) / 4)
	requiredBytes := math.Ceil(b)

	bs := ByteSlice(uint(requiredBytes))

	if len(bs) == 0 {
		return ""
	}

	// encode to URL encoded base 64 string
	// so it can be transmitted in a URL if required
	message = base64.RawURLEncoding.EncodeToString(bs)

	// since every byte produces 1.333 characters
	// we get extra characters in our string if the length doesn't cleanly divide
	// so trim the end
	return message[:length]
}
