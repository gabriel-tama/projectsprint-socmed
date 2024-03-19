package post

import "errors"

var (
	ErrUserNotFound          = errors.New("user not found")
	ErrWrongPassword         = errors.New("wrong password")
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrValidationFailed      = errors.New("validation failed")
	ErrCredentialType        = errors.New("invalid cred type")
	ErrCredentialValue       = errors.New("invalid cred value")
)
