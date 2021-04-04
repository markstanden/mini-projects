package userservice

import "github.com/markstanden/authentication"

type UserService struct {
	UserDS    authentication.UserDataStore
	SecretDS  authentication.SecretDataStore
	AccessTS  authentication.AccessTokenService  // short lived jwt
	RefreshTS authentication.RefreshTokenService // long lived opaque token
}

func (us UserService) NewUser() (u authentication.User, err error) {
	return
}
