package securerandom

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// modified version of example here:
// https://golang.org/pkg/crypto/rand/#example_Read

func getSecureByteSlice(length int) (bs []byte, err error) {
	bs = make([]byte, length)
	r, err := rand.Read(bs)
	if err != nil {
		return "", err
	}
	if r != length {
		return nil, fmt.Errorf("failed to create %v byte string", length)
	}
	return bs, nil
}

func String(length int) (message string, err error) {
	bs, err := getSecureByteSlice(length)
	if err != nil {
		return "", err
	}
	// encode to URL encoded base 64 string
	// so it can be transmitted in a URL if required
	message = base64.URLEncoding.EncodeToString(bs)
	return message, nil
}
