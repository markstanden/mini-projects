package b64

import (
	"encoding/base64"
	"errors"
)

/*
	ErrNotUrlB64 is issued when input
	string is a non compliant b64 string
*/
var ErrNotUrlB64 error = errors.New("provided string is not URL encoded b64")

/*
	FromBytes converts a []byte into a URL encoded base64 string
*/
func FromBytes(bs []byte) (b64 string) {
	return base64.RawURLEncoding.EncodeToString(bs)
}

/*
	ToBytes checks the validity of the
	supplied string and if valid returns a []byte.
	Returns ErrNotUrlB64 if the string contains invalid characters
*/
func ToBytes(b64 string) (bs []byte, err error) {

	bs, err = base64.RawURLEncoding.Strict().DecodeString(b64)
	if err != nil {
		return nil, ErrNotUrlB64
	}
	return bs, nil
}
