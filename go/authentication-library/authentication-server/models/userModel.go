package models

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
	UniqueID string
	Name string
	Email string
	hashedPassword string
	Token string
}

// UserStore specifies the requred functions of the user store
type UserStore interface {
	New() *UserStore
	Close()
	FindByID(id string) *User
	FindByEmail(email string) *User
	FindByHashedPassword(hashedPassword string) *User
}

// NewUser returns a pointer to a new, empty User struct 
func NewUser() *User {
	return &User{}
}

