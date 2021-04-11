package userservice

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/markstanden/argonhasher"
	"github.com/markstanden/authentication"
	"github.com/markstanden/securerandom"
)

var (
	ErrInvalidInput = errors.New("invalid input")

	MinInputLength = 3
	MaxInputLength = 255
)

type UserService struct {
	/*
		User data storage
	*/
	UserDS authentication.UserDataStore

	/*
		Application secret storage
	*/
	SecretDS authentication.SecretDataStore

	/*
		Password Protection Package
	*/
	PasswordHasher authentication.PasswordHash

	/*
		Session management
	*/
	AccessTS  authentication.AccessTokenService  // short lived jwt
	RefreshTS authentication.RefreshTokenService // long lived opaque token

	/*
		Config
	*/
	Config USConfig
}

type USConfig struct {
	TokenIDSize      uint
	RefreshTokenSize uint
}

/*
	** NewUserService **
	NewUserService returns a new UserService struct with system defaults
*/
func NewUserService() (us UserService) {
	//us.UserDS = pguserdatastore
	//us.SecretDS
	//us.AccessTS
	//us.RefreshTS
	//us.PasswordHasher = argonhasher
	us.Config.RefreshTokenSize = uint(100)
	us.Config.TokenIDSize = uint(100)
	return us
}

/*
	** NewUser **
	NewUser creates a new user struct,
	adds the user to the datastore,
	and returns the new user if created ok
*/
func (us UserService) NewUser(name, email, password string) (u *authentication.User, err error) {

	if emptyString(name) ||
		tooLong(name) ||
		emptyString(email) ||
		tooLong(email) ||
		!validEmail(email) ||
		emptyString(password) ||
		tooLong(password) ||
		tooShort(password) {
		return nil, ErrInvalidInput
	}

	// hash the password, default complexity
	passwordHash := argonhasher.Encode(password, 0)
	if passwordHash == "" {
		return nil, fmt.Errorf("SERVER ERROR - FAILED TO CREATE HASH")
	}

	u = new(authentication.User)
	u.Name = name
	u.Email = email
	u.HashedPassword = passwordHash
	u.TokenUserID = securerandom.String(us.Config.TokenIDSize)
	u.CurrentRefreshToken = us.NewRefreshToken()

	if err = us.UserDS.Add(u); err != nil {

		return nil, fmt.Errorf("failed to create user account :\n%v" + err.Error())
	}

	return u, nil
}

/*
	** NewRefreshToken **
	NewRefreshToken generates a new refresh token and returns the result as a string
*/
func (us UserService) NewRefreshToken() (refreshToken string) {
	refreshToken = securerandom.String(us.Config.RefreshTokenSize)
	return refreshToken
}

/*
	** UpdateRefreshToken **
	Helper function that gets a new refresh token, and updates the local and stored version.
*/
func (us UserService) UpdateRefreshToken(user *authentication.User) {
	us.UpdateUser(user, UpdateRefreshToken(us.NewRefreshToken()))
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

	/*
		check inputs to potentially save unnecessary hashing or DB lookups
	*/
	if emptyString(email) || emptyString(password) {
		return nil, ErrInvalidInput
	}

	if tooShort(email) || !validEmail(email) || tooLong(email) {
		return nil, ErrInvalidInput
	}

	if tooShort(password) || tooLong(password) {
		return nil, ErrInvalidInput
	}

	user, err = us.UserDS.Find("email", email)

	if err != nil {
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

/*
	The closure function to set the fields of the update struct
*/
type UpdateFunc func(us UserService, updates *authentication.User) error

/*
	** UpdateName **
	UpdateName is a configuration function passed to the UpdateUser
	function as argument, used to add the new name to the configuration
	struct.
*/
func UpdateName(name string) UpdateFunc {
	return UpdateFunc(func(us UserService, updates *authentication.User) error {
		updates.Name = name
		return nil
	})
}

/*
	** UpdateEmail **
	UpdateEmail is a configuration function passed to the UpdateUser
	function as argument, used to add the new email address to the configuration
	struct.
*/
func UpdateEmail(email string) UpdateFunc {
	return UpdateFunc(func(us UserService, updates *authentication.User) error {
		updates.Email = email
		return nil
	})
}

/*
	** UpdatePassword **
	UpdatePassword is a configuration function passed to the UpdateUser
	function as argument, used to add the new hashed password to the configuration
	struct.
*/
func UpdatePassword(plaintext string) UpdateFunc {
	return UpdateFunc(func(us UserService, updates *authentication.User) error {
		newHashedPW := argonhasher.Encode(plaintext, 0)
		if newHashedPW == "" {
			return ErrInvalidInput
		}
		updates.HashedPassword = newHashedPW
		return nil
	})
}

/*
	*** UpdateRefreshToken ***
	Updates the current users refresh token.
*/
func UpdateRefreshToken(refreshtoken string) UpdateFunc {
	return UpdateFunc(func(us UserService, updates *authentication.User) error {
		updates.CurrentRefreshToken = refreshtoken
		return nil
	})
}

/*
	*** UpdateUser ***
	Updates the current user with the updated fields within the struct.
*/
func (us UserService) UpdateUser(u *authentication.User, updates ...UpdateFunc) (err error) {
	updatedFields := new(authentication.User)
	for _, update := range updates {
		if err = update(us, updatedFields); err != nil {
			return err
		}

	}
	return us.UserDS.Update(u, *updatedFields)
}

/*
	*** emptyString ***
	emptyString checks for an empty input string and true if empty.
*/
func emptyString(input string) bool {
	if input == "" {
		return true
	}
	return false
}

/*
	*** tooLong ***
	tooLong returns true if the input string is too long
	(longer than MaxInputLength)
*/
func tooLong(input string) bool {
	if len(input) > MaxInputLength {
		return true
	}
	return false
}

/*
	*** tooShort ***
	tooShort returns true if the input string is too short
	(shorter than MinInputLength)
*/
func tooShort(input string) bool {
	if len(input) < MinInputLength {
		return true
	}
	return false
}
func validEmail(input string) bool {
	parts := strings.Split(input, "@")
	if len(parts) != 2 {
		return false
	}
	domain := strings.Split(parts[1], ".")
	if len(domain) != 2 {
		return false
	}
	return true
}
