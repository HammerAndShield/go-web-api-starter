package database

import "errors"

var (
	ErrRecordNotFound        = errors.New("record not found")
	ErrEditConflict          = errors.New("edit conflict")
	ErrForeignKeyViolation   = errors.New("foreign key violation, record is referenced by other records")
	ErrVersionMismatch       = errors.New("version of entity does not match version in database")
	ErrResourceHasDependents = errors.New("resource has dependents")
)

const (
	PsqlUniqueViolation      = "23505"
	PsqlForeignKeyViolation  = "23503"
	PsqlSerializationFailure = "40001"
	PsqlCheckFailure         = "23514"
)
