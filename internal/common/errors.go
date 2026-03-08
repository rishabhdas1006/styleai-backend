package common

import "errors"

var (
	ErrEmailExists        = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailNotFound      = errors.New("email not found")
	ErrMissingAuthToken   = errors.New("missing auth token")
	ErrInvalidAuthToken   = errors.New("invalid auth token")
)
