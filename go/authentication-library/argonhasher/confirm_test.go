package argonhasher

import "fmt"

func ExampleConfirm() {
	/*
		the plain text password hashed in the Encode example
	*/
	pw := "c2BDNoW38DStXvzP"

	/*
		Check password attempts against a hash produced using the Encode example
	*/
	hashedPassword := "$argon2id$v=19$t=6,m=196609,p=5$1NTc5r54Kft32HOA/SWYvOjpt6XNTE1MkGoiOsNSwjR9YhoV8guCpIWezymtmcuCODN4PqW0fylGip6yy39o1g$HB2H5fRxY+ev52xWrjoW8w"

	/*
		Incorrect password should return error
	*/
	pwAttempt1 := "password"
	valid := Confirm(pwAttempt1, hashedPassword)
	if !valid {
		fmt.Println("Mismatched Password")
	}

	/*
		Incorrect password should return error
	*/
	pwAttempt2 := "123456"
	valid = Confirm(pwAttempt2, hashedPassword)
	if !valid {
		fmt.Println("Mismatched Password")
	}

	/*
		Correct password, should return nil error
	*/
	valid = Confirm(pw, hashedPassword)
	if !valid {
		fmt.Println("Mismatched Password")
	}
	fmt.Println("Password OK?:", valid)

	// Output:
	// Mismatched Password
	// Mismatched Password
	// Password OK?: true
}
