package token

import (
	"fmt"
	"testing"
)

/*

eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJnaXRodWIuY29tL21hcmtzdGFuZGVuIiwic3ViIjoiMTIzNDU2Nzg5MCIsImF1ZCI6ImdpdGh1Yi5jb20vbWFya3N0YW5kZW4vYXV0aGVudGljYXRpb24iLCJleHAiOjE2NTAwMDAwMDAsIm5iZiI6MTYwMDAwMDAwMCwiaWF0IjoxNjAwMDAwMDAwLCJqdGkiOiJuVjFNMkIyWmwtU0MwNEdhWkpwN3FEcVA0M0duQzFaZ3R0VDBFOGR2aC1qc2VQRjBsNXAwRUVrS01IOHdJejVNMnpsenI1R0wzUi1UODltSy1OUndBUT09Iiwia2lkIjoiTVdJeFNZbl9RZFgybVBGRml3ZnUyTHVzT2lYaWRNUGpEX2lzMEtyNEJLdnZzYmdBQUUyM0xuVmRqSThVQUZXMUZ6LTlMSlBPcUs5TEFueldwWHBRcHc9PSJ9.tbQ5tU9f6TdKPwiftAAwbgst1fpqzT1kBQ2TU2d7ADt9AE632AhXVqSnAxFzET2wt6Nz47MJERCjvPVj_Pe2uQ

{
"iss":"github.com/markstanden",
"sub":"1234567890",
"aud":"github.com/markstanden/authentication",
"exp":1650000000,
"nbf":1600000000,
"iat":1600000000,
"jti":"nV1M2B2Zl-SC04GaZJp7qDqP43GnC1ZgttT0E8dvh-jsePF0l5p0EEkKMH8wIz5M2zlzr5GL3R-T89mK-NRwAQ==",
"kid":"MWIxSYn_QdX2mPFFiwfu2LusOiXidMPjD_is0Kr4BKvvsbgAAE23LnVdjI8UAFW1Fz-9LJPOqK9LAnzWpXpQpw=="
}
*/



type test []string

/* func TestDecode(t *testing.T) {

	testsShouldError := []test{

		// test input string with only one section (no dot seperators)
		{"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IlRlc3R5IE1jVGVzdGZhY2UiLCJpYXQiOjE1MTYyMzkwMjIsImV4cCI6MTUxOTk5OTk5fQ7P9PODLthZjoMYASHZtmLKSYheID6ACLoqEwHL45cX-z5YeGFRIASIbEEEj5hk2vLMeKegkXv5jwL3DcqFxIIg",
			"secretcode",
			"Failed to return an error for a token with a single section"},

		// test input string with too many sections
		{"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IlRlc3R5IE1jVGV.zdGZhY2UiLCJpYXQiOjE1MTYyMzkwMjIsImV4cCI6MTUxOTk5OTk5fQ.7P9PODLthZjoMYASHZtmLKSYheID6ACLoqEwHL45cX-z5YeGFRIASIbEEEj5hk2vLMeKegkXv5jwL3DcqFxIIg",
			"secretcode",
			"Failed to return an error for a token with 4 sections"},

		// test input string with invalid characters
		{"eyJh/GciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IlRlc3R5IE1jVGVzdGZhY2UiLCJpYXQiOjE1MTYyMzkwMjIsImV4cCI6MTUxOTk5OTk5fQ.7P9PODLthZjoMYASHZtmLKSYheID6ACLoqEwHL45cX-z5YeGFRIASIbEEEj5hk2vLMeKegkXv5jwL3DcqFxIIg",
			"secretcode",
			"Failed to return an error for a token with invalid base64 URL encoded characters"},

		// test input string with no signature
		{"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IlRlc3R5IE1jVGVzdGZhY2UiLCJpYXQiOjE1MTYyMzkwMjIsImV4cCI6MTUxOTk5OTk5fQ",
			"secretcode",
			"Failed to return an error for a token without a signature"},

		// test input string with invalid signature
		{"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IlRlc3R5IE1jVGVzdGZhY2UiLCJpYXQiOjE1MTYyMzkwMjIsImV4cCI6MTUxOTk5OTk5fQ.7P9PODLthZjoMYASHZtmLKSYheID6ACLoqEwHL45cX-z5YeGFRIASIbEEEj5hk2MeKegkXv5jwL3DcqFxIIg",
			"secretcode",
			"Failed to return an error for a token with an invalid signature"},

		// test input string with empty secret
		{"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IlRlc3R5IE1jVGVzdGZhY2UiLCJpYXQiOjE1MTYyMzkwMjIsImV4cCI6MTUxOTk5OTk5fQ.7P9PODLthZjoMYASHZtmLKSYheID6ACLoqEwHL45cX-z5YeGFRIASIbEEEj5hk2vLMeKegkXv5jwL3DcqFxIIg",
			"",
			"Failed to return an error for an empty secret"},

		// test input string with incorrect secret
		{"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IlRlc3R5IE1jVGVzdGZhY2UiLCJpYXQiOjE1MTYyMzkwMjIsImV4cCI6MTUxOTk5OTk5fQ.7P9PODLthZjoMYASHZtmLKSYheID6ACLoqEwHL45cX-z5YeGFRIASIbEEEj5hk2vLMeKegkXv5jwL3DcqFxIIg",
			"secretcde",
			"Failed to return an error for an incorrect secret"},

		// test input string with "alg" set to "none" **Known Exploit**
		{"eyJ0eXAiOiJKV1QiLCJhbGciOiJub25lIn0.eyJpZCI6MSwiaWF0IjoxNTczMzU4Mzk2fQ.",
			"",
			"Failed to return an error for a alg:none token"},

		// test input string with "alg" set to "HS256" or other algorithm
		{"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.6v6L_1zpfV-EUf-q7j9peL15ep8nsyxNnUCTHjnrSes",
			"secretcode",
			"Failed to return an error for alg:HS256"},

		{"eyJhbGciOiJIUzM4NCIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.7siIdaH0T0VCvfRWRI_lEgFWrK7tIjyWMQAYtQ72Qz5El0vbVMKAGtIPlJE2mOvT",
			"secretcode",
			"Failed to return an error for alg:HS384"},

		// test input string with "exp" set to the past
		{"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IlRlc3R5IE1jVGVzdGZhY2UiLCJpYXQiOjEwMDAwMDAwMDEsImV4cCI6MTAwMDAyMjIyMn0.1yts0Ifi921rDefqEV7rv-JZnJ33vBQBdHlhYJS4kl4KXzZyu5WOieBzbr04W1gYIP99oZ6QvIrPCbztRzPHBQ",
			"secretcode",
			"Failed to return an error for an expired token"},

		// test input string with "iat" created in the future
		{"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IlRlc3R5IE1jVGVzdGZhY2UiLCJpYXQiOjIwMDAwMDAwMDEsImV4cCI6MjAwMDAyMjIyMn0.XlF_dhnLt7tLc73_v7_T8LtFVVQMgJV3vy6tm9VoWF1BJVrit1CkyOlMzBQ2uz0iQs1Ggm7bw7JEK1l-t994Cw",
			"secretcode",
			"Failed to return an error for an issued at time from the future"},

		// test input string with "nbf" valid from the future
		{"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IlRlc3R5IE1jVGVzdGZhY2UiLCJpYXQiOjEwMDAwMDAwMDEsIm5iZiI6MjAwMDAwMDAwMSwiZXhwIjoyMDAwMDIyMjIyfQ.3iCZdcj4JC7Kbh8uG2GBmlRUCUkHNGPzLCc41ctZPn0mPx1As4XsEiDU4mzknCgEYUCJJ4NAuvea943RUrMOvw",
			"secretcode",
			"Failed to return an error for a token with a not before time from the future"},
	}

	for i, test := range testsShouldError {
		// create an empty Token struct to put the data in
		got := Token{}

		// run the decoder for this test
		err := Decode(test[0], test[1], &got)

		// if the struct is not empty OR no errors
		if (got != Token{}) || err == nil {
			t.Error("test ", i, ":\n", test[2])
		}
	}
} */

