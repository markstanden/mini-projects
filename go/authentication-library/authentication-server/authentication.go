package authentication

//User is the base struct for our User model.
//	UniqueID string
//		Each user has a unique ID within the UserStore
//	Name string
//		The user's name, for communication purposes
//	Email string
//		The user's email, for communication and ID at login
//	hashedPassword string
//		The hashed password string.  This must never be used to store a plain text password
//	Token string
//		The generated identification token stored within the ID token
type User struct {
	UniqueID       int
	Name           string
	Email          string
	HashedPassword string
	Token          string
}

// UserService specifies the requred functions of the user store
type UserService interface {
	Find(key, value string) (*User, error)
	Add(user *User) error
	FullReset() error
	//Delete(*User)
}

// PasswordHash specifies the requirements of the passord hashing module
type PasswordHash interface {
	Encode(plainTextPassword string) string
	Compare(plainTextPassword, hashedPassword string) error
}

// SecretStore is an interface for the secret storage logic to retrieve secrets,
// which will be platform dependant.
type Deployment interface {
	// Creates a map and fills it with the required information
	//GetSecrets(keys []string) (map[string] string, error)

	// Takes directly from the store
	GetSecret(project, key, version string) (string, error)
}

type TokenHandler interface {
	Create(u *User, secret string) (string, error)
	Decode(jwt, secret string) (map[string]interface{}, error)
}
