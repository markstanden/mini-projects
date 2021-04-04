package UserService

import "github.com/markstanden/authentication"

type UserService struct{
	uds authentication.UserDataStore
	ats authentication.AccessTokenService	// short lived jwt
	rts authentication.RefreshTokenService	// long lived opaque token
}