func TestDecodeJWTIO(t *testing.T) {
	// test with a valid SHA512 JWT created from jwt.io website
	test := "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJnaXRodWIuY29tL21hcmtzdGFuZGVuIiwic3ViIjoiMTIzNDU2Nzg5MCIsImF1ZCI6ImdpdGh1Yi5jb20vbWFya3N0YW5kZW4vYXV0aGVudGljYXRpb24iLCJleHAiOjE2NTAwMDAwMDAsIm5iZiI6MTYwMDAwMDAwMCwiaWF0IjoxNjAwMDAwMDAwLCJqdGkiOiJuVjFNMkIyWmwtU0MwNEdhWkpwN3FEcVA0M0duQzFaZ3R0VDBFOGR2aC1qc2VQRjBsNXAwRUVrS01IOHdJejVNMnpsenI1R0wzUi1UODltSy1OUndBUT09Iiwia2lkIjoiTVdJeFNZbl9RZFgybVBGRml3ZnUyTHVzT2lYaWRNUGpEX2lzMEtyNEJLdnZzYmdBQUUyM0xuVmRqSThVQUZXMUZ6LTlMSlBPcUs5TEFueldwWHBRcHc9PSJ9.tbQ5tU9f6TdKPwiftAAwbgst1fpqzT1kBQ2TU2d7ADt9AE632AhXVqSnAxFzET2wt6Nz47MJERCjvPVj_Pe2uQ"

	secret := "secretcode"

	var want = Payload{
		Issuer: 		"github.com/markstanden",
		Subject:        "1234567890",
		Audience:       "github.com/markstanden/authentication",
		ExpirationTime: 1650000000,
		NotBeforeTime:  1600000000,
		IssuedAtTime:   1600000000,
		TokenID:        "nV1M2B2Zl-SC04GaZJp7qDqP43GnC1ZgttT0E8dvh-jsePF0l5p0EEkKMH8wIz5M2zlzr5GL3R-T89mK-NRwAQ==",
		KeyID:          "MWIxSYn_QdX2mPFFiwfu2LusOiXidMPjD_is0Kr4BKvvsbgAAE23LnVdjI8UAFW1Fz-9LJPOqK9LAnzWpXpQpw==",
	}

	got := Token{}

	err := Decode(test, secret, &got)
	if err != nil {
		fmt.Printf("failed to decode JWT \n%t", err)
	}
	if (got.Payload != want) {
		t.Errorf("Want: \n%v\nGot: \n%v\n", want, got.Payload)
	}
}
