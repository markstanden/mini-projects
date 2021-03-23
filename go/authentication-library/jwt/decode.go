package token

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

// The allowable characters in a URL flavoured base64 encoded string.
// shamelessly stolen from the base64 package, as it is not exported.
const encodeURL = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"

// Decode turns a signed JWT into a map[string]interface (or returns an error)
// but only after checking the validity of the token.
func Decode(untrustedB64 string, secret string, trusted *Token) (err error) {
	//func Decode(untrustedB64 string, secret string) (trusted map[string]interface{}, err error) {

	// untrustedB64 should be a two or three base64 URL encoded strings, separated by dots
	// check that first

	// Verify that the JWT contains at least one period ('.') character.
	if !strings.Contains(untrustedB64, ".") {
		return fmt.Errorf("invalid token, no dots")
	}

	// Split the dot separated base64 URL encoded string into 2 or 3 segments depending on whether it contains a signature
	splitUntrustedB64 := strings.SplitAfter(untrustedB64, ".")

	// If there is 1 or less, or more than 3 sections, something has gone wrong, so exit early.
	if len(splitUntrustedB64) <= 1 || len(splitUntrustedB64) > 3 {
		return fmt.Errorf("token split error, incorrect number of sections")
	}

	// If there is only 2 sections, or the signature section is empty, exit early
	if len(splitUntrustedB64) == 2 || len(splitUntrustedB64) == 3 && splitUntrustedB64[2] == "" {
		return fmt.Errorf("invalid signature section")
	}

	// create the structures to hold the different versions of the 3 parts of the token
	var untrustedValid [][]byte
	var untrustedValidB64 = make([]string, 3)

	for i, split := range splitUntrustedB64 {

		// Trim the dot from the ends of the split base64 URL encoded string
		// and add the cleaned up version to the base64 []string
		splitTrimmed := strings.TrimSuffix(string(split), ".")
		untrustedValidB64[i] = splitTrimmed

		// First check for invalid characters in the input string
		if !checkValidity(splitTrimmed, encodeURL) {
			return fmt.Errorf("header contains invalid chars")
		}

		// Check that the header, payload, and signature are actually valid base64 strings
		// if so decode and assign for later use in the [][]byte
		decoded, err := decodeToString(splitTrimmed)
		if err != nil {
			return err
		}
		// Add the trimmed, validated, decoded []byte to the slice
		untrustedValid = append(untrustedValid, decoded)
	}

	// now the data has had basic validation, recreate the body of the jwt to be tested
	toEncode := []byte(untrustedValidB64[0] + "." + untrustedValidB64[1])

	// prep the hash
	h := hmac.New(sha512.New, []byte(secret))
	// hash and write the jwt body
	h.Write(toEncode)
	testBytes := h.Sum(nil)

	// check that the hashed body is equal to the decoded signature supplied by the jwt
	if hmac.Equal(testBytes, untrustedValid[2]) {
	} else {
		return fmt.Errorf("signature invalid")
	}

	// Now we are satisfied the token is valid, we can extract the data

	err = json.Unmarshal(untrustedValid[1], &trusted)
	if err != nil {
		return fmt.Errorf("error unmarshalling json: \n%v", err)
	}

	//fmt.Println(trusted)

	//fmt.Printf("%T", trusted)

	// the date fields (iat, nbf, exp) will unmarshall as float64
	// convert to int64, and test validity
	//fields := []string{"iat", "nbf", "exp"}
	//for _, f := range fields {
	//	if _, ok := trusted[f]; ok {
	//		trusted[f] = int64(trusted[f].(float64))
	//	}
	//}

	//now := getUnixTime(time.Now())

	// check to see if the token has expired, if it exists
	//if exp, ok := trusted["exp"]; ok {
	//	if exp.(int64) < now {
	//		return nil, fmt.Errorf("token has expired")
	//	}

	//} else {
	//	return nil, fmt.Errorf("token has no expiry date")
	//}

	// check to see if the token has a valid creation date
	//if iat, ok := trusted["iat"]; ok {
	//	if iat.(int64) > now || iat.(int64) < getUnixTime(time.Now().AddDate(-1, 0, 0)) {
	//		return nil, fmt.Errorf("token creation date invalid")
	//	}
	//}

	return nil
}

// decode checks the validity of the supplied string and if valid decodes to a []byte.
func decodeToString(untrusted string) (valid []byte, err error) {

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
