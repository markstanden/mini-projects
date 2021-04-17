package authentication

import "errors"

//User is the base struct for our User model.
//	UniqueID string
//		Each user has a unique ID within the UserStore
//	Name string
//		The user's name, for communication purposes
//	Email string
//		The user's email, for communication and ID at login
//	hashedPassword string
//		The hashed password string.  This must never be used to store a plain text password
//	TokenUserID string
//		The generated identification token stored within the ID token
type User struct {
	UniqueID             int
	Name                 string
	Email                string
	HashedPassword       string
	TokenUserID          string
	CurrentRefreshToken  string
	CurrentAccessTokenID string
}

type Secret struct {
	KeyName string
	KeyID   string
	Value   string
	Created int64
}

// UserService specifies the requred functions of the user store
type UserDataStore interface {
	DataStore
	Add(user *User) error
	Find(key, value string) (*User, error)
	Update(user *User, updated User) error
	UpdateRefreshToken(user *User, newRefreshToken string) error
	Delete(user *User) error
}

/*
	** SecretService ** is an interface for the secret storage logic to retrieve secrets,
	which will be platform dependant.
*/
type SecretDataStore interface {
	DataStore
	// Takes directly from the store
	GetSecret(keyName string) func(keyID string) (secret string)
	GetKeyID(keyName string) (keyID string)
}

/*
	** DataStore **
	A DataStore holds a connection to a datastore
*/
type DataStore interface {
	FullReset() (err error)
}

/*
	** PasswordHash **
	specifies the requirements of the passord hashing module
*/
type PasswordHash interface {
	Encode(plainTextPassword string) string
	Compare(plainTextPassword, hashedPassword string) bool
}

/*
	** Access Token Service **
	The required methods to create and verify the access tokens
	used to authorise users within areas of our site.
*/
type AccessTokenService interface {
	Create(userID string) (jwt, jwtID string, err error)
	Decode(jwt string) (userID, jwtID string, err error)
	//GetSecret(version string) (secret string)
}

/*
	** Refresh Token Service **
	The required methods to issue and verify the refresh token implementation
	for the site.  The refresh token is used as a long duration token used to identify the session,
	and is generally a high entropy private unique string stored within a cookie
*/
type RefreshTokenService interface {
	Create(userID string) (refreshToken string)
}

/*
	** UserService **
	The userservice combines all of the components of the user login and authentication system
	including creation and addition of a new user, the logging in and creation of new tokens,
	and the verification of existing tokens
*/
type UserService interface {
	NewUser(name, email, password string) (user *User, err error)
	Login(email, password string) (user *User, err error)
}

/*
	** ERRORS **
	list of accepted errors to be used and checked against
*/
var (
	// users
	ErrUserNotFound      = errors.New("user not found")
	ErrIncorrectPassword = errors.New("incorrect password")

	// tokens
	ErrExpiredToken = errors.New("expired token")

	// internal error
	ErrInternalServerError = errors.New("internal server error")
)
