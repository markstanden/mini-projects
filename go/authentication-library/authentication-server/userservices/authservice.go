package userservice

import (
	"fmt"
	"log"

	"github.com/markstanden/argonhasher"
	"github.com/markstanden/authentication"
	"github.com/markstanden/securerandom"
)

type UserService struct {
	UserDS    authentication.UserDataStore
	SecretDS  authentication.SecretDataStore
	AccessTS  authentication.AccessTokenService  // short lived jwt
	RefreshTS authentication.RefreshTokenService // long lived opaque token
}

/*
	** NewUser **
	NewUser creates a new user struct,
	adds the user to the datastore,
	and returns the new user if created ok
*/
func (us UserService) NewUser(name, email, password string) (u *authentication.User, err error) {

	// hash the password, default complexity
	passwordHash := argonhasher.Encode(password, 0)
	if passwordHash == "" {
		return nil, fmt.Errorf("SERVER ERROR - FAILED TO CREATE HASH")
	}

	u = &authentication.User{
		Name:           name,
		Email:          email,
		HashedPassword: passwordHash,
		TokenUserID:    securerandom.String(32),
	}

	if err = us.UserDS.Add(u); err != nil {
		return nil, fmt.Errorf("failed to create user account :\n%v" + err.Error())
	}
	//log.Println("User Account Created OK")
	return u, nil
}

/*
	** GetAccessToken **
	GetAccessToken takes a unique identifier as an argument and returns a unique accessTokenID used to identify the token,
	and the token itself.
*/
func (us UserService) GetAccessToken(tokenUserID string) (accessTokenID, accessToken string) {

	accessToken, accessTokenID, err := us.AccessTS.Create(tokenUserID)
	if err != nil {
		fmt.Printf("/routes/signup: error creating jwt\n%v\nError:\n%v", accessToken, err.Error())
		return "", ""
	}
	log.Println("userservice/GetAccessToken:\n\tEncoded jwt: \n", accessToken, "\n\tjwtid (From Create):\n\t", accessTokenID)
	return accessTokenID, accessToken
}

/*
	** Authenticate Access **
	AuthenticateAccess authenticates the provided access token
	and returns the user information if valid, and an error if invalid
*/
func (us UserService) AuthenticateAccess(jwt string) (user *authentication.User, err error) {

	tokenUserID, accessTokenID, err := us.AccessTS.Decode(jwt)
	if err != nil {
		log.Printf("/routes/signup: error decoding jwt\n%v\nError:\n%v", jwt, err.Error())
		return nil, err
	}

	/*
		basic check logging
	*/
	log.Println("/routes/signup: Decoded JWT OK...")
	log.Println("\n/routes/signup:\nUserID (Decoded from JWT):\n", tokenUserID)
	log.Println("\n/routes/signup:\njwtid (Decoded from JWT):\n", accessTokenID)

	/*
		Lookup the user in the user store
	*/
	user, err = us.UserDS.Find("tokenuserid", tokenUserID)
	if err != nil {
		log.Printf("/routes/signup: error looking up user\n%v\nError:\n%v", tokenUserID, err.Error())
		return nil, authentication.ErrUserNotFound
	}
	log.Printf("/routes/signup: created and decoded jwt.\nJWT String:\n%v\nUserData:\n%v", jwt, user)
	return user, nil
}

/*
	** Login **
	Login takes an email and password argument and attempts to validate the users password.
	If the password validation is successful a pointer to a user is returned,
	if the password validation fails a nil pointer and an error returned.
*/
func (us UserService) Login(email, password string) (user *authentication.User, err error) {
	user, err = us.UserDS.Find("email", email)

	if err != nil {
		fmt.Println("Invalid UserName")
		return nil, err
	}

	/*
		need to use an interface here
	*/
	valid := argonhasher.Confirm(password, user.HashedPassword)
	if !valid {
		return nil, authentication.ErrIncorrectPassword
	} else {
		return user, nil
	}
}
