package exceptions

import (
	"errors"
	"fmt"
	"net/http"
)

// HTTPError represents an error with an associated HTTP status code
type HTTPError struct {
	Code    int
	Message string
	Err     error
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("%s: %v", e.Message, e.Err)
}

// NewNotFoundError returns an error indicating a resource was not found (HTTP 404)
func NewNotFoundError(w http.ResponseWriter, message string) error {
	http.Error(w, message, http.StatusNotFound)
	return &HTTPError{
		Code:    http.StatusNotFound,
		Message: "not found",
		Err:     errors.New(message),
	}
}

// NewValidationError returns an error indicating a validation failure (HTTP 422)
func NewValidationError(w http.ResponseWriter, message string) error {
	http.Error(w, message, http.StatusUnprocessableEntity)
	return &HTTPError{
		Code:    http.StatusUnprocessableEntity,
		Message: "validation error",
		Err:     errors.New(message),
	}
}

// NewUnauthorizedError returns an error indicating unauthorized access (HTTP 401)
func NewUnauthorizedError(w http.ResponseWriter, message string) error {
	http.Error(w, message, http.StatusUnauthorized)
	return &HTTPError{
		Code:    http.StatusUnauthorized,
		Message: "unauthorized",
		Err:     errors.New(message),
	}
}

// NewForbiddenError returns an error indicating forbidden access (HTTP 403)
func NewForbiddenError(w http.ResponseWriter, message string) error {
	http.Error(w, message, http.StatusForbidden)
	return &HTTPError{
		Code:    http.StatusForbidden,
		Message: "forbidden",
		Err:     errors.New(message),
	}
}

// NewConflictError returns an error indicating a conflict (e.g., duplicate resource) (HTTP 409)
func NewConflictError(w http.ResponseWriter, message string) error {
	http.Error(w, message, http.StatusConflict)
	return &HTTPError{
		Code:    http.StatusConflict,
		Message: "conflict",
		Err:     errors.New(message),
	}
}

// NewInternalError returns an error indicating an internal server error (HTTP 500)
func NewInternalError(w http.ResponseWriter, message string) error {
	http.Error(w, message, http.StatusInternalServerError)
	return &HTTPError{
		Code:    http.StatusInternalServerError,
		Message: "internal error",
		Err:     errors.New(message),
	}
}

// NewBadRequestError returns an error indicating a bad request (HTTP 400)
func NewBadRequestError(w http.ResponseWriter, message string) error {
	http.Error(w, message, http.StatusBadRequest)
	return &HTTPError{
		Code:    http.StatusBadRequest,
		Message: "bad request",
		Err:     errors.New(message),
	}
}
