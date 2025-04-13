package auth

import "errors"

var (
	ErrUserNotFound             = errors.New("user not found")
	ErrUserAlreadyExists        = errors.New("user already exists")
	ErrIncorrectPasswordOrEmail = errors.New("incorrect password or email")
	ErrIncorrectRole            = errors.New("incorrect access level")
	ErrBadRequest               = errors.New("bad request")
)
