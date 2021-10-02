package repository

import "errors"

var (
	ErrDuplicateEmail     = errors.New("repo: duplicate email")
	ErrInvalidCredentials = errors.New("repo: invalid credentials")
	ErrNoRecord           = errors.New("repo: no record")
	ErrUserNotActivated   = errors.New("repo: not activated")
)
