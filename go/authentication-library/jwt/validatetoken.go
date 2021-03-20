package token

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"strings"
	"time"
)

// Decode checks the validity of the jwt token,
// returns the unique identifer or an error
func Validate(untrustedB64 string, secret string) (data map[string]string, err error) {

	// untrustedB64 should be a two or three base64 URL encoded strings, separated by dots

	// Verify that the JWT contains at least one period ('.') character.
	if !strings.Contains(untrustedB64, ".") {
		return nil, fmt.Errorf("invalid token")
	}

	// Split the dot separated base64 URL encoded string into 2 or 3 segments depending on whether it contains a signature
	splitUntrustedB64 := strings.SplitAfter(untrustedB64, ".")

	// If there is 1 or less, or more than 3 sections, something has gone wrong, so exit early.
	if len(splitUntrustedB64) <= 1 || len(splitUntrustedB64) > 3 {
		return nil, fmt.Errorf("token split error, incorrect number of sections")
	}

	// Trim the dot from the ends of the split base64 URL encoded string
	headerUntrustedB64 := strings.TrimSuffix(splitUntrustedB64[0], ".")
	payloadUntrustedB64 := strings.TrimSuffix(splitUntrustedB64[1], ".")
	signatureUntrustedB64 := strings.TrimSuffix(splitUntrustedB64[2], ".")

	// re-encode the header + payload.
	// This is the costly part so check if it's worth doing first

	// Check that the header, payload, and signature are actually valid base64 strings
	// decode and assign for later use.
	var (
		headerUntrustedBytes, payloadUntrustedBytes, signatureUntrustedBytes []byte
	)

	headerUntrustedBytes, err = decode(headerUntrustedB64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode JWT header :\n%v", err)
	}

	payloadUntrustedBytes, err = decode(payloadUntrustedB64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode JWT payload :\n%v", err)
	}

	signatureUntrustedBytes, err = decode(signatureUntrustedB64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode JWT signature :\n%v", err)
	}

	// create the slice of bytes to encode
	// {json header} + "." + {json payload}
	toEncode := []byte(headerUntrustedB64 + "." + payloadUntrustedB64)

	h := hmac.New(sha512.New, []byte(secret))
	h.Write(toEncode)
	testBytes := h.Sum(nil)

	if hmac.Equal(testBytes, signatureUntrustedBytes) {
		fmt.Println("Signature Verified")
	} else {
		fmt.Println("Signature invalid")
	}

	// Let the Encoded JOSE Header be the portion of the JWT before the first period ('.') character.

	fmt.Println(string(headerUntrustedBytes))
	fmt.Println(string(payloadUntrustedBytes))
	return data, nil
}

func getTime() time.Time {
	return time.Now()
}

func format(t time.Time) (string, error) {
	sb, err := t.UTC().MarshalText()
	if err != nil {
		return "", err
	}
	return string(sb), nil
}

// decode checks the validity of the supplied string and if valid decodes to a []byte.
func decode(untrusted string) (valid []byte, err error) {

	// The allowable characters in a URL flavoured base64 encoded string.
	// shamelessly stolen from the base64 package, as it is not exported.
	const encodeURL = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"

	// First check for invalid characters in the input string
	if !checkValidity(untrusted, encodeURL) {
		return nil, fmt.Errorf("header contains invalid chars")
	}
	valid, err = base64.RawURLEncoding.Strict().DecodeString(untrusted)
	if err != nil {
		return nil, fmt.Errorf("failed to decode JWT header :\n%v", err)
	}
	return valid, nil
}

// checkValidity checks whether the supplied string contains runes outside of the valid runes set supplied
func checkValidity(toCheck, validRunes string) (valid bool) {
	for _, r := range toCheck {
		valid = strings.ContainsRune(validRunes, r)
		if !valid {
			// invalid rune found
			return false
		}
	}
	return true
}

/*

   4.   Verify that the resulting octet sequence is a UTF-8-encoded
        representation of a completely valid JSON object conforming to
        RFC 7159 [RFC7159]; let the JOSE Header be this JSON object.

   5.   Verify that the resulting JOSE Header includes only parameters
        and values whose syntax and semantics are both understood and
        supported or that are specified as being ignored when not
        understood.

   6.   Determine whether the JWT is a JWS or a JWE using any of the
        methods described in Section 9 of [JWE].

	7.   Depending upon whether the JWT is a JWS or JWE, there are two
        cases:

        *  If the JWT is a JWS, follow the steps specified in [JWS] for
           validating a JWS.  Let the Message be the result of base64url
           decoding the JWS Payload.

        *  Else, if the JWT is a JWE, follow the steps specified in
           [JWE] for validating a JWE.  Let the Message be the resulting
           plaintext.

   8.   If the JOSE Header contains a "cty" (content type) value of
        "JWT", then the Message is a JWT that was the subject of nested
        signing or encryption operations.  In this case, return to Step
        1, using the Message as the JWT.

   9.   Otherwise, base64url decode the Message following the
        restriction that no line breaks, whitespace, or other additional
        characters have been used.

   10.  Verify that the resulting octet sequence is a UTF-8-encoded
        representation of a completely valid JSON object conforming to
        RFC 7159 [RFC7159]; let the JWT Claims Set be this JSON object.

*/
