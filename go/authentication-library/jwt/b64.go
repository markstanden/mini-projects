package token

import "encoding/base64"

func encodeBase64(bs []byte) (b64 string) {
	return base64.RawURLEncoding.Strict().EncodeToString(bs)
}

// decode checks the validity of the supplied string and if valid returns a []byte.
func decodeBase64(b64 string) (bs []byte, err error) {

	bs, err = base64.RawURLEncoding.Strict().DecodeString(b64)
	if err != nil {
		return nil, ErrInvalidToken
	}
	return bs, nil
}
