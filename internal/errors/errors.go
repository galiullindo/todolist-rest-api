package errors

import "errors"

var (
	ErrUnexpectedRowsAffected = errors.New("unexpected rows affected")

	ErrTaskNotFound = errors.New("task not found")
)
