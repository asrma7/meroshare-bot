package errors

import (
	"errors"
	"net/http"
)

type ErrorType string

const (
	ValidationError   ErrorType = "VALIDATION_ERROR"
	NotFoundError     ErrorType = "NOT_FOUND_ERROR"
	InternalError     ErrorType = "INTERNAL_ERROR"
	UnauthorizedError ErrorType = "UNAUTHORIZED_ERROR"
	ForbiddenError    ErrorType = "FORBIDDEN_ERROR"
	BadRequestError   ErrorType = "BAD_REQUEST_ERROR"
	ConflictError     ErrorType = "CONFLICT_ERROR"
)

type ErrorResponse struct {
	Type    ErrorType         `json:"type"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}

type AppError struct {
	Type    ErrorType         `json:"type"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
	Status  int               `json:"status"`
}

func (e *AppError) Error() string {
	return e.Message
}

type ValidationErrors struct {
	Fields map[string]string `json:"fields"`
}

func NewValidationError(field, message string) *AppError {
	return &AppError{
		Type:    ValidationError,
		Message: "Validation failed",
		Details: map[string]string{field: message},
		Status:  http.StatusUnprocessableEntity,
	}
}

func NewValidationErrors(fields map[string]string) *AppError {
	return &AppError{
		Type:    ValidationError,
		Message: "Validation failed",
		Details: fields,
		Status:  http.StatusUnprocessableEntity,
	}
}

func NewNotFoundError(message string) *AppError {
	return &AppError{
		Type:    NotFoundError,
		Message: message,
		Status:  http.StatusNotFound,
	}
}

func NewInternalError(err error) *AppError {
	return &AppError{
		Type:    InternalError,
		Message: "An Internal error occurred",
		Details: map[string]string{"internal": err.Error()},
		Status:  http.StatusInternalServerError,
	}
}

func NewUnauthorizedError(message string) *AppError {
	return &AppError{
		Type:    UnauthorizedError,
		Message: message,
		Status:  http.StatusUnauthorized,
	}
}

func NewForbiddenError(message string) *AppError {
	return &AppError{
		Type:    ForbiddenError,
		Message: message,
		Status:  http.StatusForbidden,
	}
}

func NewBadRequestError(message string) *AppError {
	return &AppError{
		Type:    BadRequestError,
		Message: message,
		Status:  http.StatusBadRequest,
	}
}

func NewConflictError(message string) *AppError {
	return &AppError{
		Type:    ConflictError,
		Message: message,
		Status:  http.StatusConflict,
	}
}

func IsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	ok := errors.As(err, &appErr)
	return appErr, ok
}

func GetErrorStatusCode(err error) int {
	if appErr, ok := IsAppError(err); ok {
		return appErr.Status
	}
	return http.StatusInternalServerError
}

func Is(err error, errorType error) bool {
	if errors.Is(err, errorType) {
		return true
	}
	return false
}

func GetErrorResponse(err error) (ErrorResponse, int) {
	appErr, ok := IsAppError(err)
	if !ok {
		return ErrorResponse{
			Type:    InternalError,
			Message: "Internal server error",
		}, http.StatusInternalServerError
	}

	errorResp := ErrorResponse{
		Type:    appErr.Type,
		Message: appErr.Message,
	}

	if appErr.Details != nil {
		errorResp.Details = appErr.Details
	}

	return errorResp, appErr.Status
}
