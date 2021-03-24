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
func Decode(untrustedB64 string, passwordLookup func(key string) (secret string, err error), data *Token) (err error) {
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

	// We need the header and key ID from the payload so we extract the data
	err = json.Unmarshal(untrustedValid[0], &data)
	if err != nil {
		return fmt.Errorf("error unmarshalling json: \n%v", err)
	}

	//fmt.Println(data)

	// check the alg
	if data.Algorithm == "none" {
		*data = Token{}
		return fmt.Errorf(`"alg":"none" specified: %v`, data.Algorithm)
	}

	if data.Algorithm != "HS512" {
		*data = Token{}
		return fmt.Errorf("invalid encoding algorithm specified: %v", data.Algorithm)
	}

	if data.TokenType != "JWT" {
		*data = Token{}
		return fmt.Errorf("invalid token type: %v", data.TokenType)
	}

	// Header ok, unmarshal payload to check time fields and obtain keyID

	err = json.Unmarshal(untrustedValid[1], &data.Payload)
	if err != nil {
		return fmt.Errorf("error unmarshalling json: \n%v", err)
	}

	// current date/time
	now := getUnixTime()

	// check to see if the token has a valid creation date
	if data.IssuedAtTime > now {
		*data = Token{}
		return fmt.Errorf("token creation date is in the future")
	}

	// check to see if the token has expired, if it exists
	if data.ExpirationTime < now {
		*data = Token{}
		return fmt.Errorf("token has expired")
	}

	// check to see if the token has a valid not before date
	if data.NotBeforeTime > now {
		*data = Token{}
		return fmt.Errorf("token not yet operational")
	}

	toEncode := []byte(untrustedValidB64[0] + "." + untrustedValidB64[1])

	secret, err := passwordLookup(data.KeyID)
	if err != nil {
		*data = Token{}
		return fmt.Errorf("failed to extract secret from callback")
	}

	// prep the hash
	h := hmac.New(sha512.New, []byte(secret))
	// hash and write the jwt body
	h.Write(toEncode)
	testBytes := h.Sum(nil)

	// check that the hashed body is equal to the decoded signature supplied by the jwt
	if !hmac.Equal(testBytes, untrustedValid[2]) {
		*data = Token{}
		return fmt.Errorf("signature invalid")
	}

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
