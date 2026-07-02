package customerrors

import (
	"errors"
	"fmt"
)

type NotFoundError struct {
	Resource string
	ID       int
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("resource '%s' with ID %d was not found", e.Resource, e.ID)
}

type ValidationError struct {
	Field  string
	Reason string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed on field '%s': %s", e.Field, e.Reason)
}

// ErrDatabaseTimeout is our sentinel error
var ErrDatabaseTimeout = errors.New("database connection timed out")