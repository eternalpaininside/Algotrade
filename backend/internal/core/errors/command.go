package coreerrors

import "errors"

var (
	ErrNotFound        = errors.New("not found")
	ErrInvalidArgument = errors.New("invalid argument")
	ErrValidation      = errors.New("validation failed")
	ErrConflict        = errors.New("conflict")
	ErrTooLarge        = errors.New("too large")
)
