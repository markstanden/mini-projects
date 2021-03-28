package b64

import (
	"encoding/base64"
	"errors"
)

var ErrNotUrlB64 error = errors.New("provided string is not URL encoded b64")

func FromBytes(bs []byte) (b64 string) {
	return base64.RawURLEncoding.Strict().EncodeToString(bs)
}

// decode checks the validity of the supplied string and if valid returns a []byte.
func ToBytes(b64 string) (bs []byte, err error) {

	bs, err = base64.RawURLEncoding.Strict().DecodeString(b64)
	if err != nil {
		return nil, ErrNotUrlB64
	}
	return bs, nil
}
