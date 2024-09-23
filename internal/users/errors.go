package users

import "errors"

var (
	ErrRoleNotFound   = errors.New("role not found in database")
	ErrDuplicateEmail = errors.New("duplicate email in users table")
)
