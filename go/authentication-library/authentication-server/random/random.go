package random

import (
	"crypto/rand"
	"fmt"
)

func String(length int) (string, error) {
	bs := make([]byte, length)
	r, err := rand.Read(bs)
	if err != nil {
		return "", err
	}
	if r != length {
		return "", fmt.Errorf("failed to create %v byte string", length)
	}
	return string(bs), nil
}
