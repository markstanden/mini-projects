package token

import (
	"strings"
	"testing"
)

type test []string

func TestDecode(t *testing.T) {

	// create a JWT using data previously parsed on the jwt.io website
	j := jwtioStruct
	validFor := j.ExpirationTime - j.IssuedAtTime
	token := NewToken(j.Issuer, j.Audience, j.UserID, j.JwtID, j.KeyID, validFor)
	jwt, err := token.CreateJWT(jwtioSecret())
	if err != nil {
		t.Fail()
	}

	testsShouldError := []test{

		// test input string with only one section (no dot seperators)
		{strings.ReplaceAll(jwt, ".", ""),
			"Failed to return an error for a token with a single section"},

		// test input string with too many sections
		{strings.Repeat(jwt, 2),
			"Failed to return an error for a token with 4 sections"},

		// test input string with invalid characters
		{strings.ReplaceAll(jwt, string(jwt[0]), "^"),
			"Failed to return an error for a token with invalid base64 URL encoded characters"},

		// test input string with no signature
		{strings.Join(strings.SplitN(jwt, ".", 2), ""),
			"Failed to return an error for a token without a signature"},

		// test input string with invalid signature
		{jwtioToken[:len(jwtioToken)-1],
			"Failed to return an error for a token with an invalid signature"},

		// test input string with "alg" set to "none" **Known Exploit**
		{"eyJ0eXAiOiJKV1QiLCJhbGciOiJub25lIn0.eyJpZCI6MSwiaWF0IjoxNTczMzU4Mzk2fQ.",
			"Failed to return an error for a alg:none token"},

		// test input string with "alg" set to "HS256" or other algorithm
		{"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.6v6L_1zpfV-EUf-q7j9peL15ep8nsyxNnUCTHjnrSes",
			"Failed to return an error for alg:HS256"},

		{"eyJhbGciOiJIUzM4NCIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.7siIdaH0T0VCvfRWRI_lEgFWrK7tIjyWMQAYtQ72Qz5El0vbVMKAGtIPlJE2mOvT",
			"Failed to return an error for alg:HS384"},

		// test input string with "exp" set to the past
		{"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IlRlc3R5IE1jVGVzdGZhY2UiLCJpYXQiOjEwMDAwMDAwMDEsImV4cCI6MTAwMDAyMjIyMn0.1yts0Ifi921rDefqEV7rv-JZnJ33vBQBdHlhYJS4kl4KXzZyu5WOieBzbr04W1gYIP99oZ6QvIrPCbztRzPHBQ",
			"Failed to return an error for an expired token"},

		// test input string with "iat" created in the future
		{"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IlRlc3R5IE1jVGVzdGZhY2UiLCJpYXQiOjIwMDAwMDAwMDEsImV4cCI6MjAwMDAyMjIyMn0.XlF_dhnLt7tLc73_v7_T8LtFVVQMgJV3vy6tm9VoWF1BJVrit1CkyOlMzBQ2uz0iQs1Ggm7bw7JEK1l-t994Cw",
			"Failed to return an error for an issued at time from the future"},

		// test input string with "nbf" valid from the future
		{"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IlRlc3R5IE1jVGVzdGZhY2UiLCJpYXQiOjEwMDAwMDAwMDEsIm5iZiI6MjAwMDAwMDAwMSwiZXhwIjoyMDAwMDIyMjIyfQ.3iCZdcj4JC7Kbh8uG2GBmlRUCUkHNGPzLCc41ctZPn0mPx1As4XsEiDU4mzknCgEYUCJJ4NAuvea943RUrMOvw",
			"Failed to return an error for a token with a not before time from the future"},
	}

	for i, test := range testsShouldError {
		// create an empty Token struct to put the data in
		got := Token{}

		// run the decoder for this test
		err := Decode(test[0], jwtioSecret(), &got)

		// if the struct is not empty OR no errors
		if (got != Token{}) || err == nil {
			t.Error("test ", i, ":\n", test[1])
		}
	}
}
