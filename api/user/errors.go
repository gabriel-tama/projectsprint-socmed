package user

import "errors"

var (
	ErrUserNotFound          = errors.New("user not found")
	ErrWrongPassword         = errors.New("wrong password")
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrValidationFailed      = errors.New("validation failed")
	ErrCredentialType        = errors.New("invalid cred type")
	ErrCredentialValue       = errors.New("invalid cred value")
	ErrEmailAlreadyExists    = errors.New("email already exists")
	ErrPhoneAlreadyExists    = errors.New("phone already exists")
	ErrInvalidToken          = errors.New("token invalid or missing")
	ErrWrongRoute            = errors.New("cant use this route to update")
)
