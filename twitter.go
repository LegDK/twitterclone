package twitterclone

import "errors"

var (
	ErrBadCredentials     = errors.New("error/password wrong combination")
	ErrValidation         = errors.New("validation error")
	ErrNotFound           = errors.New("not found")
	ErrInvalidAccessToken = errors.New("invalid access token")
	ErrGenAccessToken     = errors.New("error generate access token")
	ErrForbidden          = errors.New("forbidden")
	ErrUnauthenticated    = errors.New("unauthenticated")
)
