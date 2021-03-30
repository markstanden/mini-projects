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
//	TokenID string
//		The generated identification token stored within the ID token
type User struct {
	UniqueID       int
	Name           string
	Email          string
	HashedPassword string
	TokenID        string
}

type Secret struct {
	KeyName string
	KeyID   string
	Value   string
	Created int64
}

// UserService specifies the requred functions of the user store
type UserService interface {
	DataStore
	Find(key, value string) (*User, error)
	Add(user *User) error
	//Delete(*User)
}

// SecretService is an interface for the secret storage logic to retrieve secrets,
// which will be platform dependant.
type SecretService interface {
	DataStore
	// Takes directly from the store
	GetSecret(keyName string) func(keyID string) (secret string)
	GetKeyID(keyName string) (keyID string)
}

type DataStore interface {
	FullReset() (err error)
}

// PasswordHash specifies the requirements of the passord hashing module
type PasswordHash interface {
	Encode(plainTextPassword string) string
	Compare(plainTextPassword, hashedPassword string) bool
}

type TokenService interface {
	Create(userID string) (jwt, jwtID string, err error)
	Decode(jwt string) (userID, jwtID string, err error)
	//GetSecret(version string) (secret string)
}

var (
	// users
	ErrUserNotFound = errors.New("user not found")

	// tokens
	ErrExpiredToken = errors.New("expired token")
)
