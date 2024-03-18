package friend

import "errors"

var (
	ErrInvalidUser           = errors.New("userId is not found")
	ErrAlreadyFriends        = errors.New("userId is already user's friend or adding self as friend")
	ErrNotFriends            = errors.New("userId is not the user's friend")
	ErrValidationFailed      = errors.New("request doesn't pass validation")
	ErrNoFriendJoke_doNotUse = errors.New("u have no friend")
	ErrInvalidToken          = errors.New("request token is missing or expired")
)
