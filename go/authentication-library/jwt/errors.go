package jwt

import "errors"

// Errors
var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
	ErrFailedSecret = errors.New("failed to retrieve secret")
)
