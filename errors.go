package nasfilesapi

import "errors"

var (
	// ErrUserNotFound is thrown whenever a user is not found
	ErrUserNotFound = errors.New("User not found.")

	// ErrPasswordTooShort is thrown whenever a user is trying to change the password
	// using a new password with less than 4 characters
	ErrPasswordTooShort = errors.New("Password too small. Minimum length is 6 characters")

	// ErrWrongPassword is thrown whenever a user is trying to change the password
	// using a new password with less than 4 characters
	ErrWrongPassword = errors.New("Password is incorrect.")
)
