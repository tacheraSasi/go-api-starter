
package exceptions

import (
	"errors"
	"fmt"
)



// NewNotFoundError returns an error indicating a resource was not found
func NewNotFoundError(message string) error {
	return fmt.Errorf("not found: %w", errors.New(message))
}


// NewValidationError returns an error indicating a validation failure
func NewValidationError(message string) error {
	return fmt.Errorf("validation error: %w", errors.New(message))
}

// NewUnauthorizedError returns an error indicating unauthorized access
func NewUnauthorizedError(message string) error {
	return fmt.Errorf("unauthorized: %w", errors.New(message))
}

// NewForbiddenError returns an error indicating forbidden access
func NewForbiddenError(message string) error {
	return fmt.Errorf("forbidden: %w", errors.New(message))
}

// NewConflictError returns an error indicating a conflict (e.g., duplicate resource)
func NewConflictError(message string) error {
	return fmt.Errorf("conflict: %w", errors.New(message))
}

// NewInternalError returns an error indicating an internal server error
func NewInternalError(message string) error {
	return fmt.Errorf("internal error: %w", errors.New(message))
}