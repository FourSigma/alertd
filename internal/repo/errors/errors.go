package errors

import "errors"

var (
	ErrorDoesNotExist    = errors.New("object does not exists")
	ErrorAlreadyExists   = errors.New("object already exists")
	ErrorUniqueViolation = errors.New("unique violations")
)
